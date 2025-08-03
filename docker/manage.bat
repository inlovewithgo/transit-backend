@echo off
REM Simple Transit Backend Management Script for Windows

setlocal enabledelayedexpansion

REM Function to start all services
:start-all
echo [INFO] Starting all Transit Backend services...
echo [INFO] Starting PostgreSQL...
cd pg && docker-compose up -d && cd ..

echo [INFO] Starting Redis...
cd redis && docker-compose up -d && cd ..

echo [INFO] Starting Kafka...
cd kafka && docker-compose up -d && cd ..

echo [INFO] Waiting for services to start...
timeout /t 30 /nobreak >nul

echo [INFO] Starting main application...
docker-compose up -d

echo [INFO] All services started!
echo [INFO] Application: http://localhost:8080
echo [INFO] PostgreSQL: localhost:5432
echo [INFO] Redis: localhost:6379
echo [INFO] Kafka: localhost:9092
goto :eof

REM Function to stop all services
:stop-all
echo [INFO] Stopping all services...
docker-compose down
cd pg && docker-compose down && cd ..
cd redis && docker-compose down && cd ..
cd kafka && docker-compose down && cd ..
echo [INFO] All services stopped!
goto :eof

REM Function to show status
:status
echo [INFO] Service Status:
echo.
echo Main Application:
docker-compose ps
echo.
echo PostgreSQL:
cd pg && docker-compose ps && cd ..
echo.
echo Redis:
cd redis && docker-compose ps && cd ..
echo.
echo Kafka:
cd kafka && docker-compose ps && cd ..
goto :eof

REM Function to show logs
:logs
set service=%2
if "%service%"=="" set service=transit-app

echo [INFO] Showing logs for: %service%
docker-compose logs -f %service%
goto :eof

REM Function to clean up
:clean
echo [INFO] Cleaning up...
call :stop-all
docker system prune -f
echo [INFO] Cleanup complete!
goto :eof

REM Function to show help
:help
echo Simple Transit Backend Management
echo.
echo Commands:
echo   start     Start all services
echo   stop      Stop all services  
echo   restart   Restart all services
echo   status    Show service status
echo   logs      Show application logs
echo   clean     Clean up everything
echo   help      Show this help
goto :eof

REM Main logic
if "%1"=="start" (
    call :start-all
) else if "%1"=="stop" (
    call :stop-all
) else if "%1"=="restart" (
    call :stop-all
    timeout /t 5 /nobreak >nul
    call :start-all
) else if "%1"=="status" (
    call :status
) else if "%1"=="logs" (
    call :logs %1 %2
) else if "%1"=="clean" (
    call :clean
) else (
    call :help
)
