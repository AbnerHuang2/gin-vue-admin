FROM golang:alpine as builder

# 使用简单的工作目录，避免路径问题
WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="SliverHorn@sliver_horn@qq.com"

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk update && apk add --no-cache tzdata openntpd ca-certificates \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# 从构建阶段复制文件
COPY --from=builder /app/server ./
COPY --from=builder /app/resource ./resource/
COPY --from=builder /app/config.docker.yaml ./
COPY --from=builder /app/config.railway.yaml ./

# 创建上传目录
RUN mkdir -p uploads/file

# Railway 使用 PORT 环境变量
ENV PORT=8888
EXPOSE ${PORT}

# 默认使用 docker 配置，Railway 会通过 railway.toml 覆盖启动命令
ENTRYPOINT ["./server", "-c", "config.docker.yaml"]
