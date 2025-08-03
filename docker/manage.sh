#!/bin/bash

# Simple Transit Backend Management Script

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Start all services
start() {
    print_info "Starting all Transit Backend services..."
    
    print_info "Starting PostgreSQL..."
    cd pg && docker-compose up -d && cd ..
    
    print_info "Starting Redis..."
    cd redis && docker-compose up -d && cd ..
    
    print_info "Starting Kafka..."
    cd kafka && docker-compose up -d && cd ..
    
    print_info "Waiting for services to start..."
    sleep 30
    
    print_info "Starting main application..."
    docker-compose up -d
    
    print_info "All services started!"
    print_info "Application: http://localhost:8080"
    print_info "PostgreSQL: localhost:5432"
    print_info "Redis: localhost:6379"
    print_info "Kafka: localhost:9092"
}

# Stop all services
stop() {
    print_info "Stopping all services..."
    docker-compose down || true
    cd pg && docker-compose down && cd .. || true
    cd redis && docker-compose down && cd .. || true
    cd kafka && docker-compose down && cd .. || true
    print_info "All services stopped!"
}

# Show status
status() {
    print_info "Service Status:"
    echo
    
    echo "Main Application:"
    docker-compose ps
    echo
    
    echo "PostgreSQL:"
    cd pg && docker-compose ps && cd ..
    echo
    
    echo "Redis:"
    cd redis && docker-compose ps && cd ..
    echo
    
    echo "Kafka:"
    cd kafka && docker-compose ps && cd ..
}

# Show logs
logs() {
    service=${1:-transit-app}
    print_info "Showing logs for: $service"
    docker-compose logs -f $service
}

# Clean up
clean() {
    print_info "Cleaning up..."
    stop
    docker system prune -f
    print_info "Cleanup complete!"
}

# Help
help() {
    echo "Simple Transit Backend Management"
    echo
    echo "Commands:"
    echo "  start     Start all services"
    echo "  stop      Stop all services"
    echo "  restart   Restart all services"
    echo "  status    Show service status"
    echo "  logs      Show application logs"
    echo "  clean     Clean up everything"
    echo "  help      Show this help"
}

# Main logic
case "${1:-help}" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        stop
        sleep 5
        start
        ;;
    status)
        status
        ;;
    logs)
        logs $2
        ;;
    clean)
        clean
        ;;
    help|*)
        help
        ;;
esac
