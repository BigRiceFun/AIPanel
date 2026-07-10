# AIPanel

Lightweight AI-first Server Panel.

当前已实现的基础面板能力：

- 用户登录认证
- Dashboard 系统监控
- Docker 容器管理
- Web Terminal
- Docker / 系统日志查看
- AI Provider 基础配置
- 操作审计记录

暂未实现 AI Agent、Workflow Engine、Plugin Marketplace、自动执行任务。

## 技术栈

- Frontend: Vue 3, Vite, TypeScript, Element Plus, Pinia, Vue Router, Axios, ECharts, xterm.js
- Backend: Go, Gin, SQLite, GORM, JWT, gopsutil, Docker SDK, WebSocket

## 已实现 API

```text
POST   /api/auth/login
GET    /api/system/status
GET    /api/system/logs
GET    /api/docker/containers
GET    /api/docker/logs/:id
POST   /api/docker/start/:id
POST   /api/docker/stop/:id
POST   /api/docker/restart/:id
DELETE /api/docker/remove/:id
GET    /api/terminal/ws
GET    /api/ai/config
POST   /api/ai/config
POST   /api/ai/test
GET    /api/audit/logs
```

AI Provider 当前支持：

```text
openai_chat       OpenAI Chat Completions / OpenAI Compatible
openai_responses  OpenAI Responses API
gemini            Gemini GenerateContent
anthropic         Anthropic Messages API
```

## 本地开发

后端：

```bash
cd server
go mod tidy
$env:AIPANEL_CONFIG="../config.yaml"
go run .
```

前端：

```bash
cd frontend
npm install
npm run dev
```

默认账号：

```text
username: admin
password: admin123
```

## 宿主机运行 server

生产环境推荐把 server 直接运行在宿主机，这样 Web Terminal 操作的就是宿主机 shell。

```bash
cd server
go mod tidy
go build -o aipanel-server .
cd ..
AIPANEL_CONFIG=./config.yaml ./server/aipanel-server
```

使用 systemd：

```bash
sudo mkdir -p /opt/aipanel
sudo cp -r server deploy config.yaml /opt/aipanel/
cd /opt/aipanel/server
go build -o aipanel-server .
sudo cp /opt/aipanel/deploy/aipanel-server.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now aipanel-server
```

> 如果使用 systemd，请根据实际部署路径调整 `deploy/aipanel-server.service` 中的 `/opt/aipanel`。

## 前端 Docker 部署

```bash
docker compose up -d --build
```

当前 `docker-compose.yml` 只启动 frontend，Nginx 会把 `/api` 代理到宿主机：

```text
http://host.docker.internal:8080
```

访问：

```text
http://localhost:3000
```

后端 API 直连端口：

```text
http://localhost:8080
```

Web Terminal 通过前端 Nginx 代理 WebSocket：

```text
/api/terminal/ws
```

## 配置

根目录 `config.yaml` 控制服务端口、JWT、SQLite 路径、初始管理员、Terminal 和 AI Provider 默认配置。

生产环境请修改：

- `jwt.secret`
- 默认管理员密码
- AI Provider API Key

如果生产环境需要浏览器直接访问后端 API，可用环境变量覆盖 CORS 来源，避免把服务器 IP 或域名写进仓库：

```bash
AIPANEL_CORS_ORIGINS="http://your-domain:3000,http://localhost:5173"
```
