# k8s-pod-info-exporter

`k8s-pod-info-exporter` 是一个用于导出 Kubernetes Pod 信息的工具，并通过一个简单的 Web 界面提供文件上传和下载功能。

## 功能

- 上传 kubeconfig 文件
- 生成 Kubernetes Pod 信息的 CSV 文件
- 下载生成的 CSV 文件

## 前提条件

- [Docker](https://www.docker.com/get-started) 已安装
- [GitLab CI](https://docs.gitlab.com/ee/ci/) 配置

## 快速开始

### 使用 Docker

1. 构建 Docker 镜像：

    ```sh
    docker build -t k8s-pod-info-exporter .
    ```

2. 运行 Docker 容器：

    ```sh
    docker run -p 8080:8080 k8s-pod-info-exporter
    ```

3. 打开浏览器并访问 `http://localhost:8080`。



..... 待
