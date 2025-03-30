FROM golang:1.20-alpine

# 安装必要的系统依赖
RUN apk add --no-cache git postgresql-client

# 设置 GOPROXY 以加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct
# 禁用 Go modules 校验，跳过 checksum 验证
ENV GOSUMDB=off
ENV GO111MODULE=on

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download -x

# 复制源代码和等待脚本
COPY . .

# 创建存储目录
RUN mkdir -p /app/storage /app/graphql/logs

# 构建应用
RUN go build -o spurtcms .

# 确保等待脚本可执行
RUN chmod +x /app/wait-for-postgres.sh

# 暴露端口
EXPOSE 8082 8083 8084

# 使用等待脚本确保数据库准备好后再启动应用
CMD ["/app/wait-for-postgres.sh", "postgres", "./spurtcms"]
