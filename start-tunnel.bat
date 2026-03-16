@echo off
REM Cloudflare Tunnel + Backend Launcher
REM 
REM Prerequisites:
REM 1. Install cloudflared: winget install cloudflare.cloudflared
REM 2. Login: cloudflared tunnel login
REM 3. Create tunnel: cloudflared tunnel create rena-backend
REM 4. Configure DNS: cloudflared tunnel route dns rena-backend rena.yourdomain.com
REM 
REM Usage:
REM - Run locally: start-tunnel.bat
REM - Then update VITE_API_URL in .env to your tunnel URL

echo Starting Go backend...
start "Go Backend" cmd /k "cd rena\backend && go run main.go"

echo Waiting for backend to start...
timeout /t 5 /nobreak > nul

echo Starting Cloudflare Tunnel...
echo Once connected, your backend will be accessible via the tunnel URL
cloudflared tunnel run --url http://localhost:8080 rena-backend
