# 📚 API Documentation - Mini App Bot Telegram Backend

Dokumentasi lengkap API untuk sistem backend Mini App Bot Telegram.

## 🔧 Base Information

- **Base URL**: `http://localhost:3000` (Development)
- **API Version**: v1
- **Content-Type**: `application/json`
- **Authentication**: API Key (`X-API-Key` header) atau JWT Token (`Authorization: Bearer <token>`)

## 🔐 Authentication

### 1. API Key Authentication
```bash
# Gunakan header X-API-Key
curl -H "X-API-Key: your-api-key-here" http://localhost:3000/v1/endpoint
```

### 2. JWT Authentication
```bash
# Gunakan header Authorization dengan Bearer token
curl -H "Authorization: Bearer your-jwt-token-here" http://localhost:3000/v1/endpoint
```

## 🛠️ API Endpoints

### Health Check

#### `GET /health`
**Deskripsi**: Check if the API is running  
**Authentication**: None  

**Response**:
```json
{
  "status": "ok"
}
```

**cURL Example**:
```bash
curl -X GET http://localhost:3000/health
```

---

## 🔐 Authentication Endpoints

### Login

#### `POST /v1/auth/login`
**Deskripsi**: Login dengan username dan password untuk mendapatkan JWT token  
**Authentication**: None  

**Request Body**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expire_at": "2025-09-03 22:31:05"
  },
  "message": "Login successful"
}
```

**Response Error (400)**:
```json
{
  "code": 400,
  "data": [
    {"username": "Username is required"}
  ],
  "message": "Invalid payload"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:3000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Register

#### `POST /v1/auth/register`
**Deskripsi**: Registrasi user baru  
**Authentication**: None  

**Request Body**:
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Validation Rules**:
- `username`: Required, min 3 characters, max 20 characters
- `email`: Required, valid email format
- `password`: Required, min 8 characters

**Response Success (201)**:
```json
{
  "code": 201,
  "data": {
    "id": 3,
    "username": "testuser",
    "email": "testuser@example.com",
    "created_at": "2025-09-02T22:30:51.800578+07:00",
    "updated_at": "2025-09-02T22:30:51.800578+07:00"
  },
  "message": "User created successfully"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:3000/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "testuser@example.com",
    "password": "password123"
  }'
```

---

## 👤 User Management Endpoints

### Get Current User Details

#### `GET /v1/user/detail-me`
**Deskripsi**: Mendapatkan detail user yang sedang login  
**Authentication**: Required (API Key atau JWT)  

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "user": {
      "id": 1,
      "username": "nyuuk",
      "email": "adnan@nyuuk.my.id",
      "created_at": "2025-08-27T15:16:35.554792+07:00",
      "updated_at": "2025-08-27T15:16:35.554792+07:00"
    },
    "auth_info": {
      "user_id": 1,
      "auth_type": "api_key",
      "expire_at": "0001-01-01T00:00:00Z"
    }
  },
  "message": "User retrieved successfully"
}
```

**cURL Example**:
```bash
curl -X GET http://localhost:3000/v1/user/detail-me \
  -H "X-API-Key: your-api-key-here"
```

### Get All Users

#### `GET /v1/user/`
**Deskripsi**: Mendapatkan semua user (admin only)  
**Authentication**: Required (API Key atau JWT)  

**Response Success (200)**:
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "username": "nyuuk",
      "email": "adnan@nyuuk.my.id",
      "created_at": "2025-08-27T15:16:35.554792+07:00",
      "updated_at": "2025-08-27T15:16:35.554792+07:00"
    }
  ],
  "message": "Users retrieved successfully"
}
```

### Get User by ID

#### `GET /v1/user/{id}`
**Deskripsi**: Mendapatkan detail user berdasarkan ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `id` (path): User ID (integer)

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "id": 3,
    "username": "testuser",
    "email": "testuser@example.com",
    "created_at": "2025-09-02T22:30:51.800578+07:00",
    "updated_at": "2025-09-02T22:30:51.800578+07:00"
  },
  "message": "User retrieved successfully"
}
```

**cURL Example**:
```bash
curl -X GET http://localhost:3000/v1/user/3 \
  -H "X-API-Key: your-api-key-here"
```

### Delete User

#### `DELETE /v1/user/{id}`
**Deskripsi**: Menghapus user berdasarkan ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `id` (path): User ID (integer)

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "id": 4,
    "username": "testuser2",
    "email": "testuser2@example.com",
    "created_at": "2025-09-02T22:31:21.070289+07:00",
    "updated_at": "2025-09-02T22:31:21.070289+07:00"
  },
  "message": "User deleted successfully"
}
```

### Get API Keys

#### `GET /v1/user/api-key`
**Deskripsi**: Mendapatkan semua API key dari user yang sedang aktif  
**Authentication**: Required (API Key atau JWT)  

**Response Success (200)**:
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "api_key": "933a7c601ddd95a888f1cfe802e66561736eda27ec0a9da94cffb591ca345cb7",
      "description": "My API Key for testing",
      "is_active": true,
      "created_at": "2025-08-27T15:58:39.086484+07:00",
      "expired_at": "2025-08-27T22:04:05+07:00"
    }
  ],
  "message": "API key retrieved successfully"
}
```

### Create API Key

#### `POST /v1/user/api-key`
**Deskripsi**: Membuat API key baru  
**Authentication**: Required (API Key atau JWT)  

**Request Body**:
```json
{
  "description": "string",
  "is_active": true,
  "expired_at": "2025-12-31T23:59:59Z"
}
```

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "description": "Test API Key",
    "is_active": true,
    "expired_at": "2025-12-31T23:59:59Z",
    "api_key": "0c407e385b046e67f60c679793718dac074902aa6a768e1aab8c0b55269c2823"
  },
  "message": "API key created successfully"
}
```

---

## 📱 Telegram User Management

### Create Telegram User

#### `POST /v1/telegram/`
**Deskripsi**: Membuat telegram user baru untuk user yang sedang aktif  
**Authentication**: Required (API Key atau JWT)  

**Request Body**:
```json
{
  "telegram_id": 999888777,
  "first_name": "Test",
  "last_name": "User",
  "username": "testuser_tg"
}
```

**Response Success (201)**:
```json
{
  "code": 201,
  "data": {
    "id": 7,
    "user_id": 1,
    "telegram_id": 999888777,
    "username": "testuser_tg",
    "first_name": "Test",
    "last_name": "User",
    "created_at": "2025-09-02T22:31:51.069348+07:00",
    "updated_at": "2025-09-02T22:31:51.069348+07:00"
  },
  "message": "Telegram user created successfully"
}
```

### Get All Telegram Users

#### `GET /v1/telegram/`
**Deskripsi**: Mendapatkan semua telegram user untuk user yang sedang aktif  
**Authentication**: Required (API Key atau JWT)  

**Response Success (200)**:
```json
{
  "code": 200,
  "data": [
    {
      "id": 6,
      "user_id": 1,
      "telegram_id": 1234567892,
      "username": "nyuuk",
      "first_name": "Nyuuk",
      "last_name": "Nyuuk",
      "created_at": "2025-08-28T10:50:12.750974+07:00",
      "updated_at": "2025-08-28T10:50:12.750974+07:00"
    }
  ],
  "message": "Telegram user found successfully"
}
```

### Get Telegram User by ID

#### `GET /v1/telegram/{telegram_id}`
**Deskripsi**: Mendapatkan telegram user berdasarkan telegram ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `telegram_id` (path): Telegram ID (integer)

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "id": 7,
    "user_id": 1,
    "telegram_id": 999888777,
    "username": "testuser_tg",
    "first_name": "Test",
    "last_name": "User",
    "created_at": "2025-09-02T22:31:51.069348+07:00",
    "updated_at": "2025-09-02T22:31:51.069348+07:00"
  },
  "message": "Telegram user found successfully"
}
```

### Update Telegram User

#### `PUT /v1/telegram/{telegram_id}`
**Deskripsi**: Update telegram user berdasarkan telegram ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `telegram_id` (path): Telegram ID (integer)

**Request Body**:
```json
{
  "telegram_id": 999888777,
  "first_name": "Updated Test",
  "last_name": "Updated User",
  "username": "updated_testuser_tg"
}
```

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "id": 7,
    "user_id": 1,
    "telegram_id": 999888777,
    "username": "updated_testuser_tg",
    "first_name": "Updated Test",
    "last_name": "Updated User",
    "created_at": "2025-09-02T22:31:51.069348+07:00",
    "updated_at": "2025-09-02T22:31:51.069348+07:00"
  },
  "message": "Telegram user updated successfully"
}
```

### Delete Telegram User

#### `DELETE /v1/telegram/{telegram_id}`
**Deskripsi**: Hapus telegram user berdasarkan telegram ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `telegram_id` (path): Telegram ID (integer)

**Response Success (200)**:
```json
{
  "code": 200,
  "data": null,
  "message": "Telegram user deleted successfully"
}
```

---

## ⏰ Overtime Management

### Create Overtime Record

#### `POST /v1/overtime/`
**Deskripsi**: Membuat record overtime baru  
**Authentication**: Required (API Key atau JWT)  

**Request Body**:
```json
{
  "telegram_id": 1234567892,
  "date": "2025-01-17",
  "time_start": "09:00:00",
  "time_stop": "18:00:00",
  "break_duration": 1.0,
  "duration": 8.0,
  "description": "Testing overtime endpoint",
  "category": "Testing"
}
```

**Field Descriptions**:
- `telegram_id`: ID telegram user
- `date`: Tanggal overtime (YYYY-MM-DD)
- `time_start`: Waktu mulai (HH:MM:SS)
- `time_stop`: Waktu selesai (HH:MM:SS)
- `break_duration`: Durasi istirahat dalam jam (decimal)
- `duration`: Total durasi overtime dalam jam (decimal)
- `description`: Deskripsi pekerjaan
- `category`: Kategori pekerjaan

**Response Success (201)**:
```json
{
  "code": 201,
  "data": {
    "time_start": "09:00:00",
    "time_stop": "18:00:00",
    "id": 9,
    "telegram_user_id": 6,
    "date": "2025-01-17T00:00:00+07:00",
    "break_duration": 1,
    "duration": 8,
    "description": "Testing overtime endpoint",
    "category": "Testing",
    "updated_at": "2025-09-02T22:32:21.966943+07:00",
    "created_at": "2025-09-02T22:32:21.966943+07:00"
  },
  "message": "Overtime record created successfully"
}
```

### Get All Overtime Records by Telegram ID

#### `GET /v1/overtime/telegram/{telegram_id}`
**Deskripsi**: Mendapatkan semua record overtime berdasarkan telegram ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `telegram_id` (path): Telegram ID (integer)

**Response Success (200)**:
```json
{
  "code": 200,
  "data": [
    {
      "time_start": "09:00:00",
      "time_stop": "18:00:00",
      "id": 9,
      "telegram_user_id": 6,
      "date": "2025-01-17T00:00:00Z",
      "break_duration": 1,
      "duration": 8,
      "description": "Testing overtime endpoint",
      "category": "Testing",
      "updated_at": "2025-09-02T22:32:21.966943+07:00",
      "created_at": "2025-09-02T22:32:21.966943+07:00"
    }
  ],
  "message": "Overtime records retrieved successfully"
}
```

### Get Overtime Record by Date

#### `POST /v1/overtime/by-date`
**Deskripsi**: Mendapatkan record overtime berdasarkan tanggal spesifik  
**Authentication**: Required (API Key atau JWT)  

**Request Body**:
```json
{
  "date": "2025-01-17",
  "telegram_id": 1234567892
}
```

**Response Success (200)**:
```json
{
  "code": 200,
  "data": [
    {
      "time_start": "09:00:00",
      "time_stop": "18:00:00",
      "id": 9,
      "telegram_user_id": 6,
      "date": "2025-01-17T00:00:00Z",
      "break_duration": 1,
      "duration": 8,
      "description": "Testing overtime endpoint",
      "category": "Testing",
      "updated_at": "2025-09-02T22:32:21.966943+07:00",
      "created_at": "2025-09-02T22:32:21.966943+07:00"
    }
  ],
  "message": "Overtime record retrieved successfully"
}
```

**Response Not Found (404)**:
```json
{
  "code": 404,
  "data": null,
  "message": "Overtime record not found for this date"
}
```

### Get Overtime Records Between Dates

#### `POST /v1/overtime/between-dates`
**Deskripsi**: Mendapatkan record overtime dalam rentang tanggal  
**Authentication**: Required (API Key atau JWT)  

**Request Body**:
```json
{
  "telegram_id": 1234567892,
  "start_date": "2025-01-15",
  "end_date": "2025-01-17"
}
```

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "period": {
      "end_date": "2025-01-17T00:00:00+07:00",
      "start_date": "2025-01-15T00:00:00+07:00"
    },
    "records": [
      {
        "time_start": "09:00:00",
        "time_stop": "18:00:00",
        "id": 9,
        "telegram_user_id": 6,
        "date": "2025-01-17T00:00:00Z",
        "break_duration": 1,
        "duration": 8,
        "description": "Testing overtime endpoint",
        "category": "Testing",
        "updated_at": "2025-09-02T22:32:21.966943+07:00",
        "created_at": "2025-09-02T22:32:21.966943+07:00"
      }
    ],
    "records_count": 3,
    "total_duration": 24
  },
  "message": "Overtime records retrieved successfully"
}
```

### Get Overtime Record by ID

#### `GET /v1/overtime/{id}`
**Deskripsi**: Mendapatkan record overtime berdasarkan ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `id` (path): Overtime record ID (integer)

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "time_start": "09:00:00",
    "time_stop": "18:00:00",
    "id": 9,
    "telegram_user_id": 6,
    "date": "2025-01-17T00:00:00Z",
    "break_duration": 1,
    "duration": 8,
    "description": "Testing overtime endpoint",
    "category": "Testing",
    "updated_at": "2025-09-02T22:32:21.966943+07:00",
    "created_at": "2025-09-02T22:32:21.966943+07:00"
  },
  "message": "Overtime record retrieved successfully"
}
```

### Update Overtime Record

#### `PUT /v1/overtime/{id}`
**Deskripsi**: Update record overtime berdasarkan ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `id` (path): Overtime record ID (integer)

**Request Body**:
```json
{
  "telegram_id": 1234567892,
  "date": "2025-01-17",
  "time_start": "08:00:00",
  "time_stop": "19:00:00",
  "break_duration": 1.5,
  "duration": 9.5,
  "description": "Updated testing overtime endpoint",
  "category": "Updated Testing"
}
```

**Response Success (200)**:
```json
{
  "code": 200,
  "data": {
    "time_start": "08:00:00",
    "time_stop": "19:00:00",
    "id": 0,
    "telegram_user_id": 6,
    "date": "2025-01-17T00:00:00+07:00",
    "break_duration": 1.5,
    "duration": 9.5,
    "description": "Updated testing overtime endpoint",
    "category": "Updated Testing",
    "updated_at": "0001-01-01T00:00:00Z",
    "created_at": "0001-01-01T00:00:00Z"
  },
  "message": "Overtime record updated successfully"
}
```

### Delete Overtime Record

#### `DELETE /v1/overtime/{id}`
**Deskripsi**: Hapus record overtime berdasarkan ID  
**Authentication**: Required (API Key atau JWT)  

**Parameters**:
- `id` (path): Overtime record ID (integer)

**Response Success (200)**:
```json
{
  "code": 200,
  "data": null,
  "message": "Overtime record deleted successfully"
}
```

---

## 📊 Response Codes

| Code | Description |
|------|-------------|
| 200  | OK - Request berhasil |
| 201  | Created - Resource berhasil dibuat |
| 400  | Bad Request - Request tidak valid |
| 401  | Unauthorized - Authentication diperlukan |
| 403  | Forbidden - Akses ditolak |
| 404  | Not Found - Resource tidak ditemukan |
| 500  | Internal Server Error - Error server |

## 🔍 Error Response Format

Semua error response menggunakan format standar:

```json
{
  "code": 400,
  "data": [
    {"field": "error message"}
  ],
  "message": "Error description"
}
```

atau

```json
{
  "code": 500,
  "data": null,
  "message": "Internal server error"
}
```

## 🚀 Quick Start

### 1. Get API Key
```bash
# Register user
curl -X POST http://localhost:3000/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# Login to get JWT
curl -X POST http://localhost:3000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# Create API Key (using JWT from login)
curl -X POST http://localhost:3000/v1/user/api-key \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "description": "My API Key",
    "is_active": true,
    "expired_at": "2025-12-31T23:59:59Z"
  }'
```

### 2. Create Telegram User
```bash
curl -X POST http://localhost:3000/v1/telegram/ \
  -H "Content-Type: application/json" \
  -H "X-API-Key: YOUR_API_KEY" \
  -d '{
    "telegram_id": 123456789,
    "first_name": "John",
    "last_name": "Doe",
    "username": "johndoe"
  }'
```

### 3. Create Overtime Record
```bash
curl -X POST http://localhost:3000/v1/overtime/ \
  -H "Content-Type: application/json" \
  -H "X-API-Key: YOUR_API_KEY" \
  -d '{
    "telegram_id": 123456789,
    "date": "2025-01-17",
    "time_start": "09:00:00",
    "time_stop": "18:00:00",
    "break_duration": 1.0,
    "duration": 8.0,
    "description": "Working on project development",
    "category": "Development"
  }'
```

## 🛠️ Tools & Testing

### Postman Collection
Import file collection yang tersedia di `docs/rest-api/` untuk testing dengan Postman.

### cURL Scripts
Semua contoh cURL di atas dapat disalin dan dijalankan langsung di terminal.

### HTTP Files
Gunakan file `.http` di direktori `docs/rest-api/` untuk testing dengan VS Code REST Client extension.

## 📞 Support

Untuk bantuan atau pertanyaan terkait API:
1. Periksa dokumentasi ini
2. Cek endpoint `/health` untuk memastikan server berjalan
3. Periksa logs aplikasi untuk debugging
4. Hubungi tim development untuk bantuan lebih lanjut

---

**Last Updated**: September 2, 2025  
**API Version**: 1.0
