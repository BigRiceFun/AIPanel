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

- Frontend: Vue 3, Vite, TypeScript, Element Plus, Pinia, Vue Router, Axios, ECharts
- Backend: Go, Gin, SQLite, GORM, JWT, gopsutil, Docker SDK

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

## Docker 部署

```bash
docker compose up -d --build
```

访问：

```text
http://localhost:3000
```

后端 API 直连端口：

```text
http://localhost:8080
```

## 配置

根目录 `config.yaml` 控制服务端口、JWT、SQLite 路径和初始管理员。
生产环境请修改 `jwt.secret` 和默认管理员密码。

如果生产环境需要浏览器直接访问后端 API，可用环境变量覆盖 CORS 来源，避免把服务器 IP 或域名写进仓库：

```bash
AIPANEL_CORS_ORIGINS="http://your-domain:3000,http://localhost:5173"
```
