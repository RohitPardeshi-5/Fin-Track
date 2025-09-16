@echo off
echo 🚀 FinTrack - Starting All Services
echo Reading ports from .env file...
echo.

REM Kill existing processes
taskkill /IM go.exe /F 2>nul

echo Starting services...
start "Ai Service" cmd /k "go run simple-ai-service.go"
timeout /t 3 >nul

start "User Service" cmd /k "go run simple-user-service.go"
timeout /t 3 >nul

start "Expense Service" cmd /k "go run simple-expense-service.go"  
timeout /t 3 >nul

start "Report Service" cmd /k "go run simple-report-service.go"
timeout /t 3 >nul

start "Web Frontend" cmd /k "go run cmd/web-frontend/main.go"
 
echo.
echo ✅ All services started!
echo 📝 Edit .env file to change ports
echo 🌐 Access: http://localhost:8000
echo.
pause