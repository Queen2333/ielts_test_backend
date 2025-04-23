# 使用官方的 Golang 镜像作为基础镜像
# FROM golang:latest AS builder
FROM golang:1.21-alpine AS builder

# 在容器中创建一个目录来存放项目代码
WORKDIR /app

# 复制 go.mod 和 go.sum，下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 将本地代码复制到容器中的工作目录
COPY . .

# 使用 go mod 下载项目依赖
# RUN go mod download

# 编译 Go 应用
# RUN CGO_ENABLED=0 GOOS=linux go build -o app .
# 编译 Go 程序，禁用 CGO，针对 ARM 架构
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o app .

# 使用轻量级 Alpine 镜像运行
FROM alpine:3.18
# 创建一个小镜像
# FROM alpine:latest

# 在容器中创建一个目录来存放应用程序
WORKDIR /root/

# 从 builder 镜像中将编译好的应用程序复制到容器中
COPY --from=builder /app/app .

# 声明服务端口
EXPOSE 8081

# 运行应用程序
CMD ["./app"]
