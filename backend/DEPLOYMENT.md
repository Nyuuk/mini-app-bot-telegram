# 🚀 Deployment Guide - Mini App Bot Telegram Backend

Panduan lengkap untuk mendeploy aplikasi backend Mini App Bot Telegram.

## 📋 Prerequisites

### System Requirements

- **Go**: Version 1.19 atau lebih baru
- **PostgreSQL**: Version 12 atau lebih baru
- **Git**: Untuk cloning repository
- **Docker** (Opsional): Untuk containerized deployment

### Environment Variables

Buat file `.env` di root directory dengan konfigurasi berikut:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres

# Environment
ENV=production

# Log Level (trace, debug, info, warn, error, fatal, panic)
LOG_LEVEL=debug

# JWT
JWT_SECRET_KEY=your_jwt_secret
JWT_EXPIRATION=24
```

## 🔧 Development Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd mini-app-bot-telegram/backend
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Setup Database

```sql
-- Buat database PostgreSQL
CREATE DATABASE mini_app_bot_telegram;
```

### 4. Setup Environment

```bash
cp .env.example .env
# Edit .env dengan konfigurasi yang sesuai
```

### 5. Run Database Migration

```bash
# Uncomment migration code di main.go (line 28-35)
go run main.go
```

### 6. Run Application

```bash
go run main.go
```

Server akan berjalan di `http://localhost:3000`

## 🐳 Docker Deployment

### 1. Docker Files Structure

Project sudah dilengkapi dengan:
- `backend/Dockerfile` - Multi-stage build untuk backend Go application
- `backend/.dockerignore` - File untuk mengexclude files yang tidak diperlukan
- `docker-compose.yml` - Development environment
- `docker-compose.prod.yml` - Production environment
- `backend/init.sql` - Database initialization script

### 2. Architecture Support

**Multi-Architecture Compatibility:**
- **AMD64**: Supported (Intel/AMD processors)
- **ARM64**: Supported (Apple M1/M2, ARM servers)
- **Auto-detection**: Docker automatically builds for host architecture

### 3. Development dengan Docker

```bash
# Clone repository
git clone <repository-url>
cd mini-app-bot-telegram

# Build and run development environment
docker-compose up --build -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down

# Clean up (remove volumes)
docker-compose down -v
```

### 4. Production Deployment

```bash
# Set environment variables
export POSTGRES_PASSWORD=your_secure_password
export JWT_SECRET_KEY=your_super_secure_jwt_key
export VERSION=v1.0.0

# Deploy production environment
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Stop production services
docker-compose -f docker-compose.prod.yml down
```

### 5. Environment Variables for Production

Buat file `.env.prod` untuk production:

```env
# Database Configuration
POSTGRES_DB=mini_app_bot_telegram_prod
POSTGRES_USER=app_user
POSTGRES_PASSWORD=your_super_secure_password
POSTGRES_PORT=5432

# Backend Configuration
BACKEND_PORT=3000
ENV=production
LOG_LEVEL=warn
JWT_SECRET_KEY=your_super_secure_jwt_key_for_production
JWT_EXPIRATION=24

# Docker Registry (for GitHub Actions)
DOCKER_REGISTRY=ghcr.io
GITHUB_REPOSITORY=your-username/mini-app-bot-telegram
VERSION=latest
```

## 🚀 GitHub Actions CI/CD

Project ini dilengkapi dengan automated CI/CD pipeline menggunakan GitHub Actions.

### 1. Workflow Structure

- **`.github/workflows/ci.yml`** - Continuous Integration untuk setiap push/PR
- **`.github/workflows/build-and-deploy.yml`** - Build dan deploy berdasarkan tag

### 2. CI/CD Features

#### Continuous Integration (CI)
- **Backend Testing**: Unit tests, linting, security scanning
- **Frontend Testing**: Unit tests, linting, build verification
- **Docker Build**: Test Docker image building
- **Code Quality**: Static analysis dengan golangci-lint dan gosec

#### Build and Deploy (CD)
- **Smart Detection**: Hanya build komponen yang berubah
- **Multi-platform Build**: Support untuk linux/amd64 dan linux/arm64
- **Container Registry**: Push ke GitHub Container Registry (ghcr.io)
- **Version Tagging**: Otomatis tag dengan format `v1.0.0-backend` atau `v1.0.0-frontend`

### 3. Deployment Process

#### Step 1: Create Tag
```bash
# Tag untuk backend changes
git tag v1.0.0
git push origin v1.0.0

# Tag untuk specific version
git tag v1.1.0-beta
git push origin v1.1.0-beta
```

#### Step 2: Automatic Detection
GitHub Actions akan:
1. Detect perubahan files antara tag sebelumnya dan current tag
2. Build hanya komponen yang berubah:
   - Jika ada changes di `backend/` → build backend image
   - Jika ada changes di `frontend/` → build frontend image
   - Jika keduanya berubah → build kedua images

#### Step 3: Image Tagging
Images akan di-tag dengan format:
- Backend: `ghcr.io/username/mini-app-bot-telegram:v1.0.0-backend`
- Frontend: `ghcr.io/username/mini-app-bot-telegram:v1.0.0-frontend`

### 4. Setup GitHub Actions

#### Required Secrets
Tidak ada secrets tambahan yang diperlukan. GitHub Actions menggunakan built-in `GITHUB_TOKEN`.

#### Repository Settings
1. Enable GitHub Container Registry:
   - Go to repository Settings → Actions → General
   - Enable "Read and write permissions" for GITHUB_TOKEN

2. Optional: Create Environment
   - Go to Settings → Environments
   - Create "production" environment
   - Add protection rules jika diperlukan

### 5. Manual Deployment

Jika ingin deploy manually:

```bash
# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Pull latest images
docker pull ghcr.io/username/mini-app-bot-telegram:v1.0.0-backend

# Run with docker-compose production
VERSION=v1.0.0 docker-compose -f docker-compose.prod.yml up -d
```

### 6. Monitoring Deployment

#### Check Workflow Status
- Go to repository → Actions tab
- Monitor workflow progress dan logs

#### Check Container Registry
- Go to repository → Packages tab  
- View published container images

#### Health Check
```bash
# Check application health
curl https://your-domain.com/health

# Check specific version
curl -H "User-Agent: deployment-check" https://your-domain.com/health
```

### 7. Rollback Process

```bash
# Rollback to previous version
export PREVIOUS_VERSION=v0.9.0
docker-compose -f docker-compose.prod.yml down
VERSION=$PREVIOUS_VERSION docker-compose -f docker-compose.prod.yml up -d

# Verify rollback
curl https://your-domain.com/health
```

### 8. Troubleshooting CI/CD

#### Build Failures
```bash
# Check workflow logs di GitHub Actions tab
# Common issues:
# 1. Go module issues → Check go.mod/go.sum
# 2. Docker build issues → Check Dockerfile
# 3. Test failures → Check test coverage
```

#### Deployment Issues
```bash
# Check container logs
docker-compose -f docker-compose.prod.yml logs backend

# Check image availability
docker images | grep mini-app-bot

# Verify environment variables
docker-compose -f docker-compose.prod.yml config
```

## ☁️ Production Deployment

### 1. VPS/Cloud Server Setup

#### Ubuntu/Debian Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Install Nginx (as reverse proxy)
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

### 2. Database Setup

```bash
# Switch to postgres user
sudo -u postgres psql

# Create database and user
CREATE DATABASE mini_app_bot_telegram;
CREATE USER app_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE mini_app_bot_telegram TO app_user;
\q
```

### 3. Application Deployment

```bash
# Clone repository
git clone <repository-url>
cd mini-app-bot-telegram/backend

# Setup environment
cp .env.example .env
nano .env  # Edit dengan konfigurasi production

# Build application
go build -o main .

# Create systemd service
sudo nano /etc/systemd/system/mini-app-bot.service
```

#### Systemd Service File

```ini
[Unit]
Description=Mini App Bot Telegram Backend
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/mini-app-bot-telegram/backend
ExecStart=/home/ubuntu/mini-app-bot-telegram/backend/main
Restart=always
RestartSec=10
Environment=PATH=/usr/local/go/bin:/usr/bin:/bin

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable mini-app-bot
sudo systemctl start mini-app-bot
sudo systemctl status mini-app-bot
```

### 4. Nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/mini-app-bot
```

```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Health check endpoint
    location /health {
        access_log off;
        proxy_pass http://localhost:3000/health;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/mini-app-bot /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 5. SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Get SSL certificate
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

## 🔧 Environment Configurations

### Development Environment

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mini_app_bot_telegram_dev
DB_USER=dev_user
DB_PASSWORD=dev_password
JWT_SECRET_KEY=dev_secret_key
PORT=3000
LOG_LEVEL=debug
AUTO_MIGRATE=true
```

### Staging Environment

```env
DB_HOST=staging-db.example.com
DB_PORT=5432
DB_NAME=mini_app_bot_telegram_staging
DB_USER=staging_user
DB_PASSWORD=staging_secure_password
JWT_SECRET_KEY=staging_secret_key
PORT=3000
LOG_LEVEL=info
AUTO_MIGRATE=false
```

### Production Environment

```env
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_NAME=mini_app_bot_telegram_prod
DB_USER=prod_user
DB_PASSWORD=super_secure_production_password
JWT_SECRET_KEY=super_secure_production_jwt_key
PORT=3000
LOG_LEVEL=warn
AUTO_MIGRATE=false
TIMEZONE=Asia/Jakarta
```

## 📊 Monitoring & Maintenance

### 1. Health Check

```bash
# Check application status
curl http://localhost:3000/health

# Check systemd service
sudo systemctl status mini-app-bot

# View logs
sudo journalctl -u mini-app-bot -f
```

### 2. Database Backup

```bash
# Create backup script
nano backup.sh
```

```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/home/ubuntu/backups"
DB_NAME="mini_app_bot_telegram"

mkdir -p $BACKUP_DIR

pg_dump -h localhost -U app_user -d $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# Keep only last 7 backups
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete

echo "Backup completed: backup_$DATE.sql"
```

```bash
chmod +x backup.sh

# Add to crontab for daily backup
crontab -e
# Add: 0 2 * * * /home/ubuntu/backup.sh
```

### 3. Log Rotation

```bash
# Create logrotate config
sudo nano /etc/logrotate.d/mini-app-bot
```

```
/var/log/mini-app-bot/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 644 ubuntu ubuntu
    postrotate
        systemctl reload mini-app-bot
    endscript
}
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Database Connection Error

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check connection
psql -h localhost -U app_user -d mini_app_bot_telegram

# Reset password if needed
sudo -u postgres psql
ALTER USER app_user PASSWORD 'new_password';
```

#### 2. Port Already in Use

```bash
# Check what's using port 3000
sudo netstat -tulpn | grep :3000
sudo lsof -i :3000

# Kill process if needed
sudo kill -9 <PID>
```

#### 3. Permission Issues

```bash
# Fix file permissions
sudo chown -R ubuntu:ubuntu /home/ubuntu/mini-app-bot-telegram
chmod +x /home/ubuntu/mini-app-bot-telegram/backend/main
```

#### 4. Environment Variables Not Loading

```bash
# Check .env file
cat .env

# Ensure proper file location
ls -la /home/ubuntu/mini-app-bot-telegram/backend/.env
```

## 📈 Performance Optimization

### 1. Database Optimization

```sql
-- Add indexes for frequently queried fields
CREATE INDEX idx_overtime_telegram_user_id ON overtimes(telegram_user_id);
CREATE INDEX idx_overtime_date ON overtimes(date);
CREATE INDEX idx_telegram_users_telegram_id ON telegram_users(telegram_id);
CREATE INDEX idx_api_keys_api_key ON api_keys(api_key);
```

### 2. Application Optimization

- Enable GZIP compression in Nginx
- Implement connection pooling for database
- Add Redis for caching (optional)
- Set up CDN for static assets

### 3. Security Hardening

- Use strong JWT secrets
- Implement rate limiting
- Set up firewall rules
- Regular security updates
- Use HTTPS only
- Implement proper CORS settings

---

## 📞 Support

Jika mengalami masalah dalam deployment, periksa:

1. Logs aplikasi: `sudo journalctl -u mini-app-bot -f`
2. Database connectivity
3. Environment variables
4. File permissions
5. Network connectivity

Untuk bantuan lebih lanjut, buat issue di repository atau hubungi tim development.
