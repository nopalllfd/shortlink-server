---
# 🚀 Short Link Backend

Backend service untuk aplikasi short link built dengan Go, PostgreSQL, Redis, dan Docker (WSL2 Native Ready).
---

## ⚙️ Tech Stack

- **Backend**: Go 1.26+
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL
- **Cache**: Redis (JWT Blacklist & caching)
- **Migration Tool**: golang-migrate
- **Containerization**: Docker (WSL2 Native / Docker Desktop WSL2 Integration)
- **Authentication**: JWT (Redis Blacklist)

---

## 📁 Struktur Proyek

```
server/
├── cmd/
│   └── main.go
├── db/
│   └── migrations/
├── internal/
│   ├── config/
│   ├── controller/
│   ├── dto/
│   ├── errs/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── response/
│   ├── route/
│   └── service/
├── pkg/
├── public/
├── .env.example
├── docker-compose.yml
├── Makefile
├── go.mod
└── go.sum
```

---

## 🧰 Prasyarat

Pastikan environment berikut sudah terinstall:

- Go 1.26+
- WSL2 (Ubuntu recommended)
- Docker (WSL Native atau Docker Desktop dengan WSL2 integration)
- Docker Compose

---

## 🚀 Setup Project

### 1. Clone Repository

```bash
git clone https://github.com/nopalllfd/shortlink-server.git
cd shortlink-server
```

---

### 2. Setup Environment

```bash
cp .env.example .env
```

---

### 3. Jalankan Docker (WSL2)

Pastikan Docker daemon aktif di WSL:

```bash
sudo service docker start
```

Jika pakai systemd:

```bash
sudo systemctl start docker
```

Lalu jalankan container:

```bash
docker compose up -d
```

---

### 📦 Services yang berjalan

- PostgreSQL → `localhost:5432`
- Redis → `localhost:6379`

---

### 4. Install Dependencies Go

```bash
go mod download
```

---

### 5. Jalankan Migration

```bash
make migrate-up
```

---

### 6. Run Server

```bash
go run cmd/main.go
```

Server berjalan di:

```
http://localhost:8080
```

---

# 🔐 Authentication

Menggunakan JWT + Redis Blacklist system

---

## Authorization Header

```http
Authorization: Bearer <token>
```

---

## Auth Endpoints

### 📝 Register

```http
POST /api/auth/register
```

Request:

```json
{
  "email": "user@example.com",
  "password": "password123" // min 8 chars
}
```

Response:

```json
{
  "status": "success",
  "message": "registered success",
  "data": null
}
```

---

### 🔑 Login

```http
POST /api/auth/login
```

Request:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "status": "success",
  "message": "login success",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com"
    },
    "token": "jwt_token_here"
  }
}
```

---

### 🚪 Logout (JWT Blacklist via Redis)

```http
DELETE /api/auth/logout
```

Headers:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "status": "success",
  "message": "logout success",
  "data": null
}
```

### ⚙️ Behavior

- Token dimasukkan ke Redis blacklist
- Middleware akan mengecek blacklist setiap request
- Token tetap valid sampai expired, tapi tidak bisa dipakai lagi

---

# 🔗 Links API

---

## ➕ Create Short Link

```http
POST /api/links
```

Headers:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "link": "https://example.com/very-long-url",
  "slug": "my-custom-slug" // optional, min 6 chars
}
```

Response:

```json
{
  "status": "success",
  "message": "success create link",
  "data": {
    "id": 1,
    "slug": "my-custom-slug",
    "original_url": "https://example.com/very-long-url",
    "short_link": "http://localhost:8080/my-custom-slug",
    "clicks": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 📄 Get All Links

```http
GET /api/links?page=1&limit=10
```

Headers:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "status": "success",
  "message": "success get all links",
  "data": {
    "links": [
      {
        "id": 1,
        "slug": "my-custom-slug",
        "original_url": "https://example.com/very-long-url",
        "short_link": "http://localhost:8080/my-custom-slug",
        "clicks": 5,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "total": 1,
      "total_pages": 1,
      "limit": 10,
      "next_link": null,
      "prev_link": null
    }
  }
}
```

---

## 🗑️ Get Deleted Links

```http
GET /api/links/deleted?page=1&limit=10
```

Headers:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "status": "success",
  "message": "success get all deleted links",
  "data": {
    "links": [
      {
        "id": 2,
        "slug": "old-slug",
        "original_url": "https://example.com/old-url",
        "short_link": "http://localhost:8080/old-slug",
        "clicks": 10,
        "created_at": "2024-01-02T00:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "total": 1,
      "total_pages": 1,
      "limit": 10,
      "next_link": null,
      "prev_link": null
    }
  }
}
```

### Behavior

- Mengambil soft-deleted links
- Tidak muncul di list aktif
- Bisa digunakan untuk audit / restore feature di masa depan

---

## ❌ Delete Link

```http
DELETE /api/links/:id
```

Headers:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "status": "success",
  "message": "success delete link",
  "data": null
}
```

---

## 🔎 Check Slug Availability

```http
GET /api/links/check-slug/:slug
```

Response:

```json
{
  "status": "success",
  "message": "slug succcesfully check",
  "data": false // true jika sudah ada
}
```

---

# 🔄 Redirect System

```http
GET /:slug
```

### Behavior

- Redirect 301 ke original URL
- Return 404 jika slug tidak ditemukan
- **Click Counter**: Setiap redirect akan menambah kolom `clicks` di database secara async (tidak block redirect user)

---

# 👤 Profile

## Get Profile

```http
GET /api/profiles
```

Headers:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "status": "success",
  "message": "success to get profile",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

# ⚙️ Make Commands

```bash
make migrate-create NAME=table_name
make migrate-up
make migrate-down
make migrate-force VERSION=1
```

---

# 🧠 Notes (IMPORTANT)

## 🐧 WSL2 Environment

Project ini berjalan di WSL2 Linux environment.

### Pastikan:

- Docker daemon aktif di WSL
- PostgreSQL & Redis diakses via `localhost`
- Tidak perlu konfigurasi network manual

### Cek Docker:

```bash
sudo service docker status
```

Restart:

```bash
sudo service docker restart
```

---

## 🧠 Architecture

- Controller → Handle HTTP request
- Service → Business logic
- Repository → Database layer
- Middleware → Auth & validation

---

## ⚡ Redis Usage

- JWT Blacklist (logout system)
- Optional caching layer untuk optimasi

---

## 🗃️ Database

- Menggunakan PostgreSQL
- Soft delete untuk links (flag `deleted_at`)
- Kolom `clicks` untuk tracking berapa kali link diakses

---

# 🚀 Final Result

```
Backend Shortlink System
✔ JWT Auth + Redis Blacklist
✔ Short URL Generator
✔ Soft Delete System
✔ Clean Architecture (Go)
✔ Dockerized WSL2 Ready
```
