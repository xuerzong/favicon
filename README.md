# Favicon

A lightweight⚡️ api to get favicons from any website.

> **Languages:** [English](README.md) | [简体中文](README.zh-CN.md)

- [Playground](https://xuerzong.github.io/favicon/)
- [Demo API](https://favicon.xuco.me/)

## How to use

Optional [imgproxy](https://github.com/imgproxy/imgproxy) params (with `GET /{siteUrl}`):

- `size` — resize dimension
- `format` — output format
- `quality` — quality (0-100)
- `rotate` — rotation (multiples of 90)

Examples:

```
GET /example.com
GET /example.com?size=32&format=png
GET /example.com?size=32&format=png&rotate=180
GET /example.com?size=32&format=png&rotate=180&quality=100
```

## How to dev

### API

```bash
cp .env.example .env
go run .            # or scripts/dev.sh
```

### Web

```bash
cd web
npm install
npm run dev         # default API at http://localhost:8080
```

The frontend API base URL is configured via Vite env files (`web/.env.development` / `web/.env.production`):

```
VITE_API_BASE_URL=http://localhost:8080/
```

## How to deploy

```bash
docker compose -f docker/docker-compose.yml up -d
```

- favicon service: `:8866`
- imgproxy: `:8867`

## Environment Variables

| Variable              | Default                 | Description                                  |
| --------------------- | ----------------------- | -------------------------------------------- |
| `ADDRESS`             | `:8080`                 | Listen address                               |
| `SHUTDOWN_TIMEOUT`    | `10`                    | Graceful shutdown timeout (seconds)          |
| `GIN_MODE`            | `debug`                 | Gin run mode                                 |
| `IMAGE_SAVE_PATH`     | `images`                | Local image cache directory                  |
| `LOCAL_CACHE`         | `true`                  | Enable local disk cache                      |
| `CACHE_TTL`           | `1h`                    | Cache TTL                                    |
| `APP_BASE_URL`        | `http://127.0.0.1:8080` | Public service URL                           |
| `IMGPROXY_URL`        | empty                   | imgproxy URL; enables imgproxy path when set |
| `IMGPROXY_SOURCE_URL` | empty                   | imgproxy upstream source URL                 |
