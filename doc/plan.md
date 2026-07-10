# AIPanel

> Lightweight AI-first Server Panel

AIPanel 是一个面向开发者和个人服务器用户的轻量级 AI 运维面板。

项目目标并不是替代宝塔，而是提供：

* 服务器监控
* Docker 管理
* Web Terminal
* 日志分析
* Workflow 自动化
* 任意 AI API 接入

核心理念：

```text
AI Suggestion + Human Confirmation + Workflow Execution
```

避免传统 Agent 的不可控执行问题。

---

# 项目定位

## 为什么做 AIPanel

现有方案存在以下问题：

### 宝塔

* 功能过重
* 占用资源较高
* 偏传统网站运维

### Portainer

* 仅关注 Docker
* 缺少 AI 能力

### CasaOS

* 偏家庭 NAS 场景

### 各类 Agent

* 容易出现误操作
* 调试困难
* Token 消耗高
* 执行不可预测

AIPanel 希望成为：

> 一个轻量、可扩展、支持任意大模型的 AI 运维助手。

---

# 核心功能

## Dashboard

实时展示：

* CPU
* Memory
* Disk
* Network
* Docker Container Count
* System Load

支持：

* WebSocket 实时刷新
* 多服务器监控（后续）

---

## Docker 管理

功能：

* 容器列表
* 查看日志
* 启动 / 停止 / 重启
* 进入容器终端
* 镜像管理
* Docker Compose 执行

后端直接调用 Docker SDK。

---

## Web Terminal

提供浏览器终端：

支持：

* bash
* zsh
* powershell

技术方案：

```text
xterm.js
      ↓
WebSocket
      ↓
PTY
```

---

## 日志分析

支持分析：

* Docker Logs
* journalctl
* Nginx Logs
* Application Logs

示例：

```text
为什么网站返回 502？
```

AI 自动分析：

```text
检测结果：

1. nginx upstream 不可达
2. docker 容器已退出
3. 80端口未监听
```

---

## AI Chat

支持任意兼容 OpenAI 的模型。

支持：

* OpenAI
* Claude
* Gemini
* OpenRouter
* DeepSeek
* Qwen
* GLM
* Ollama

统一配置：

```yaml
provider:
base_url:
api_key:
model:
```

统一调用接口：

```http
POST /api/chat
```

---

# Workflow 系统

AIPanel 不直接让 AI 执行命令。

执行流程：

```text
用户问题
    ↓
AI生成执行计划
    ↓
展示执行步骤
    ↓
用户确认
    ↓
执行 Workflow
```

---

## Workflow 示例

### 部署项目

```yaml
steps:
  - git clone
  - docker compose pull
  - docker compose up -d
```

---

### 清理磁盘

```yaml
steps:
  - docker image prune -f
  - docker volume prune -f
```

---

### 数据库备份

```yaml
steps:
  - mysqldump
  - upload to R2
```

---

# 插件系统

目录：

```text
plugins/

docker/
nginx/
mysql/
redis/
cloudflare/
```

每个插件提供：

```yaml
name:
tools:
workflows:
permissions:
```

示例：

```yaml
name: docker

tools:
  - logs
  - restart
  - exec

workflows:
  - cleanup
  - update
```

---

# 系统架构

```text
Browser
    ↓
Frontend (Vue3)
    ↓
API Server (Go)
    ↓
─────────────────────
Docker
System Metrics
PTY
Workflow Engine
AI Provider
Plugin Runtime
```

---

# 技术栈

## Frontend

* Vue3
* Element Plus
* Pinia
* Vue Router
* xterm.js
* Monaco Editor
* ECharts

---

## Backend

推荐：

* Go 1.25+

主要依赖：

* Gin / Fiber
* Docker SDK
* SQLite
* WebSocket
* PTY

---

## Database

```text
SQLite
```

无需 MySQL。

---

# 项目结构

```text
AIPanel/

├── frontend/
├── server/
├── plugins/
├── storage/
├── docs/
├── scripts/
├── config.yaml
└── docker-compose.yml
```

后端结构：

```text
server/internal/

├── ai/
├── auth/
├── docker/
├── system/
├── terminal/
├── workflow/
├── plugin/
└── websocket/
```

---

# API 设计

## 获取系统状态

```http
GET /api/system
```

---

## Docker 容器列表

```http
GET /api/docker/containers
```

---

## 获取日志

```http
GET /api/logs
```

---

## AI 对话

```http
POST /api/chat
```

---

## 执行 Workflow

```http
POST /api/workflows/run
```

---

## Web Terminal

```http
WS /api/terminal
```

---

# 数据库设计

## servers

```sql
id
name
host
type
status
```

---

## ai_configs

```sql
id
provider
base_url
api_key
model
```

---

## workflows

```sql
id
name
content
created_at
```

---

## workflow_logs

```sql
id
workflow_id
status
logs
```

---

## chat_history

```sql
id
role
content
created_at
```

---

# 部署方式

## Docker

```bash
docker compose up -d
```

默认端口：

```text
3000
```

---

# 第一阶段（MVP）

预计开发周期：

2 周。

## Week 1

* 登录系统
* Dashboard
* Docker 管理
* Web Terminal
* AI 配置

## Week 2

* AI Chat
* Workflow
* 日志分析
* 插件系统
* Docker 部署

---

# 第二阶段

支持：

* 多服务器管理
* SSH 节点接管
* Cloudflare 集成
* R2 备份
* MCP 集成
* 手机端适配

---

# 第三阶段

支持：

* Workflow Marketplace
* 插件市场
* 团队协作
* 告警通知
* 自动巡检
* Kubernetes 插件

---

# 项目原则

## 轻量

不做：

* FTP
* 邮件服务器
* LNMP 一键安装
* 网站文件管理器

只关注：

```text
Server + Docker + AI
```

---

## 安全

AI 默认不直接执行命令。

所有高风险操作：

```text
生成计划
↓
用户确认
↓
执行
```

---

## 可扩展

支持：

* 任意 AI Provider
* 任意 Workflow
* 任意 Plugin

---

# 愿景

打造一个：

> 面向开发者的轻量级 AI 运维平台。

让用户可以通过自然语言完成：

```text
分析问题
执行运维任务
部署项目
管理服务器
自动化工作流
```

而无需复杂的命令和脚本。
