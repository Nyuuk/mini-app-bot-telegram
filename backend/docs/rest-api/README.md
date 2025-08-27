# REST API Testing dengan HTTP Files

## Overview
Dokumentasi ini menjelaskan cara menggunakan variable dan environment di file .http untuk testing REST API.

## Cara Menggunakan Variable di File .http

### 1. **Variable Lokal (dalam file yang sama)**

```http
### Define variables di awal file
@baseUrl = http://localhost:3000
@apiVersion = v1
@jwtToken = your_jwt_token_here
@apiKey = your_api_key_here

### Gunakan variable dengan {{variableName}}
GET {{baseUrl}}/{{apiVersion}}/user/detail-me
Authorization: Bearer {{jwtToken}}
```

### 2. **Environment Variables (http-client.env.json)**

Buat file `http-client.env.json` di folder yang sama:

```json
{
  "development": {
    "baseUrl": "http://localhost:3000",
    "apiVersion": "v1",
    "jwtToken": "your_jwt_token_here",
    "apiKey": "your_api_key_here",
    "telegramToken": "your_telegram_bot_token",
    "chatId": "123456789",
    "userId": "987654321"
  },
  "production": {
    "baseUrl": "https://your-production-domain.com",
    "apiVersion": "v1",
    "jwtToken": "your_production_jwt_token",
    "apiKey": "your_production_api_key",
    "telegramToken": "your_production_telegram_bot_token",
    "chatId": "123456789",
    "userId": "987654321"
  }
}
```

Kemudian gunakan di file .http:

```http
### Gunakan environment variables
GET {{baseUrl}}/{{apiVersion}}/user/detail-me
Authorization: Bearer {{jwtToken}}
```

### 3. **Import Variable dari File Lain**

**File: telegram.http**
```http
### Variables untuk Telegram Bot API
@baseUrl = http://localhost:3000
@apiVersion = v1
@telegramToken = your_telegram_bot_token_here
@chatId = 123456789
@userId = 987654321
@jwtToken = your_jwt_token_here
@apiKey = your_api_key_here

### Telegram Bot Endpoints
POST {{baseUrl}}/{{apiVersion}}/telegram/send-message
Content-Type: application/json
Authorization: Bearer {{jwtToken}}

{
    "chat_id": {{chatId}},
    "text": "Hello from REST API!"
}
```

**File: auth.http**
```http
### Import variables dari telegram.http
# @import "telegram.http"

### Auth Endpoints menggunakan variable dari telegram.http
POST {{baseUrl}}/{{apiVersion}}/auth/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123"
}
```

### 4. **Dynamic Variables (Response Variables)**

```http
### Login dan simpan token
# @name login
POST {{baseUrl}}/{{apiVersion}}/auth/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123"
}

### Gunakan token dari response login
GET {{baseUrl}}/{{apiVersion}}/user/detail-me
Authorization: Bearer {{login.response.body.data.token}}
```

### 5. **Conditional Variables**

```http
### Variables dengan kondisi
@baseUrl = {{$processEnv NODE_ENV === 'production' ? 'https://api.production.com' : 'http://localhost:3000'}}
@apiVersion = v1
@jwtToken = {{$processEnv JWT_TOKEN}}

### Gunakan variable
GET {{baseUrl}}/{{apiVersion}}/user/detail-me
Authorization: Bearer {{jwtToken}}
```

## Struktur File

```
backend/docs/rest-api/
├── http-client.env.json          # Environment variables
├── telegram.http                 # Telegram Bot API endpoints
├── auth.http                     # Authentication endpoints
├── user.http                     # User management endpoints
└── README.md                     # Dokumentasi ini
```

## Best Practices

### 1. **Organisasi Variable**
- Gunakan prefix yang jelas: `@baseUrl`, `@apiVersion`, `@jwtToken`
- Kelompokkan variable berdasarkan fungsi
- Gunakan environment untuk environment-specific values

### 2. **Security**
- Jangan commit token atau API key ke repository
- Gunakan environment variables untuk sensitive data
- Gunakan `.env` file untuk local development

### 3. **Maintainability**
- Gunakan satu file environment untuk semua variable
- Dokumentasikan setiap variable
- Gunakan naming convention yang konsisten

### 4. **Testing**
- Buat test cases untuk berbagai scenario
- Test error cases dan edge cases
- Gunakan different environments untuk testing

## Contoh Penggunaan Lengkap

### Environment File (http-client.env.json)
```json
{
  "development": {
    "baseUrl": "http://localhost:3000",
    "apiVersion": "v1",
    "jwtToken": "dev_jwt_token",
    "apiKey": "dev_api_key",
    "telegramToken": "dev_telegram_token",
    "chatId": "123456789",
    "userId": "987654321"
  },
  "staging": {
    "baseUrl": "https://staging-api.example.com",
    "apiVersion": "v1",
    "jwtToken": "staging_jwt_token",
    "apiKey": "staging_api_key",
    "telegramToken": "staging_telegram_token",
    "chatId": "123456789",
    "userId": "987654321"
  },
  "production": {
    "baseUrl": "https://api.example.com",
    "apiVersion": "v1",
    "jwtToken": "prod_jwt_token",
    "apiKey": "prod_api_key",
    "telegramToken": "prod_telegram_token",
    "chatId": "123456789",
    "userId": "987654321"
  }
}
```

### Main Test File (api-tests.http)
```http
### Import environment
# @import "http-client.env.json"

### Test Authentication
# @name login
POST {{baseUrl}}/{{apiVersion}}/auth/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123"
}

### Test User Profile dengan JWT
GET {{baseUrl}}/{{apiVersion}}/user/detail-me
Authorization: Bearer {{login.response.body.data.token}}

### Test User Profile dengan API Key
GET {{baseUrl}}/{{apiVersion}}/user/detail-me
X-API-Key: {{apiKey}}

### Test Telegram Bot
POST {{baseUrl}}/{{apiVersion}}/telegram/send-message
Content-Type: application/json
Authorization: Bearer {{login.response.body.data.token}}

{
    "chat_id": {{chatId}},
    "text": "Test message from API"
}
```

## Tips dan Trik

1. **Gunakan VS Code REST Client extension** untuk testing yang lebih mudah
2. **Gunakan environment switching** untuk test di berbagai environment
3. **Gunakan response variables** untuk chaining requests
4. **Gunakan comments** untuk dokumentasi
5. **Gunakan conditional variables** untuk dynamic values

## Troubleshooting

### Variable Not Found
- Pastikan variable didefinisikan sebelum digunakan
- Pastikan nama variable benar (case sensitive)
- Pastikan environment file ada dan benar

### Import Not Working
- Pastikan path import benar
- Pastikan file yang di-import ada
- Pastikan syntax import benar

### Environment Not Loading
- Pastikan nama environment benar
- Pastikan file `http-client.env.json` ada
- Pastikan format JSON valid
