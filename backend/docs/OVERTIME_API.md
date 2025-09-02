# Overtime Management API

## Overview
API untuk mengelola data overtime/lembur karyawan yang terintegrasi dengan Telegram user. Menggunakan pola arsitektur Repository-Service-Controller dengan logging menggunakan MyLogger. API ini menggunakan **Telegram ID** (ID asli dari Telegram) sebagai identifier, bukan ID database internal.

**Timezone**: API menggunakan timezone **Asia/Jakarta** untuk parsing datetime. Client dapat mengirim format datetime tanpa timezone indicator dan akan diparsing sebagai Asia/Jakarta timezone.

## Features
- ✅ Create new overtime record
- ✅ Get all overtime records by telegram user ID
- ✅ Get overtime record by specific date
- ✅ Get overtime records between date range
- ✅ Get overtime record by ID
- ✅ Update overtime record
- ✅ Delete overtime record

## API Endpoints

### Base URL
```
http://localhost:3000/v1/overtime
```

### Authentication
Semua endpoint memerlukan authentication menggunakan JWT token atau API Key.

### 1. Create New Overtime Record
**POST** `/`

**Request Body:**
```json
{
  "telegram_id": 123456789,
  "date": "2024-01-15",
  "time_start": "2024-01-15T09:00:00",
  "time_stop": "2024-01-15T18:00:00",
  "break_duration": 1.0,
  "duration": 8.0,
  "description": "Working on project development",
  "category": "Development"
}
```

**Note**: Semua datetime dikirim dalam format Asia/Jakarta timezone tanpa timezone indicator.

**Response:**
```json
{
  "status": "success",
  "message": "Overtime record created successfully",
  "data": {
    "id": 1,
    "telegram_user_id": 1,
    "date": "2024-01-15T00:00:00Z",
    "time_start": "2024-01-15T09:00:00Z",
    "time_stop": "2024-01-15T18:00:00Z",
    "break_duration": 1.0,
    "duration": 8.0,
    "description": "Working on project development",
    "category": "Development",
    "created_by_user_id": 1,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "TelegramUser": {
      "id": 1,
      "telegram_id": 123456789,
      "username": "john_doe",
      "first_name": "John",
      "last_name": "Doe"
    }
  }
}
```

### 2. Get All Overtime Records by Telegram ID
**GET** `/telegram/{telegram_id}`

**Response:**
```json
{
  "status": "success",
  "message": "Overtime records retrieved successfully",
  "data": [
    {
      "id": 1,
      "telegram_user_id": 1,
      "date": "2024-01-15T00:00:00Z",
      "time_start": "2024-01-15T09:00:00Z",
      "time_stop": "2024-01-15T18:00:00Z",
      "break_duration": 1.0,
      "duration": 8.0,
      "description": "Working on project development",
      "category": "Development",
      "created_by_user_id": 1,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z",
      "User": {...},
      "TelegramUser": {...}
    }
  ]
}
```

### 3. Get Overtime Record by Date
**POST** `/by-date`

**Request Body:**
```json
{
  "telegram_id": 123456789,
  "date": "2024-01-15"
}
```

### 4. Get Overtime Records Between Dates
**POST** `/between-dates`

**Request Body:**
```json
{
  "telegram_id": 123456789,
  "start_date": "2024-01-01",
  "end_date": "2024-01-31"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Overtime records retrieved successfully",
  "data": {
    "records": [...],
    "total_duration": 160.5,
    "records_count": 20,
    "period": {
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-01-31T00:00:00Z"
    }
  }
}
```

### 5. Get Overtime Record by ID
**GET** `/{id}`

### 6. Update Overtime Record
**PUT** `/{id}`

**Request Body:** Same as create request

### 7. Delete Overtime Record
**DELETE** `/{id}`

## Data Validation

### CreateNewRecordOvertime
- `telegram_id`: Required, must be valid telegram ID (int64)
- `date`: Required, format "YYYY-MM-DD" (Asia/Jakarta timezone)
- `time_start`: Required, format "YYYY-MM-DDTHH:MM:SS" (Asia/Jakarta timezone)
- `time_stop`: Required, format "YYYY-MM-DDTHH:MM:SS" (Asia/Jakarta timezone)
- `break_duration`: Optional, must be >= 0
- `duration`: Required, must be > 0
- `description`: Optional, 3-255 characters if provided
- `category`: Optional, 3-255 characters if provided

### Timezone Configuration
- Default timezone: **Asia/Jakarta**
- Environment variable: `TIMEZONE=Asia/Jakarta`
- Supported formats:
  - Date: `"2024-01-15"`
  - DateTime: `"2024-01-15T09:00:00"` or `"2024-01-15 09:00:00"`
  - Time only: `"09:00:00"` or `"09:00"` (disimpan sebagai time only di database)

### Time Storage Behavior
- **Full DateTime**: Disimpan sebagai datetime lengkap
- **Time Only**: Disimpan sebagai time only (tanpa tanggal) di database
- **Database Type**: `time_start` dan `time_stop` menggunakan PostgreSQL `TIME` type

## Error Responses

### 400 Bad Request
```json
{
  "status": "error",
  "message": "Invalid payload",
  "data": [
    {
      "telegram_id": "Telegram ID is required"
    }
  ]
}
```

### 404 Not Found
```json
{
  "status": "error",
  "message": "Overtime record not found",
  "data": null
}
```

### 409 Conflict
```json
{
  "status": "error",
  "message": "Overtime record already exists for this date",
  "data": null
}
```

### 500 Internal Server Error
```json
{
  "status": "error",
  "message": "Internal server error",
  "data": null
}
```

## Database Schema

### Table: `overtimes`
| Field | Type | Description |
|-------|------|-------------|
| `id` | `uint` | Primary key |
| `telegram_user_id` | `uint` | Foreign key to telegram_users |
| `date` | `date` | Overtime date |
| `time_start` | `time` | Start time |
| `time_stop` | `time` | Stop time |
| `break_duration` | `decimal(4,2)` | Break duration in hours |
| `duration` | `decimal(4,2)` | Total overtime duration |
| `description` | `text` | Overtime description |
| `category` | `varchar(255)` | Overtime category |
| `created_by_user_id` | `uint` | Foreign key to users |
| `created_at` | `timestamp` | Creation timestamp |
| `updated_at` | `timestamp` | Update timestamp |

### Relationships
- `User` (belongs_to): Created by user
- `TelegramUser` (belongs_to): Associated telegram user

## Logging
Semua operasi dicatat menggunakan MyLogger dengan context:
- **Module**: "OvertimeManagement"
- **Action**: Method name (e.g., "CreateNewRecordOvertime")
- **Layer**: "repository", "service", "controller"
- **Level**: "debug", "info", "error"

## Testing
Test cases tersedia di `docs/rest-api/overtime.http` untuk testing manual menggunakan REST Client.

## Architecture
```
Controller -> Service -> Repository -> Database
     |          |           |
   Validation  Business   Data Access
   Logging     Logic      Logging
   HTTP        Logging    Transaction
```
