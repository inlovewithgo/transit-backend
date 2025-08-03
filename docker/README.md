# Transit Backend Docker Setup

This directory contains all the Docker configuration files for the Transit Backend application.

## Quick Start

### Prerequisites
- Docker Engine 20.10+
- Docker Compose 2.0+
- Git

### Development Environment

```bash
git clone https://github.com/inlovewithgo/transit-backend.git
cd transit-backend/docker

./manage.sh dev
manage.bat dev
```

### Production Environment

```bash
./manage.sh build
./manage.sh start

manage.bat build
manage.bat start
```

## Files Overview

### Docker Files
- `Dockerfile` - Multi-stage Docker build for the Go application
- `docker-compose.yaml` - Production environment with all services
- `docker-compose.dev.yaml` - Development environment (simplified)
- `.dockerignore` - Files to ignore during Docker build

### Management Scripts
- `manage.sh` - Linux/macOS management script
- `manage.bat` - Windows management script

### Configuration
- `monitoring/prometheus.yml` - Prometheus monitoring configuration

## Services

The production environment includes the following services:

### Core Application
- **transit-app**: Main Go application (Port 3030)

### Infrastructure
- **postgres**: PostgreSQL database (Port 5432)
- **redis**: Redis cache (Port 6379)
- **kafka**: Apache Kafka message broker (Port 9092)
- **zookeeper**: Kafka dependency (Port 2181)

### Monitoring
- **prometheus**: Metrics collection (Port 9090)
- **grafana**: Metrics visualization (Port 3000)

## Environment Variables

### Application Configuration
```bash
PORT=3030                  
GO_ENV=production         
```

### Database Configuration
```bash
DB_HOST=postgres      
DB_PORT=5432          
DB_USER=postgres       
DB_PASSWORD=postgres      
DB_NAME=transit_wallet   
```

### Redis Configuration
```bash
REDIS_HOST=redis   
REDIS_PORT=6379   
```

### Kafka Configuration
```bash
KAFKA_BROKERS=kafka:9092
```

## Usage Examples

### Management Commands

```bash
./manage.sh build
./manage.sh start
./manage.sh dev
./manage.sh logs transit-app
./manage.sh logs postgres
./manage.sh status
./manage.sh stop
./manage.sh clean
./manage.sh backup
./manage.sh restore backup_20231201_120000.sql
```

### Manual Docker Commands

```bash
docker build -t transit-backend:latest -f Dockerfile ..
docker-compose up -d
docker-compose -f docker-compose.dev.yaml up -d
docker-compose logs -f transit-app
docker-compose down

docker-compose down -v --remove-orphans
docker system prune -f
```

## Development Workflow

### 1. Development Environment
```bash
./manage.sh dev
```

### 2. Testing
```bash
docker-compose -f docker-compose.dev.yaml exec transit-app-dev go test ./...
docker-compose -f docker-compose.dev.yaml exec transit-app-dev go test -coverprofile=coverage.out ./...
```

### 3. Debugging
```bash
docker-compose -f docker-compose.dev.yaml exec transit-app-dev sh
./manage.sh logs transit-app-dev
dlv attach --listen=:2345 --headless=true --api-version=2 --accept-multiclient <pid>
```

## Production Deployment

### 1. Build and Test
```bash
./manage.sh build
docker run --rm -p 3030:3030 transit-backend:latest
```

### 2. Deploy
```bash
./manage.sh start
curl http://localhost:3030/health
```

### 3. Monitor
- **Application Health**: http://localhost:3030/health
- **Metrics**: http://localhost:9090 (Prometheus)
- **Dashboards**: http://localhost:3000 (Grafana - admin/admin)

## Troubleshooting

### Common Issues

1. **Port conflicts**
   ```bash
   # Check what's using the port
   netstat -tulpn | grep :3030
   
   # Stop conflicting services or change ports in docker-compose.yaml
   ```

2. **Database connection issues**
   ```bash
   # Check database status
   ./manage.sh logs postgres
   
   # Verify database is ready
   docker-compose exec postgres pg_isready -U postgres
   ```

3. **Out of disk space**
   ```bash
   # Clean up Docker resources
   ./manage.sh clean
   
   # More aggressive cleanup
   docker system prune -a
   ```

4. **Application won't start**
   ```bash
   # Check application logs
   ./manage.sh logs transit-app
   
   # Check if all dependencies are ready
   ./manage.sh status
   ```

### Health Checks

```bash
curl http://localhost:3030/health

docker-compose exec postgres pg_isready -U postgres

docker-compose exec redis redis-cli ping

docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --list
```

## Security Considerations

1. **Environment Variables**: Use Docker secrets or external secret management in production
2. **Network Security**: Configure proper firewall rules
3. **Image Security**: Regularly scan images for vulnerabilities
4. **Database Security**: Use strong passwords and enable SSL
5. **Access Control**: Implement proper authentication and authorization

## Performance Tuning

### Resource Limits
Add resource limits to docker-compose.yaml:
```yaml
services:
  transit-app:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '1.0'
          memory: 512M
```

### Database Optimization
```yaml
postgres:
  environment:
    - POSTGRES_SHARED_PRELOAD_LIBRARIES=pg_stat_statements
    - POSTGRES_MAX_CONNECTIONS=200
```

## Monitoring and Logging

### Prometheus Metrics
- Application metrics: http://localhost:3030/metrics
- Prometheus UI: http://localhost:9090

### Grafana Dashboards
- Grafana UI: http://localhost:3000 (admin/admin)
- Import dashboard ID: 1860 for Node Exporter

### Log Aggregation
For production, consider using:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Fluentd + Elasticsearch
- Loki + Grafana

## Backup and Recovery

### Database Backup
```bash
./manage.sh backup

docker-compose exec postgres pg_dump -U postgres transit_wallet > backup.sql
```

### Volume Backup
```bash
docker run --rm -v transit-backend_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz -C /data .
```

## CI/CD Integration

This setup integrates with Jenkins pipeline (see `../Jenkinsfile`) which:
1. Builds and tests the application
2. Creates Docker images
3. Runs security scans
4. Deploys to staging/production
5. Performs health checks

## Contributing

1. Make changes to the application code
2. Test locally using development environment
3. Update Docker configurations if needed
4. Test production build
5. Submit pull request

The Jenkins pipeline will automatically test and deploy approved changes.
