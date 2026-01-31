# SILAJU Backend - Feature Review

Backend API untuk aplikasi pelaporan kondisi jalan dengan sistem manajemen worker dan admin.

## 🏗️ Arsitektur

- **Framework**: Gin (Go)
- **Database**: PostgreSQL dengan GORM
- **Cloud Storage**: Cloudinary
- **Authentication**: JWT Bearer Token
- **Documentation**: Swagger UI

---

## 👥 Role-Based Access Control (RBAC)

| Role     | Deskripsi                                |
| -------- | ---------------------------------------- |
| `user`   | Pengguna umum yang dapat membuat laporan |
| `worker` | Petugas lapangan yang menangani laporan  |
| `admin`  | Administrator sistem                     |

---

## 🔐 Authentication API

### Public Endpoints

| Method | Endpoint                    | Deskripsi                |
| ------ | --------------------------- | ------------------------ |
| `POST` | `/api/auth/user/register`   | Registrasi user baru     |
| `POST` | `/api/auth/user/verify-otp` | Verifikasi OTP via email |
| `POST` | `/api/auth/user/login`      | Login user               |
| `POST` | `/api/auth/google`          | Login via Google OAuth   |
| `POST` | `/api/auth/admin/login`     | Login admin              |
| `POST` | `/api/auth/worker/login`    | Login worker             |

### Protected Endpoints

| Method | Endpoint                  | Role   | Deskripsi                 |
| ------ | ------------------------- | ------ | ------------------------- |
| `GET`  | `/api/auth/me`            | All    | Mendapatkan profil user   |
| `GET`  | `/api/auth/admin/users`   | Admin  | Mendapatkan semua users   |
| `GET`  | `/api/auth/admin/workers` | Admin  | Mendapatkan semua workers |
| `GET`  | `/api/auth/worker/me`     | Worker | Mendapatkan profil worker |
| `GET`  | `/api/auth/user/me`       | User   | Mendapatkan profil user   |

---

## 📝 Report API

### Public Endpoints

| Method | Endpoint          | Deskripsi                                                           |
| ------ | ----------------- | ------------------------------------------------------------------- |
| `GET`  | `/api/get_report` | Mendapatkan semua laporan yang sudah selesai dengan status non-good |

### User Endpoints

| Method | Endpoint              | Deskripsi                                                |
| ------ | --------------------- | -------------------------------------------------------- |
| `POST` | `/api/user/report`    | Membuat laporan baru (multipart/form-data dengan gambar) |
| `GET`  | `/api/user/report/me` | Mendapatkan semua laporan user (pagination)              |

### Worker Endpoints

| Method  | Endpoint                        | Deskripsi                                                              |
| ------- | ------------------------------- | ---------------------------------------------------------------------- |
| `PATCH` | `/api/worker/report`            | Menyelesaikan laporan dengan upload foto after                         |
| `GET`   | `/api/worker/report/assign/me`  | Mendapatkan laporan yang di-assign ke worker (pagination)              |
| `GET`   | `/api/worker/report/history/me` | Mendapatkan history laporan worker (pagination, filter `verify_admin`) |

### Admin Endpoints

| Method  | Endpoint                   | Deskripsi                                       |
| ------- | -------------------------- | ----------------------------------------------- |
| `PATCH` | `/api/admin/report/assign` | Assign worker ke laporan                        |
| `GET`   | `/api/admin/report/assign` | Mendapatkan semua laporan yang sudah di-assign  |
| `PATCH` | `/api/admin/report/verify` | Verifikasi laporan yang sudah dikerjakan worker |

---

## 📊 Report Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        REPORT LIFECYCLE                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. USER creates report          →  Status: "complete"           │
│                                                                 │
│  2. ADMIN assigns worker         →  Status: "assigned"         │
│     (+ admin_notes, deadline)                                   │
│                                                                 │
│  3. WORKER completes task        →  Status: "finish by worker" │
│     (+ uploads after image)                                     │
│                                                                 │
│  4. ADMIN verifies completion    →  Status: "finished"         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🗄️ Database Schema

### User Entity

| Field      | Type         | Deskripsi             |
| ---------- | ------------ | --------------------- |
| `id`       | UUID         | Primary key           |
| `username` | VARCHAR(100) | Username unik         |
| `fullname` | VARCHAR(100) | Nama lengkap          |
| `email`    | VARCHAR(100) | Email unik            |
| `role`     | VARCHAR(20)  | user/admin/worker     |
| `password` | VARCHAR(255) | Hashed password       |
| `verified` | BOOLEAN      | Status verifikasi OTP |

### Report Entity

| Field              | Type      | Deskripsi                                     |
| ------------------ | --------- | --------------------------------------------- |
| `id`               | TEXT      | Primary key (format: uuid_timestamp_long_lat) |
| `user_id`          | UUID      | FK ke User                                    |
| `worker_id`        | UUID      | FK ke Worker (nullable)                       |
| `longitude`        | NUMERIC   | Koordinat longitude                           |
| `latitude`         | NUMERIC   | Koordinat latitude                            |
| `road_name`        | TEXT      | Nama jalan                                    |
| `before_image_url` | TEXT      | URL gambar sebelum                            |
| `after_image_url`  | TEXT      | URL gambar sesudah                            |
| `description`      | TEXT      | Deskripsi laporan                             |
| `destruct_class`   | TEXT      | Klasifikasi kerusakan                         |
| `location_score`   | NUMERIC   | Skor lokasi                                   |
| `total_score`      | NUMERIC   | Skor total                                    |
| `status`           | TEXT      | Status laporan                                |
| `admin_notes`      | TEXT      | Catatan dari admin                            |
| `deadline`         | TIMESTAMP | Batas waktu pengerjaan                        |

---

## 🔧 Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Cloud Storage**: Cloudinary
- **Documentation**: Swagger (swaggo)
- **Deployment**: Docker, Hugging Face Spaces

---

## 📁 Project Structure

```
BE-MainServer/
├── config/          # Konfigurasi database, secret
├── controllers/     # Handler HTTP request
├── docs/            # Swagger documentation
├── middleware/      # Auth, RBAC middleware
├── models/
│   ├── dto/         # Data Transfer Objects
│   ├── entity/      # Database entities + constants
│   └── error/       # Centralized error definitions
├── provider/        # Dependency injection
├── repositories/    # Database operations
├── router/          # Route definitions
├── services/        # Business logic
├── utils/           # Helper functions
├── main.go          # Entry point
└── Dockerfile       # Container config
```

---

## 🚀 Quick Start

```bash
# Install dependencies
go mod tidy

# Generate Swagger docs
swag init

# Run development server
go run main.go

# Build for production
go build -o main main.go
```

---

## 📖 API Documentation

Swagger UI tersedia di: `http://localhost:8080/swagger/index.html`

---

## 🔒 Security Features

- JWT Bearer Token Authentication
- Role-Based Access Control (RBAC)
- Password Hashing (bcrypt)
- OTP Email Verification
- Google OAuth Integration
- Gzip Compression

---
title: DINACOM 11.0 Backend
emoji: 🤖
colorFrom: green
colorTo: blue
sdk: docker
pinned: false
---
