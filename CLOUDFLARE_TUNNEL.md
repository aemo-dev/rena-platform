# Cloudflare Tunnel Setup Script
# Run this to expose your local Go backend via Cloudflare Tunnel

# 1. Install cloudflared (if not already installed)
# winget install cloudflare.cloudflared

# 2. Login to Cloudflare
# cloudflared tunnel login

# 3. Create a tunnel (one-time)
# cloudflared tunnel create rena-backend

# 4. Configure tunnel to point to your Go backend (default: localhost:8080)
# cloudflared tunnel route dns rena-backend rena.yourdomain.com
# cloudflared tunnel run --url http://localhost:8080 rena-backend

# For quick testing without a domain:
# cloudflared tunnel --url http://localhost:8080
