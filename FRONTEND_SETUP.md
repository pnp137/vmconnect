# Frontend Backend Setup

## Use These Scripts

Start backend:

```bash
cd vmconnect
./scripts/start-backend.sh
```

Stop backend:

```bash
./scripts/stop-backend.sh
```

Restart backend:

```bash
./scripts/restart-backend.sh
```

View API logs:

```bash
./scripts/logs.sh api
```

View database logs:

```bash
./scripts/logs.sh db
```

Check running containers:

```bash
./scripts/logs.sh status
```

## Check Backend

```bash
curl http://localhost:8080/health
```

## API URL

```text
http://localhost:8080
```

## Seed Data

```bash
cat database_seed.sql | docker exec -i vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect
```

## Database Login

```bash
docker exec -it vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect
```

## Docker Commands

Use these only if scripts are not available.

```bash
docker compose -f build/docker-compose.yml up -d --build
```

```bash
docker compose -f build/docker-compose.yml down
```
