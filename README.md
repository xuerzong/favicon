# Favicon

A simple Favicon API service: given any website domain, return its favicon, with local caching and imgproxy image processing.

> **Languages:** [English](README.md) | [简体中文](README.zh-CN.md)

- Web demo: `web/`
- Online demo: https://xuerzong.github.io/favicon/
- API: https://favicon.xuco.me/

## Architecture

```
favicon/
├── main.go            # HTTP server entry
├── server.go          # Graceful shutdown
├── internal/
│   ├── cache/         # In-memory TTL cache
│   ├── config/        # Environment config
│   └── service/       # Favicon fetching, cleanup, imgproxy
├── pkg/
│   ├── core/          # Favicon retrieval logic
│   └── util/          # Domain parsing, Content-Type
├── docker/            # Docker deployment
└── web/               # Vite + TS frontend (Svelte template)
```

## API

| Route | Description |
| --- | --- |
| `GET /{siteUrl}` | Return the site favicon (via imgproxy when configured) |
| `GET /raw/{siteUrl}` | Return the raw favicon |
| `GET /` | Return the default placeholder image |

Optional imgproxy params (with `GET /{siteUrl}`):

- `size` — resize dimension
- `format` — output format
- `quality` — quality (0-100)
- `rotate` — rotation (multiples of 90)

Examples:

```
GET /example.com
GET /example.com?size=32&format=png
GET /raw/example.com
```

## Local Development

### Backend

```bash
cp .env.example .env
go run .            # or scripts/dev.sh
```

### Frontend

```bash
cd web
npm install
npm run dev         # default API at http://localhost:8080
```

The frontend API base URL is configured via Vite env files (`web/.env.development` / `web/.env.production`):

```
VITE_API_BASE_URL=http://localhost:8080/
```

## Docker Deployment

```bash
docker compose -f docker/docker-compose.yml up -d
```

- favicon service: `:8866`
- imgproxy: `:8867`

### Dokploy

See `docker/docker-compose.dokploy.yml`.

## GitHub Pages Deployment

The frontend is deployed automatically via GitHub Actions (`.github/workflows/deploy.yml`) on every push to `main`:

1. Build `web/` (the `base` is set to `/favicon/` automatically in CI)
2. Upload the artifact and publish to GitHub Pages

First-time setup: enable Pages manually at **Settings → Pages → Source → GitHub Actions**.

## Environment Variables

| Variable | Default | Description |
| --- | --- | --- |
| `ADDRESS` | `:8080` | Listen address |
| `SHUTDOWN_TIMEOUT` | `10` | Graceful shutdown timeout (seconds) |
| `GIN_MODE` | `debug` | Gin run mode |
| `IMAGE_SAVE_PATH` | `images` | Local image cache directory |
| `LOCAL_CACHE` | `true` | Enable local disk cache |
| `CACHE_TTL` | `1h` | Cache TTL |
| `APP_BASE_URL` | `http://127.0.0.1:8080` | Public service URL |
| `IMGPROXY_URL` | empty | imgproxy URL; enables imgproxy path when set |
| `IMGPROXY_SOURCE_URL` | empty | imgproxy upstream source URL |

## License

MIT
