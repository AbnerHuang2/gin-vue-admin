FROM golang:alpine as builder

WORKDIR /go/src/github.com/AbnerHuang2/gin-vue-admin/server
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

WORKDIR /go/src/github.com/flipped-aurora/gin-vue-admin/server

COPY --from=0 /go/src/github.com/flipped-aurora/gin-vue-admin/server/server ./
COPY --from=0 /go/src/github.com/flipped-aurora/gin-vue-admin/server/resource ./resource/
COPY --from=0 /go/src/github.com/flipped-aurora/gin-vue-admin/server/config.docker.yaml ./
COPY --from=0 /go/src/github.com/flipped-aurora/gin-vue-admin/server/config.railway.yaml ./

# 创建上传目录
RUN mkdir -p uploads/file

# 挂载目录：如果使用了sqlite数据库，容器命令示例：docker run -d -v /宿主机路径/gva.db:/go/src/github.com/flipped-aurora/gin-vue-admin/server/gva.db -p 8888:8888 --name gva-server-v1 gva-server:1.0
# VOLUME ["/go/src/github.com/flipped-aurora/gin-vue-admin/server"]

# Railway 使用 PORT 环境变量
ENV PORT=8888
EXPOSE ${PORT}

# 默认使用 docker 配置，Railway 会通过 railway.toml 覆盖启动命令
ENTRYPOINT ./server -c config.docker.yaml
