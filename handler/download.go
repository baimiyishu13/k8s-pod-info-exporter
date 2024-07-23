package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// DownloadHandler 处理CSV文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取文件名
	filename := chi.URLParam(r, "filename")
	filepath := filepath.Join("output", filename) // 构建文件路径

	// 设置头部信息以启动文件下载
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	// 将文件发送给客户端
	http.ServeFile(w, r, filepath)
}
