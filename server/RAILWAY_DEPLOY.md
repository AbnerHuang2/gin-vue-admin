# Railway 部署指南

## 快速部署步骤

### 1. 登录 Railway

访问 [railway.app](https://railway.app) 并登录

### 2. 创建新项目

1. 点击 **"New Project"**
2. 选择 **"Deploy from GitHub repo"**
3. 选择你的 `gin-vue-admin` 仓库

### 3. 配置项目

1. 点击项目进入设置
2. 在 **Settings** → **Root Directory** 中设置为 `server`
3. Railway 会自动检测 Dockerfile

### 4. 添加 MySQL 数据库

1. 在项目中点击 **"+ New"**
2. 选择 **"Database"** → **"Add MySQL"**
3. Railway 会自动创建并注入以下环境变量：
   - `MYSQL_HOST`
   - `MYSQL_PORT`
   - `MYSQL_DATABASE`
   - `MYSQL_USER`
   - `MYSQL_PASSWORD`

### 5. 配置环境变量

在项目的 **Variables** 中添加以下环境变量：

#### 必需的环境变量

| 变量名 | 说明 | 示例值 |
|--------|------|--------|
| `JWT_SIGNING_KEY` | JWT 签名密钥（请使用随机字符串） | `your-super-secret-jwt-key-2024` |
| `FRONTEND_URL` | 前端 Vercel 部署地址（用于 CORS） | `https://your-app.vercel.app` |

#### 可选的环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `PORT` | 服务端口（Railway 自动设置） | `8888` |
| `USE_REDIS` | 是否启用 Redis | `false` |
| `REDIS_URL` | Redis 连接地址 | - |
| `REDIS_PASSWORD` | Redis 密码 | - |
| `MCP_PORT` | MCP 服务端口 | `8889` |

### 6. 部署

配置完成后，Railway 会自动部署。你可以在 **Deployments** 中查看部署日志。

### 7. 获取后端 URL

部署成功后，在 **Settings** → **Networking** → **Public Networking** 中生成公开域名。

---

## 数据库初始化

首次部署后，需要通过前端进行数据库初始化：

1. 访问前端页面
2. 系统会自动跳转到初始化页面
3. 数据库信息会自动填充（从环境变量读取）
4. 点击初始化完成数据库配置

---

## 前端 Vercel 配置

在 Vercel 中配置以下环境变量：

```env
VITE_BASE_API=/api
VITE_BASE_PATH=https://你的railway域名
VITE_SERVER_PORT=443
VITE_CLI_PORT=8080
VITE_POSITION=close
```

同时在 `web/vercel.json` 中配置 API 代理：

```json
{
  "rewrites": [
    {
      "source": "/api/:path*",
      "destination": "https://你的railway域名/:path*"
    },
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ]
}
```

---

## 添加 Redis（可选）

如果需要使用 Redis：

1. 在 Railway 项目中点击 **"+ New"** → **"Database"** → **"Add Redis"**
2. 添加环境变量：
   - `USE_REDIS=true`
   - Railway 会自动注入 `REDIS_URL`

---

## 常见问题

### Q: 部署失败，提示找不到配置文件？

确保 `railway.toml` 中的启动命令正确：
```toml
[deploy]
startCommand = "./server -c config.railway.yaml"
```

### Q: 数据库连接失败？

1. 确保 MySQL 服务已添加
2. 检查环境变量是否正确注入
3. 查看部署日志获取详细错误信息

### Q: CORS 错误？

确保 `FRONTEND_URL` 环境变量设置正确，且包含完整的 URL（如 `https://your-app.vercel.app`）

### Q: 健康检查失败？

- 健康检查端点：`/health`
- 确保应用正常启动
- 检查端口配置是否正确

---

## 环境变量完整参考

```bash
# 必需
JWT_SIGNING_KEY=your-jwt-signing-key
FRONTEND_URL=https://your-frontend.vercel.app

# MySQL（Railway 自动注入）
MYSQL_HOST=xxx
MYSQL_PORT=3306
MYSQL_DATABASE=railway
MYSQL_USER=root
MYSQL_PASSWORD=xxx

# 可选 - Redis
USE_REDIS=false
REDIS_URL=redis://xxx
REDIS_PASSWORD=xxx

# 可选 - 其他
MCP_PORT=8889
```

---

## 本地测试 Railway 配置

```bash
# 设置环境变量
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_DATABASE=gva
export MYSQL_USER=root
export MYSQL_PASSWORD=your_password
export JWT_SIGNING_KEY=test-key
export PORT=8888

# 运行
go run . -c config.railway.yaml
```

