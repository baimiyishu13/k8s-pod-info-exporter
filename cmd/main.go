package main

import (
	"k8s-pod-info-exporte/handler"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Logger中间件，用于记录HTTP请求日志
	r.Use(middleware.Recoverer) // Recoverer中间件，用于从panic中恢复

	r.Post("/api/upload", handler.UploadHandler)               // 文件上传的路由
	r.Get("/api/download/{filename}", handler.DownloadHandler) // 文件下载的路由
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html") // 主页路由，提供前端页面
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r) // 在8080端口启动HTTP服务器
}
