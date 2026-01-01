# 构建阶段
FROM node:18-alpine as builder

WORKDIR /app

# 先只复制 package 文件（利用 Docker 缓存）
# 如果 package.json 没变，npm install 这一层会被缓存
COPY package*.json ./

# 安装依赖（使用缓存）
RUN npm install --prefer-offline --no-audit --progress=false

# 再复制源代码（源代码变化不会重新 npm install）
COPY . .

# 构建生产版本
RUN npm run build

# 生产阶段 - 使用 nginx 服务静态文件
FROM nginx:alpine

# 复制 nginx 配置模板
COPY nginx.conf.template /etc/nginx/nginx.conf.template

# 从构建阶段复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 创建启动脚本
RUN echo '#!/bin/sh' > /docker-entrypoint.sh && \
    echo 'set -e' >> /docker-entrypoint.sh && \
    echo 'export PORT=${PORT:-80}' >> /docker-entrypoint.sh && \
    echo 'export BACKEND_HOST=${BACKEND_HOST:-skitii-gva-api.up.railway.app}' >> /docker-entrypoint.sh && \
    echo 'envsubst '"'"'${PORT} ${BACKEND_HOST}'"'"' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf' >> /docker-entrypoint.sh && \
    echo 'exec nginx -g "daemon off;"' >> /docker-entrypoint.sh && \
    chmod +x /docker-entrypoint.sh

ENV PORT=80
ENV BACKEND_HOST=skitii-gva-api.up.railway.app

EXPOSE ${PORT}

ENTRYPOINT ["/docker-entrypoint.sh"]
