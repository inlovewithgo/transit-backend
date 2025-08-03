# Transit Backend - Simple Setup Guide

## Quick Setup on Your PC

### Prerequisites
- Docker Desktop installed and running
- Git (optional, for cloning)

### Step 1: Clone & Navigate
```bash
git clone https://github.com/inlovewithgo/transit-backend.git
cd transit-backend/docker
```

### Step 2: Start Everything (Windows)
```cmd
manage.bat start
```

### Step 2: Start Everything (Linux/Mac)
```bash
chmod +x manage.sh
./manage.sh start
```

### Step 3: Test
Visit: http://localhost:8080/health

## What's Running?

### Services (All Separate)
- **PostgreSQL**: localhost:5432 (in docker/pg/)
- **Redis**: localhost:6379 (in docker/redis/)  
- **Kafka**: localhost:9092 (in docker/kafka/)
- **Your App**: localhost:8080 (main docker-compose.yaml)

### File Structure
```
docker/
├── docker-compose.yaml     # Main app (production)
├── docker-compose.dev.yaml # Development version
├── Dockerfile              # App container
├── manage.bat              # Windows script
├── manage.sh               # Linux/Mac script
├── pg/
│   └── docker-compose.yaml # PostgreSQL only
├── redis/
│   └── docker-compose.yaml # Redis only
└── kafka/
    └── docker-compose.yaml # Kafka + Zookeeper
```

## Simple Commands

### Windows (manage.bat)
```cmd
manage.bat start    # Start everything
manage.bat stop     # Stop everything  
manage.bat status   # Check status
manage.bat logs     # View app logs
manage.bat clean    # Clean up
```

### Linux/Mac (manage.sh)
```bash
./manage.sh start    # Start everything
./manage.sh stop     # Stop everything
./manage.sh status   # Check status  
./manage.sh logs     # View app logs
./manage.sh clean    # Clean up
```

## Development vs Production

### Development (docker-compose.dev.yaml)
- Code hot reload
- Connected to localhost services
- Port 3030

### Production (docker-compose.yaml)
- Optimized build
- Connected to localhost services  
- Port 8080

## Jenkins Setup

### Simple Jenkins Pipeline
The Jenkinsfile is now much simpler:
1. Test your Go code
2. Build binary
3. Build Docker image  
4. Deploy to dev/production

### Jenkins Requirements
- Go installed
- Docker installed
- That's it!

## Manual Setup (No Scripts)

### 1. Start Infrastructure
```bash
cd docker/pg && docker-compose up -d
cd ../redis && docker-compose up -d  
cd ../kafka && docker-compose up -d
```

### 2. Start Your App
```bash
cd .. && docker-compose up -d
```

### 3. Stop Everything
```bash
docker-compose down
cd pg && docker-compose down
cd ../redis && docker-compose down
cd ../kafka && docker-compose down
```

## Troubleshooting

### "Port already in use"
```bash
# Check what's using the port
netstat -ano | findstr :5432
netstat -ano | findstr :6379
netstat -ano | findstr :9092
netstat -ano | findstr :8080
```

### "Service won't start"
```bash
manage.bat logs     # Check app logs
docker logs transit-postgres  # Check postgres
docker logs transit-redis     # Check redis
docker logs transit-kafka     # Check kafka
```

### "Clean everything"
```bash
manage.bat clean
# Or manually:
docker system prune -a
docker volume prune
```

## Environment Variables

Your app reads these from docker-compose.yaml:
```yaml
- PORT=3030
- DB_HOST=localhost     # Points to your PC's localhost
- DB_PORT=5432
- DB_USER=postgres
- DB_PASSWORD=postgres
- DB_NAME=transit_wallet
- REDIS_HOST=localhost  # Points to your PC's localhost  
- REDIS_PORT=6379
- KAFKA_BROKERS=localhost:9092  # Points to your PC's localhost
```

## That's It!
- Each service runs separately like before
- Your app connects to localhost (your PC)
- Simple management scripts
- Simple Jenkins pipeline
- Everything works on your local PC

No complex orchestration, no huge config files, just simple Docker containers that work!
