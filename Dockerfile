# 使用 golang 官方提供的镜像作为基础
FROM golang:1.21

#ENV http_proxy http://10.98.163.254
ENV GOPROXY=https://goproxy.io,direct

# 设置工作目录
WORKDIR /go/src/app

RUN go mod init go-project

# 将项目文件复制到镜像中
COPY . .

# 编译应用程序
RUN go build -o app

# 暴露应用程序的端口（如果有需要）
EXPOSE 8080

# 安装Python 3
#RUN apt-get update && apt-get install -y python3 && apt-get -y install python3-pip

#RUN pip3 install RPi.GPIO

# 设置容器启动时的命令
CMD ["./app"]
