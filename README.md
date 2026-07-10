# AIPanel

Lightweight AI-first Server Panel.

当前第一阶段聚焦基础面板能力：

- 用户登录认证
- Dashboard 系统监控
- Docker 容器管理

暂未实现 AI Chat、Workflow、Plugin、Terminal、日志分析。

## 技术栈

- Frontend: Vue 3, Vite, TypeScript, Element Plus, Pinia, Vue Router, Axios, ECharts
- Backend: Go, Gin, SQLite, GORM, JWT, gopsutil, Docker SDK

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

## 配置

根目录 `config.yaml` 控制服务端口、JWT、SQLite 路径和初始管理员。
生产环境请修改 `jwt.secret` 和默认管理员密码。
