# VMConnect Backend

Backend API for vendor and merchant apps.

## Frontend Setup

Use [FRONTEND_SETUP.md](FRONTEND_SETUP.md) to run the backend locally.

## Common Commands

```bash
docker compose -f build/docker-compose.yml up -d --build
```

```bash
curl http://localhost:8080/health
```

```bash
docker compose -f build/docker-compose.yml down
```

## API Base URL

```text
http://localhost:8080
```
