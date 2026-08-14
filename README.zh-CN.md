# Favicon

一个简单的 Favicon API 服务：给定任意网站域名，返回其 favicon 图标，支持本地缓存与 imgproxy 图片处理。

> **语言切换：** [English](README.md) | [简体中文](README.zh-CN.md)

- 前端演示: `web/`
- 在线演示: https://xuerzong.github.io/favicon/
- API: https://favicon.xuco.me/

## 架构

```
favicon/
├── main.go            # HTTP 服务入口
├── server.go          # 优雅停机
├── internal/
│   ├── cache/         # 内存 TTL 缓存
│   ├── config/        # 环境变量配置
│   └── service/       # favicon 抓取、清理、imgproxy
├── pkg/
│   ├── core/          # favicon 获取逻辑
│   └── util/          # 域名解析、Content-Type
├── docker/            # Docker 部署
└── web/               # Vite + TS 前端（Svelte 模板）
```

## API

| 路由 | 说明 |
| --- | --- |
| `GET /{siteUrl}` | 返回站点 favicon（配置后默认经 imgproxy 处理） |
| `GET /raw/{siteUrl}` | 返回原始 favicon |
| `GET /` | 返回默认占位图 |

imgproxy 可选参数（配合 `GET /{siteUrl}`）：

- `size` — 缩放尺寸
- `format` — 输出格式
- `quality` — 质量（0-100）
- `rotate` — 旋转角度（90 的倍数）

示例：

```
GET /example.com
GET /example.com?size=32&format=png
GET /raw/example.com
```

## 本地开发

### 后端

```bash
cp .env.example .env
go run .            # 或 scripts/dev.sh
```

### 前端

```bash
cd web
npm install
npm run dev         # 默认 API 地址 http://localhost:8080
```

前端 API 地址通过 Vite env 配置（`web/.env.development` / `web/.env.production`）：

```
VITE_API_BASE_URL=http://localhost:8080/
```

## Docker 部署

```bash
docker compose -f docker/docker-compose.yml up -d
```

- favicon 服务: `:8866`
- imgproxy: `:8867`

### Dokploy

见 `docker/docker-compose.dokploy.yml`。

## GitHub Pages 部署

前端通过 GitHub Actions 自动部署（`.github/workflows/deploy.yml`），推送到 `main` 即触发：

1. 构建 `web/`（`base` 在 CI 中自动设为 `/favicon/`）
2. 上传产物并发布到 GitHub Pages

首次需要手动启用：仓库 **Settings → Pages → Source → GitHub Actions**。

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `ADDRESS` | `:8080` | 监听地址 |
| `SHUTDOWN_TIMEOUT` | `10` | 优雅停机超时（秒） |
| `GIN_MODE` | `debug` | Gin 运行模式 |
| `IMAGE_SAVE_PATH` | `images` | 本地图片缓存目录 |
| `LOCAL_CACHE` | `true` | 是否启用本地磁盘缓存 |
| `CACHE_TTL` | `1h` | 缓存 TTL |
| `APP_BASE_URL` | `http://127.0.0.1:8080` | 服务对外地址 |
| `IMGPROXY_URL` | 空 | imgproxy 地址，设置后走 imgproxy |
| `IMGPROXY_SOURCE_URL` | 空 | imgproxy 上游源地址 |

## License

MIT
