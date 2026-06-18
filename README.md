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
  "password": "password123"
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
  "token": "jwt_token_here"
}
```

---

### 🚪 Logout (JWT Blacklist via Redis)

```http
POST /api/auth/logout
```

Headers:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "status": "success",
  "message": "success logout"
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

---

## 📄 Get All Links

```http
GET /api/links?page=1&limit=10
```

---

## 🗑️ Get Deleted Links

```http
GET /api/links/deleted
```

Headers:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "status": "success",
  "message": "success get deleted links",
  "data": {
    "links": []
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

---

## 🔎 Check Slug Availability

```http
GET /api/links/check-slug/:slug
```

---

# 🔄 Redirect System

```http
GET /:slug
```

### Behavior

- Redirect 301 ke original URL
- Return 404 jika slug tidak ditemukan

---

# 👤 Profile

## Get Profile

```http
GET /api/profile
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
- Soft delete untuk links (deleted_links table / flag)

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
