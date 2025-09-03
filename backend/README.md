# 🚀 Mini App Bot Telegram Backend

Backend API system untuk Mini App Bot Telegram yang dibangun dengan Go Fiber dan PostgreSQL.

## 📚 Dokumentasi

### 📖 API Documentation
Kami menyediakan beberapa cara untuk mengakses dokumentasi API:

1. **Swagger UI** (Interactive): http://localhost:3000/swagger/index.html
2. **API Documentation (Manual)**: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
3. **HTTP Files**: Lihat folder `docs/rest-api/` untuk testing dengan REST Client

### 🚀 Deployment Guide
Panduan lengkap deployment tersedia di: [DEPLOYMENT.md](./DEPLOYMENT.md)

## 🔧 Quick Start

### Prerequisites
- Go 1.19+
- PostgreSQL 12+
- Git

### 1. Clone & Setup
```bash
git clone <repository-url>
cd mini-app-bot-telegram/backend

# Install dependencies
go mod download

# Setup environment
cp .env.example .env
# Edit .env sesuai konfigurasi Anda
```

### 2. Database Setup
```sql
-- Buat database PostgreSQL
CREATE DATABASE mini_app_bot_telegram;
```

### 3. Run Application
```bash
# Development mode
go run main.go

# Build dan run
go build -o main .
./main
```

Server akan berjalan di `http://localhost:3000`

## 📊 API Features

### ✅ Fitur Utama
- **Authentication**: API Key & JWT Token
- **User Management**: CRUD operations untuk users
- **Telegram Integration**: Manajemen telegram users
- **Overtime Management**: Sistem pencatatan overtime
- **Logging**: Comprehensive logging system
- **Swagger Documentation**: Interactive API documentation

### 🔐 Authentication Methods
1. **API Key**: Header `X-API-Key: your-api-key`
2. **JWT Token**: Header `Authorization: Bearer your-jwt-token`

## 🛠️ Development Tools

### Swagger Documentation
```bash
# Generate swagger docs (jika ada perubahan anotasi)
~/go/bin/swag init -g main.go --output docs

# Akses Swagger UI
open http://localhost:3000/swagger/index.html
```

### Testing Endpoints
```bash
# Health check
curl http://localhost:3000/health

# Login example
curl -X POST http://localhost:3000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Protected endpoint example
curl -X GET http://localhost:3000/v1/user/detail-me \
  -H "X-API-Key: your-api-key-here"
```

## 📁 Project Structure

```
backend/
├── app/
│   ├── controllers/     # HTTP handlers
│   ├── entities/        # Database models
│   ├── middlewares/     # HTTP middlewares
│   ├── payloads/        # Request/Response structures
│   ├── pkg/
│   │   ├── database/    # Database connection & migration
│   │   └── helpers/     # Utility functions
│   ├── repositories/    # Database access layer
│   └── services/        # Business logic layer
├── docs/
│   ├── rest-api/        # HTTP test files
│   ├── swagger.json     # Generated swagger spec
│   ├── swagger.yaml     # Generated swagger spec
│   └── docs.go          # Generated swagger docs
├── main.go              # Application entry point
├── go.mod               # Go dependencies
├── go.sum               # Go dependency checksums
├── .env                 # Environment variables
├── API_DOCUMENTATION.md # Manual API documentation
├── DEPLOYMENT.md        # Deployment guide
└── README.md            # This file
```

## 🔄 Available Endpoints

### Health & Documentation
- `GET /health` - Health check
- `GET /swagger/*` - Swagger UI

### Authentication (Public)
- `POST /v1/auth/login` - User login
- `POST /v1/auth/register` - User registration

### User Management (Protected)
- `GET /v1/user/detail-me` - Get current user details
- `GET /v1/user/` - Get all users
- `POST /v1/user/` - Create user
- `GET /v1/user/{id}` - Get user by ID
- `DELETE /v1/user/{id}` - Delete user
- `GET /v1/user/api-key` - Get user's API keys
- `POST /v1/user/api-key` - Create new API key

### Telegram Management (Protected)
- `POST /v1/telegram/` - Create telegram user
- `GET /v1/telegram/` - Get all telegram users
- `GET /v1/telegram/{telegram_id}` - Get telegram user by ID
- `PUT /v1/telegram/{telegram_id}` - Update telegram user
- `DELETE /v1/telegram/{telegram_id}` - Delete telegram user

### Overtime Management (Protected)
- `POST /v1/overtime/` - Create overtime record
- `GET /v1/overtime/telegram/{telegram_id}` - Get all overtime records
- `POST /v1/overtime/by-date` - Get overtime by specific date
- `POST /v1/overtime/between-dates` - Get overtime between dates
- `GET /v1/overtime/{id}` - Get overtime by ID
- `PUT /v1/overtime/{id}` - Update overtime record
- `DELETE /v1/overtime/{id}` - Delete overtime record

## 🔍 Testing

### API Testing Tools
1. **Swagger UI**: Interactive testing - http://localhost:3000/swagger/
2. **REST Client**: VS Code extension dengan file `.http`
3. **Postman**: Import collection dari `docs/rest-api/`
4. **cURL**: Command line testing

### Example API Key
Untuk testing, gunakan API key yang valid atau buat baru melalui endpoint `/v1/user/api-key`.

## 🔧 Environment Variables

Buat file `.env` berdasarkan `.env.example`:

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

### Environment Variables Description

- `DATABASE_URL`: PostgreSQL connection string
- `ENV`: Application environment (development, staging, production)
- `LOG_LEVEL`: Logging level (trace, debug, info, warn, error, fatal, panic)
- `JWT_SECRET_KEY`: Secret key untuk JWT token generation
- `JWT_EXPIRATION`: JWT token expiration dalam jam

## 🐳 Docker Support

### Multi-Architecture Support
Project ini mendukung build untuk different architectures (ARM64 & AMD64):

```bash
# Build image
docker build -t mini-app-bot-telegram-backend .

# Run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f backend
```

## 📊 Performance & Monitoring

### Logging
- **Request Logging**: Semua HTTP requests dicatat
- **Database Logging**: Query performance monitoring
- **Security Logging**: Authentication attempts
- **Performance Logging**: Response time tracking

### Health Check
```bash
curl http://localhost:3000/health
```

## 🚨 Troubleshooting

### Common Issues

1. **Database Connection Error**
   ```bash
   # Check PostgreSQL status
   sudo systemctl status postgresql
   
   # Test connection
   psql -h localhost -U your_user -d mini_app_bot_telegram
   ```

2. **Port Already in Use**
   ```bash
   # Check what's using port 3000
   lsof -i :3000
   
   # Kill process
   kill -9 <PID>
   ```

3. **Swagger Not Loading**
   ```bash
   # Regenerate swagger docs
   ~/go/bin/swag init -g main.go --output docs
   
   # Restart application
   go run main.go
   ```

4. **JWT Token Expired**
   ```bash
   # Login again to get new token
   curl -X POST http://localhost:3000/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"your_username","password":"your_password"}'
   ```

## 📈 Development Workflow

### 1. Add New Endpoint
1. Create handler in `controllers/`
2. Add Swagger annotations
3. Update routes in `main.go`
4. Regenerate docs: `~/go/bin/swag init -g main.go --output docs`
5. Test endpoint

### 2. Database Changes
1. Update entities in `entities/`
2. Update repositories in `repositories/`
3. Update services in `services/`
4. Run migration (if needed)

### 3. Testing
1. Use Swagger UI for interactive testing
2. Create `.http` files for automated testing
3. Test with different authentication methods

## 🔗 Useful Links

- **Swagger UI**: http://localhost:3000/swagger/
- **Health Check**: http://localhost:3000/health
- **API Documentation**: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
- **Deployment Guide**: [DEPLOYMENT.md](./DEPLOYMENT.md)

## 📞 Support

Untuk bantuan atau pertanyaan:
1. Periksa [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
2. Test dengan Swagger UI
3. Periksa logs aplikasi
4. Hubungi tim development

---

**Tech Stack**: Go Fiber, PostgreSQL, GORM, JWT, Swagger  
**Version**: 1.0  
**Last Updated**: September 2, 2025
