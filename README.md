# 🤖 Mini App Bot Telegram

Full-stack aplikasi Mini App Bot Telegram dengan backend Go Fiber dan frontend modern.

## 📁 Project Structure

```
mini-app-bot-telegram/
├── backend/                 # Go Fiber API Backend
│   ├── app/                # Application logic
│   ├── docs/               # API Documentation & Swagger
│   ├── Dockerfile          # Backend container
│   ├── .dockerignore       # Docker ignore rules
│   ├── go.mod              # Go dependencies
│   ├── main.go             # Application entry point
│   └── README.md           # Backend documentation
├── frontend/               # Frontend application (in development)
├── .github/
│   └── workflows/          # GitHub Actions CI/CD
│       ├── ci.yml          # Continuous Integration
│       └── build-and-deploy.yml # Build & Deploy
├── docker-compose.yml      # Development environment
├── docker-compose.prod.yml # Production environment
└── README.md              # This file
```

## 🚀 Quick Start

### Prerequisites

- **Docker & Docker Compose** (Recommended)
- **Go 1.21+** (untuk development)
- **PostgreSQL 15+** (jika tidak menggunakan Docker)

### 1. Clone Repository

```bash
git clone <repository-url>
cd mini-app-bot-telegram
```

### 2. Development dengan Docker (Recommended)

```bash
# Start development environment
docker-compose up --build -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### 3. Manual Development Setup

```bash
# Setup backend
cd backend
cp .env.example .env
# Edit .env dengan konfigurasi yang sesuai

# Install dependencies & run
go mod download
go run main.go
```

## 📊 Services

### Backend API
- **Port**: 3000
- **Health Check**: http://localhost:3000/health
- **Swagger UI**: http://localhost:3000/swagger/
- **Technology**: Go Fiber, PostgreSQL, JWT, GORM

### Database
- **PostgreSQL**: Port 5432
- **Database**: postgres
- **User**: postgres

## 🔧 Environment Configuration

Buat file `.env` di `backend/` folder berdasarkan `.env.example`:

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

## 🐳 Docker Deployment

### Development
```bash
docker-compose up --build -d
```

### Production
```bash
# Set environment variables
export POSTGRES_PASSWORD=secure_password
export JWT_SECRET_KEY=secure_jwt_key
export VERSION=v1.0.0

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

## 🚀 CI/CD dengan GitHub Actions

Project ini menggunakan automated CI/CD pipeline:

### Features
- **Smart Build Detection**: Hanya build komponen yang berubah
- **Multi-platform Support**: linux/amd64 dan linux/arm64
- **Automated Testing**: Unit tests, linting, security scanning
- **Container Registry**: Push ke GitHub Container Registry

### Deployment Process
1. **Create Tag**: `git tag v1.0.0 && git push origin v1.0.0`
2. **Auto Detection**: GitHub Actions detect changes di backend/ atau frontend/
3. **Build & Push**: Build Docker images dengan tag format `v1.0.0-backend`
4. **Deploy**: Deploy ke production environment

### Image Tags
- Backend: `ghcr.io/username/mini-app-bot-telegram:v1.0.0-backend`
- Frontend: `ghcr.io/username/mini-app-bot-telegram:v1.0.0-frontend`

## 📖 Documentation

- **Backend API**: [backend/README.md](./backend/README.md)
- **API Documentation**: [backend/API_DOCUMENTATION.md](./backend/API_DOCUMENTATION.md)
- **Deployment Guide**: [backend/DEPLOYMENT.md](./backend/DEPLOYMENT.md)
- **Swagger UI**: http://localhost:3000/swagger/

## 🔗 API Endpoints

### Health & Documentation
- `GET /health` - Health check
- `GET /swagger/*` - Interactive API documentation

### Authentication (Public)
- `POST /v1/auth/login` - User login
- `POST /v1/auth/register` - User registration

### Protected Endpoints
- **User Management**: `/v1/user/*`
- **Telegram Integration**: `/v1/telegram/*`
- **Overtime Management**: `/v1/overtime/*`

## 🛠️ Development

### Backend Development
```bash
cd backend

# Install dependencies
go mod download

# Run with hot reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Generate swagger docs
swag init -g main.go --output docs
```

### Frontend Development
```bash
cd frontend
# (Frontend masih dalam tahap pengembangan)
```

## 🔍 Testing

### API Testing
1. **Swagger UI**: http://localhost:3000/swagger/
2. **REST Client**: Files di `backend/docs/rest-api/`
3. **cURL**: Command line testing
4. **Postman**: Import dari Swagger spec

### Example Request
```bash
# Health check
curl http://localhost:3000/health

# Login
curl -X POST http://localhost:3000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

## 🔐 Authentication

Aplikasi mendukung dua metode authentication:
1. **API Key**: Header `X-API-Key: your-api-key`
2. **JWT Token**: Header `Authorization: Bearer your-jwt-token`

## 📈 Monitoring

### Health Checks
```bash
# Application health
curl http://localhost:3000/health

# Database connectivity (via application)
curl -H "X-API-Key: your-key" http://localhost:3000/v1/user/detail-me
```

### Logs
```bash
# Development
docker-compose logs -f backend

# Production
docker-compose -f docker-compose.prod.yml logs -f backend
```

## 🚨 Troubleshooting

### Common Issues

1. **Port 3000 already in use**
   ```bash
   lsof -i :3000
   kill -9 <PID>
   ```

2. **Database connection error**
   ```bash
   docker-compose logs postgres
   docker-compose restart postgres
   ```

3. **Docker build fails**
   ```bash
   docker-compose down
   docker system prune -f
   docker-compose up --build
   ```

## 🤝 Contributing

1. Fork repository
2. Create feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push branch: `git push origin feature/new-feature`
5. Submit Pull Request

## 📞 Support

- **Backend Issues**: Check [backend/README.md](./backend/README.md)
- **API Documentation**: http://localhost:3000/swagger/
- **Deployment Issues**: Check [backend/DEPLOYMENT.md](./backend/DEPLOYMENT.md)

---

**Tech Stack**: Go Fiber, PostgreSQL, Docker, GitHub Actions  
**Version**: 1.0  
**Status**: Backend ✅ | Frontend 🚧 (In Development)
