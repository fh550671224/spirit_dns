# 使用官方Go镜像作为构建环境
FROM golang:1.18 as builder

# 设置工作目录
WORKDIR /app

# 复制Go模块的依赖文件
COPY go.mod .
COPY go.sum .

# 下载依赖项
RUN go mod download

# 复制源代码到容器中
COPY . .

# 编译Go程序为二进制文件。使用CGO_ENABLED=0来构建一个静态链接的二进制文件。
RUN CGO_ENABLED=0 GOOS=linux go build -v -o spirit_dns

# 使用scratch作为最终镜像的基础镜像，这是一个空白的镜像
FROM scratch

# 将工作目录设置为 /
WORKDIR /

# 从构建阶段的镜像中复制二进制文件到当前目录
COPY --from=builder /app/spirit_dns .

# 运行编译好的二进制文件
CMD ["./spirit_dns"]