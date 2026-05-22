# Backend Setup Guide for Frontend Team

Welcome! This guide will help you set up and run the backend locally using Docker.

## Prerequisites
- **Docker Desktop** ([Download here](https://www.docker.com/products/docker-desktop))
- That's all! No need to install Go, MySQL, Node, or anything else locally.

---

## ⚡ Quick Start (30 seconds)

```bash
cd vmconnect
docker compose -f build/docker-compose.yml up -d --build
```

**Done!** Your backend is running at `http://localhost:8080` ✅

---

## 🚀 Full Setup Guide

### 1️⃣ Start the Backend
```bash
cd vmconnect
docker compose -f build/docker-compose.yml up -d --build
```

**What this does:**
- Downloads and starts MySQL database
- Builds and runs the backend API
- Sets up network and volumes
- Runs database migrations automatically

**That's it!** The backend is now running.

### 2️⃣ Verify It's Working
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "component": "api",
  "hostname": "...",
  "service": "vmconnect-api",
  "services": {}
}
```

### 3️⃣ Stop the Backend (When Done)
```bash
cd vmconnect
docker compose -f build/docker-compose.yml down
```

---

## 🎯 Convenient Scripts

### Start Backend Quickly
```bash
# Option 1: Manual
cd vmconnect
docker compose -f build/docker-compose.yml up -d --build

# Option 2: Using the provided script
./scripts/start-backend.sh
```

### Stop Backend
```bash
# Option 1: Manual
cd vmconnect
docker compose -f build/docker-compose.yml down

# Option 2: Using the provided script
./scripts/stop-backend.sh
```

### View Live Logs
```bash
# See all container logs in real-time
./scripts/logs.sh

# Or see specific service:
./scripts/logs.sh api      # API logs only
./scripts/logs.sh db       # Database logs only
./scripts/logs.sh status   # Container status
```

### Restart Backend
```bash
./scripts/restart-backend.sh
```

### Check Status
```bash
docker ps | grep vmconnect
```

You should see:
```
vmconnect-api    (running on port 8080)
vmconnect-db     (running on port 3306)
```

---

## 🔌 Using the Backend in Your Frontend

### API Base URL
```
http://localhost:8080
```

### Example: Login Request
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"merchant@test.com","password":"merchant123"}'
```

### JavaScript/Fetch Example
```javascript
const API_URL = 'http://localhost:8080';

async function login(email, password) {
  const response = await fetch(`${API_URL}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  return response.json();
}
```

---

## 💾 Database

### Database Credentials
```
Host:     localhost
Port:     3306
Database: vmconnect
User:     vmconnect_user
Password: vmconnect_pass
```

### View Database
```bash
# Using Docker (easiest)
docker exec -it vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect

# Or if MySQL CLI installed locally
mysql -h localhost -u vmconnect_user -pvmconnect_pass vmconnect
```

## 🗄️ Database Info

**Host:** localhost  
**Port:** 3306  
**Database:** vmconnect  
**User:** vmconnect_user  
**Password:** vmconnect_pass

### View Database (Optional)
```bash
# Install MySQL CLI (if you have it)
mysql -h localhost -u vmconnect_user -pvmconnect_pass vmconnect

# Or using Docker
docker exec -it vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect
```

---

## 🚨 Troubleshooting

### Backend not responding?
```bash
# Check if it's running
docker ps

# You should see:
# - vmconnect-api (running)
# - vmconnect-db (running)
```

### Port 8080 in use?
```bash
# Either stop other services or use a different port
# Edit build/docker-compose.yml and change 8080 to 8081
```

### Database error?
```bash
# Check the logs
docker logs vmconnect-db
docker logs vmconnect-api
```

---

## 📚 API Endpoints

Common endpoints you'll use:

```
POST   /api/auth/login                    - Login
POST   /api/merchant/register             - Register as Merchant
POST   /api/vendor/register               - Register as Vendor
GET    /api/merchant/profile              - Get merchant details
GET    /api/vendor/profile                - Get vendor details
GET    /api/products                      - List products
POST   /api/cart/add                      - Add to cart
POST   /api/orders                        - Create order
```

See full API documentation in `bruno/` folder (API testing collection).

---

## 🐛 Common Issues

| Problem | Solution |
|---------|----------|
| `docker not found` | Install Docker Desktop |
| `Port 8080 in use` | Stop other services or change port in build/docker-compose.yml |
| `Cannot connect to database` | Wait 10 seconds and try again (MySQL takes time to start) |
| `API returning 500 errors` | Check logs: `docker logs vmconnect-api` |

---

## 📝 Tips

- Keep the backend running in one terminal while developing frontend
- Check logs with: `docker logs vmconnect-api -f` (live updates)
- Restart backend quickly: `docker compose -f build/docker-compose.yml up -d --build`
- Database data persists between restarts (won't lose data)

---

## Need Help?

1. Check the logs: `docker logs vmconnect-api`
2. Make sure Docker is running
3. Try restarting: `docker compose -f build/docker-compose.yml down && docker compose -f build/docker-compose.yml up -d --build`
4. Ask the backend team 🙂

---

**Happy coding! 🚀**
