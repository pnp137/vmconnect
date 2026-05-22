# ✅ Backend Setup Checklist (Frontend Team)

Use this checklist to quickly set up the backend.

## Prerequisites ✓
- [ ] Docker Desktop installed ([Download](https://www.docker.com/products/docker-desktop))
- [ ] Backend code cloned: `vmconnect/` folder
- [ ] Terminal open in the project root

---

## Startup (Do This First Time)

- [ ] **Step 1:** Open terminal in `vmconnect/` folder
- [ ] **Step 2:** Run this command:
  ```bash
  docker compose -f build/docker-compose.yml up -d --build
  ```
- [ ] **Step 3:** Wait 10 seconds for MySQL to start
- [ ] **Step 4:** Verify it works:
  ```bash
  curl http://localhost:8080/health
  ```
- [ ] **Step 5:** Should see JSON response (green ✅)

---

## Using the Backend

- [ ] API is available at: `http://localhost:8080`
- [ ] Database: `localhost:3306` (vmconnect_user / vmconnect_pass)
- [ ] Test login endpoint:
  ```bash
  curl -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"merchant@test.com","password":"merchant123"}'
  ```

---

## Loading Test Data (Optional)

- [ ] Run this command to load sample data:
  ```bash
  cat database_seed.sql | docker exec -i vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect
  ```
- [ ] Test with:
  - Email: `merchant@test.com`
  - Password: `merchant123`

---

## Daily Usage

### ⏰ When You Start Work
```bash
docker compose -f build/docker-compose.yml up -d --build
```

### 🛑 When You're Done
```bash
docker compose -f build/docker-compose.yml down
```

### 🔄 If Something Breaks
```bash
docker compose -f build/docker-compose.yml down
docker compose -f build/docker-compose.yml up -d --build
```

### 📊 View Live Logs
```bash
docker logs vmconnect-api -f
```

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| Docker command not found | Install Docker Desktop |
| Backend not responding | Wait 10 seconds, MySQL starting up |
| Port 8080 in use | `docker compose -f build/docker-compose.yml down` and try again |
| Cannot login | Make sure seed data loaded |

---

## Database Access (Advanced)

**View tables:**
```bash
docker exec vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect -e "SHOW TABLES;"
```

**View users:**
```bash
docker exec vmconnect-db mysql -u vmconnect_user -pvmconnect_pass vmconnect -e "SELECT * FROM users;"
```

---

## Quick Reference

```
🌐 API URL:        http://localhost:8080
🗄️  DB Host:       localhost
🔐 DB User:        vmconnect_user
🔑 DB Password:    vmconnect_pass
📦 DB Name:        vmconnect
```

---

## Notes

- ✅ No need to install Go, MySQL, or any dependencies
- ✅ Everything runs in Docker (clean environment)
- ✅ Database data persists between restarts
- ✅ Easy to reset: just `docker compose -f build/docker-compose.yml down` and up again
- ✅ Share this guide with other frontend team members!

---

**Questions? Check README.md or ask the backend team 🙂**
