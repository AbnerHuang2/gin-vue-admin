# 构建阶段
FROM node:18-alpine as builder

WORKDIR /app

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY . .

# 构建生产版本
RUN npm run build

# 生产阶段 - 使用 nginx 服务静态文件
FROM nginx:alpine

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/nginx.conf

# 从构建阶段复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# Railway 使用 PORT 环境变量
ENV PORT=80
EXPOSE ${PORT}

# 启动 nginx（使用 envsubst 替换环境变量）
CMD envsubst '${PORT} ${BACKEND_URL}' < /etc/nginx/nginx.conf > /etc/nginx/nginx.conf.tmp \
    && mv /etc/nginx/nginx.conf.tmp /etc/nginx/nginx.conf \
    && nginx -g 'daemon off;'
