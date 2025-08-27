# Logging System Documentation

## Overview
Sistem logging menggunakan zerolog dengan fitur-fitur yang komprehensif untuk tracking, debugging, dan monitoring aplikasi.

## Konfigurasi

### Environment Variables
```env
# Environment
ENV=development

# Log Level (trace, debug, info, warn, error, fatal, panic)
LOG_LEVEL=info
```

### Log Levels
- `trace`: Logging paling detail, untuk debugging mendalam
- `debug`: Informasi debugging
- `info`: Informasi umum aplikasi
- `warn`: Peringatan yang tidak fatal
- `error`: Error yang terjadi
- `fatal`: Error fatal yang menyebabkan aplikasi berhenti
- `panic`: Panic yang terjadi

## Fitur Logging

### 1. HTTP Request Logging
Otomatis mencatat semua HTTP request dengan informasi:
- Method, Path, IP Address
- User Agent, Duration
- Status Code, User ID

### 2. Database Operation Logging
Mencatat operasi database:
- Operation type (create, read, update, delete)
- Table name
- Duration
- Rows affected
- Error (jika ada)

### 3. Authentication Logging
Mencatat event autentikasi:
- Login attempts
- Login success/failure
- User ID tracking

### 4. Business Logic Logging
Mencatat event bisnis:
- User creation
- API key generation
- Business operations

### 5. Security Logging
Mencatat event keamanan:
- Failed authentication
- Missing user agent
- Suspicious activities

### 6. Performance Logging
Mencatat metrik performa:
- Request duration > 1 second
- Database query performance
- Resource usage

## Cara Penggunaan

### 1. Initialize Logger
```go
func main() {
    helpers.InitLogger()
    // ... rest of your code
}
```

### 2. Apply Middleware
```go
app := fiber.New()

// Apply logging middlewares
app.Use(middlewares.LoggingMiddleware())
app.Use(middlewares.DatabaseLoggingMiddleware())
app.Use(middlewares.PerformanceLoggingMiddleware())
app.Use(middlewares.SecurityLoggingMiddleware())
```

### 3. Log Events di Controller
```go
// Log authentication
helpers.LogAuth("login_attempt", "anonymous", false, map[string]interface{}{
    "ip_address": c.IP(),
    "user_agent": c.Get("User-Agent"),
    "username":   payload.Username,
})

// Log business events
helpers.LogBusiness("user_created", userID, map[string]interface{}{
    "email":    payload.Email,
    "username": payload.Username,
})

// Log errors
helpers.LogError(err, "database_connection", userID, map[string]interface{}{
    "operation": "create_user",
    "table":     "users",
})
```

### 4. Log Database Operations di Repository
```go
func (r *UserRepository) CreateUser(user *entities.User, tx *gorm.DB, c *fiber.Ctx) error {
    start := time.Now()
    
    if err := tx.WithContext(c.Context()).Create(&user).Error; err != nil {
        helpers.LogDatabase("create", user.TableName(), time.Since(start), 0, err)
        return err
    }
    
    helpers.LogDatabase("create", user.TableName(), time.Since(start), 1, nil)
    return nil
}
```

### 5. Direct Logger Usage (untuk kasus khusus)
```go
// Untuk logging yang tidak masuk kategori di atas
helpers.Logger.Info().
    Str("custom_field", "value").
    Int("count", 42).
    Msg("Custom log message")

// Untuk error logging
helpers.Logger.Error().
    Err(err).
    Str("context", "custom_error").
    Msg("Custom error message")
```

## Output Format

### Development Mode
```
2024-01-15T10:30:45Z INF app/middlewares/logging_middleware.go:25 > HTTP Request duration=1860.027375 method=GET path=/v1/user remote_addr=127.0.0.1 status_code=200 type=http_request user_agent=vscode-restclient user_id=anonymous
2024-01-15T10:30:45Z INF app/repositories/user_repository.go:15 > Database Operation operation=create table=users duration=25ms rows_affected=1
2024-01-15T10:30:45Z DBG app/controllers/auth_controller.go:42 > Authentication Success event=login_success user_id=123 success=true
```

### Production Mode (JSON)
```json
{
  "level": "info",
  "time": "2024-01-15T10:30:45Z",
  "caller": "app/middlewares/logging_middleware.go:25",
  "type": "http_request",
  "method": "POST",
  "path": "/v1/auth/login",
  "remote_addr": "127.0.0.1",
  "duration": 150000000,
  "status_code": 200,
  "user_id": "anonymous",
  "message": "HTTP Request"
}
```

## Caller Information

Sistem logging memiliki dua jenis caller information:

### 1. Direct Logger Usage
```go
helpers.Logger.Info().Msg("Direct log")
```
- **Caller**: Menampilkan file dan line number dari code yang memanggil `Logger.Info()` langsung
- **Contoh**: `app/controllers/auth_controller.go:42`

### 2. Custom Logging Functions
```go
helpers.LogAuth("login_attempt", userID, false, details)
```
- **Caller**: Menampilkan file dan line number dari code yang memanggil `LogAuth()` function
- **Contoh**: `app/controllers/auth_controller.go:42`

### Perbedaan Caller Information:
- **Direct Logger**: Caller langsung dari code yang memanggil
- **Custom Functions**: Caller dari code yang memanggil function custom (bukan dari dalam function logger)

## Best Practices

1. **Gunakan Log Level yang Tepat**
   - `info` untuk operasi normal
   - `warn` untuk kondisi yang perlu diperhatikan
   - `error` untuk error yang perlu ditangani
   - `debug` untuk informasi debugging

2. **Include Context yang Relevan**
   - User ID untuk tracking
   - IP Address untuk security
   - Duration untuk performance
   - Error details untuk debugging

3. **Structured Logging**
   - Gunakan map[string]interface{} untuk additional fields
   - Konsisten dalam naming fields
   - Jangan log sensitive data (password, tokens)

4. **Performance Considerations**
   - Log level yang tepat di production
   - Hindari logging berlebihan
   - Gunakan async logging jika diperlukan

5. **Security**
   - Sistem otomatis filter field sensitif
   - Jangan log API keys, passwords, tokens
   - Log hash atau ID untuk tracking

6. **Pemilihan Method Logging**
   - **Custom Functions**: Untuk logging yang terstruktur dan konsisten
   - **Direct Logger**: Untuk logging yang fleksibel dan custom

## Monitoring dan Alerting

### Log Aggregation
- Gunakan tools seperti ELK Stack, Fluentd, atau Logstash
- Centralized logging untuk multiple services
- Log retention policy

### Alerting
- Set up alerts untuk error rate tinggi
- Monitor performance degradation
- Security event notifications

### Metrics
- Request rate dan response time
- Error rate per endpoint
- Database performance metrics
- User activity patterns

## Troubleshooting

### Caller Information Tidak Benar
- **Direct Logger**: Caller akan langsung dari code yang memanggil
- **Custom Functions**: Caller akan dari code yang memanggil function custom

### Format Log Tidak Rapi
- Pastikan `FormatFieldName` menambahkan `=`
- Pastikan `FormatFieldValue` menambahkan spasi

### Sensitive Data Ter-log
- Sistem otomatis filter field dengan nama sensitif
- Tambahkan field name ke `isSensitiveField` jika diperlukan
