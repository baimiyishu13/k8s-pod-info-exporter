package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// clearDirectory 清空指定目录下的所有文件
func clearDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}

// getCurrentDate 获取当前日期，格式为YYYYMMDD
func getCurrentDate() string {
	return time.Now().Format("20060102")
}

// UploadHandler 处理文件上传和命令执行
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 清空 ./config 目录
	configDir := "./config"
	err := clearDirectory(configDir)
	if err != nil {
		log.Printf("Error clearing config directory: %v", err)
		http.Error(w, "Failed to clear config directory", http.StatusInternalServerError)
		return
	}

	// 清空 ./output 目录
	outputDir := "./output"
	err = clearDirectory(outputDir)
	if err != nil {
		log.Printf("Error clearing output directory: %v", err)
		http.Error(w, "Failed to clear output directory", http.StatusInternalServerError)
		return
	}

	// 解析multipart表单，允许最多10 MB的文件上传
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// 从表单数据中获取文件
	file, header, err := r.FormFile("kubeconfig")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		http.Error(w, "Unable to get the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 使用上传文件的原始名称
	originalFilename := header.Filename
	filePath := filepath.Join(configDir, originalFilename) // 将文件保存在 ./config 目录中

	// 在指定路径创建一个新文件
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 将上传文件的内容复制到新文件中
	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}

	// config 目录下文件的相对路径
	relPath := filepath.Join("config", originalFilename)

	// 获取当前日期
	currentDate := getCurrentDate()
	outputFileName := fmt.Sprintf("k8s_resources_%s.csv", currentDate)

	// 确保输出目录存在，不存在则创建
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			log.Printf("Error creating output directory: %v", err)
			http.Error(w, "Unable to create output directory", http.StatusInternalServerError)
			return
		}
	}

	// 使用指定参数执行k8s-resource-exporter命令
	cmd := exec.Command("./bin/k8s-resource-exporter", "--config", relPath, "--output", outputDir)
	err = cmd.Run()
	if err != nil {
		log.Printf("Error executing command: %v", err)
		http.Error(w, "Failed to execute command", http.StatusInternalServerError)
		return
	}

	// 处理完后清理上传的kubeconfig文件
	defer os.Remove(filePath)

	// 将CSV文件的路径写入到响应中，以便前端可以下载
	w.Header().Set("HX-Redirect", fmt.Sprintf("/api/download/%s", outputFileName))
	w.WriteHeader(http.StatusOK)
}
