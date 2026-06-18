# Short Link Backend

Backend service untuk aplikasi short link built dengan Go, PostgreSQL, Redis, dan Docker (WSL Native).

---

## Tech Stack

- **Backend**: Go 1.26+
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL
- **Cache**: Redis (JWT Blacklist & caching)
- **Migration Tool**: golang-migrate
- **Containerization**: Docker (WSL Native)
- **Authentication**: JWT (Redis Blacklist)

---

## Struktur Proyek

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

## Prasyarat

- Go 1.26+
- WSL2 (Ubuntu recommended)
- Docker (install langsung di WSL, bukan Docker Desktop)
- Docker Compose (via WSL)

---

## Setup Project

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

### 3. Jalankan Docker (WSL)

```bash
docker compose up -d
```

Services:

- PostgreSQL → `localhost:5432`
- Redis → `localhost:6379`

---

### 4. Install Dependencies

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

Menggunakan JWT + Redis Blacklist

## Authorization Header

```http
Authorization: Bearer <token>
```

---

## Auth Endpoints

### Register

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

### Login

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

### Logout (JWT Blacklist via Redis)

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

### Behavior

- Token dimasukkan ke Redis blacklist
- Middleware akan memblokir token yang sudah logout
- Token tetap valid sampai expired, tetapi tidak bisa digunakan lagi

---

# 🔗 Links API

---

## Create Short Link

```http
POST /api/links
```

---

## Get All Links

```http
GET /api/links?page=1&limit=10
```

---

## ❌ Get Deleted Links

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

- Mengambil data soft-deleted links
- Tidak muncul di list aktif
- Bisa digunakan untuk audit atau future restore feature

---

## Delete Link

```http
DELETE /api/links/:id
```

---

## Check Slug Availability

```http
GET /api/links/check-slug/:slug
```

---

# 🔄 Redirect

```http
GET /:slug
```

### Behavior

- 301 redirect ke original URL
- 404 jika slug tidak ditemukan

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

# 🧠 Notes

- Redis digunakan untuk:
  - JWT blacklist (logout system)
  - caching (optional improvement)

- Soft delete digunakan untuk deleted links
- Full project berjalan di WSL + Docker native
- Clean architecture: controller → service → repository

```

```
