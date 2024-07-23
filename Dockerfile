# 第一阶段：构建 Go 应用程序
FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 禁用 CGO，构建与 Alpine 兼容的二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o k8s-pod-info-exporter ./cmd

# 第二阶段：创建最终镜像
FROM alpine:latest

WORKDIR /app

# 安装 wget 和 ca-certificates
RUN apk add --no-cache wget ca-certificates

# 复制第一阶段构建的二进制文件
COPY --from=builder /app/k8s-pod-info-exporter .
COPY . .
# 确认文件存在并具有执行权限
RUN ls -l /app/k8s-pod-info-exporter
RUN chmod +x /app/k8s-pod-info-exporter
RUN ls -l /app/k8s-pod-info-exporter

# 下载和解压 k8s-resource-exporter
RUN mkdir -p bin \
    && wget -O bin/k8s-resource-exporter-adm64.tar.gz https://gitlab.com/baimiyishu13/k8s-resourc-exporter/-/jobs/7402538260/artifacts/raw/k8s-resource-exporter-adm64.tar.gz \
    && tar -xzf bin/k8s-resource-exporter-adm64.tar.gz -C bin \
    && mv bin/k8s-resource-exporter-adm64 bin/k8s-resource-exporter \
    && rm -f bin/k8s-resource-exporter-adm64.tar.gz

ENTRYPOINT ["./k8s-pod-info-exporter"]
EXPOSE 8080
