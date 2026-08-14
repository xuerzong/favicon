# Favicon

一个轻量级⚡️ API，用于获取任意网站的 favicon。

> **语言切换：** [English](README.md) | [简体中文](README.zh-CN.md)

- [在线演示](https://xuerzong.github.io/favicon/)
- [Demo API](https://favicon.xuco.me/)

## 使用方法

可选 [imgproxy](https://github.com/imgproxy/imgproxy) 参数（配合 `GET /{siteUrl}`）：

- `size` — 缩放尺寸
- `format` — 输出格式
- `quality` — 质量（0-100）
- `rotate` — 旋转角度（90 的倍数）

示例：

```
GET /example.com
GET /example.com?size=32&format=png
GET /example.com?size=32&format=png&rotate=180
GET /example.com?size=32&format=png&rotate=180&quality=100
```

## 本地开发

### API

```bash
cp .env.example .env
go run .            # 或 scripts/dev.sh
```

### Web

```bash
cd web
npm install
npm run dev         # 默认 API 地址 http://localhost:8080
```

前端 API 地址通过 Vite env 配置（`web/.env.development` / `web/.env.production`）：

```
VITE_API_BASE_URL=http://localhost:8080/
```

## 部署

```bash
docker compose -f docker/docker-compose.yml up -d
```

- favicon 服务: `:8866`
- imgproxy: `:8867`

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
