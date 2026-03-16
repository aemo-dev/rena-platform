@echo off
setlocal

:: Get the directory of the script
set "ROOT_DIR=%~dp0"

echo ==================================================
echo                 - RENA PLATFORM -
echo ==================================================
echo.

:: Kill existing processes on 8080 and 8081 to avoid conflicts
netstat -aon | findstr :8080 | findstr LISTENING >nul 2>&1 && (for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8080 ^| findstr LISTENING') do taskkill /f /pid %%a >nul 2>&1)
netstat -aon | findstr :8081 | findstr LISTENING >nul 2>&1 && (for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8081 ^| findstr LISTENING') do taskkill /f /pid %%a >nul 2>&1)

:: Start Backend (in background of the same window)
echo [Rena] Starting Backend (Go) on :8080...
if exist "%ROOT_DIR%rena\backend\go.mod" (
    start /b cmd /c "cd /d %ROOT_DIR%rena\backend && go run main.go"
) else (
    echo [Rena] Backend not initialized yet. Skipping...
)

:: Start Frontend (in the current window)
echo [Rena] Starting Frontend (Vite) on :5173...
echo [Rena] Note: Backend logs will also appear here.
cd /d %ROOT_DIR%rena\frontend && npm run dev -- --host
