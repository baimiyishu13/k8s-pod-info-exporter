FROM golang:1.22.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o k8s-pod-info-exporter ./cmd

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y ca-certificates wget tar && rm -rf /var/lib/apt/lists/*
RUN useradd -m appuser
WORKDIR /app

# Create bin directory and download/extract the tarball
RUN mkdir -p /app/k8s-pod-info-exporter/bin
RUN wget -O /app/k8s-pod-info-exporter/bin/k8s-resource-exporter-mac.tar.gz https://gitlab.com/baimiyishu13/k8s-resourc-exporter/-/jobs/7402538260/artifacts/raw/k8s-resource-exporter-mac.tar.gz
RUN tar -xzf /app/k8s-pod-info-exporter/bin/k8s-resource-exporter-mac.tar.gz -C /app/k8s-pod-info-exporter/bin

COPY --from=builder /app/k8s-pod-info-exporter /app/k8s-pod-info-exporter
COPY templates /app/templates
RUN chown -R appuser /app
USER appuser
ENTRYPOINT ["/app/k8s-pod-info-exporter"]
EXPOSE 8080
