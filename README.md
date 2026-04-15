# TaskFlow — Backend (Go + PostgreSQL)

A production-style task management backend built using Go, PostgreSQL, and Docker.
It supports authentication, project management, and task tracking with clean architecture and secure APIs.

---

# 🚀 Tech Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Routing:** chi
* **Auth:** JWT + bcrypt
* **DB Driver:** pgx
* **Migrations:** golang-migrate
* **Containerization:** Docker + Docker Compose
* **Logging:** slog

---

# 🧠 Architecture

This project follows a **clean layered architecture**:

```text
Handler → Service → Repository → Database
```

### Layers

* **Handler**

  * HTTP handling, validation, responses
* **Service**

  * Business logic, authorization
* **Repository**

  * SQL queries using pgx
* **Models**

  * Structs mapped to DB schema

---

# 📁 Project Structure

```text
backend/
├── cmd/server
├── internal/
│   ├── config
│   ├── db
│   ├── models
│   ├── repository
│   ├── service
│   ├── handler
│   ├── middleware
│   └── utils
├── migrations
├── seed
├── Dockerfile
```

---

# ⚙️ Setup & Run

### 1. Clone repo

```bash
git clone <your-repo-url>
cd taskflow
```

---

### 2. Setup env

```bash
cp .env.example .env
```

---

### 3. Run everything

```bash
docker compose up --build
```

---

### 4. Server

```text
http://localhost:8080
```

---

# 🔁 Migrations

* Automatically run on startup
* Uses `golang-migrate`
* No manual step required

---

# 🌱 Seed Data

Seed runs automatically after migrations.

### Test credentials:

```text
Email:    test@example.com
Password: password123
```

---

# 🔐 Authentication

* JWT-based auth
* Token expiry: 24 hours

### Header:

```http
Authorization: Bearer <token>
```

---

# 📡 API Endpoints

---

## Auth

| Method | Endpoint         |
| ------ | ---------------- |
| POST   | `/auth/register` |
| POST   | `/auth/login`    |

---

## Projects

| Method | Endpoint                 |
| ------ | ------------------------ |
| GET    | `/projects?page=&limit=` |
| POST   | `/projects`              |
| GET    | `/projects/{projectID}`  |
| PATCH  | `/projects/{projectID}`  |
| DELETE | `/projects/{projectID}`  |

---

## Tasks

| Method | Endpoint                                                     |
| ------ | ------------------------------------------------------------ |
| GET    | `/projects/{projectID}/tasks?page=&limit=&status=&assignee=` |
| POST   | `/projects/{projectID}/tasks`                                |
| GET    | `/tasks/{taskID}`                                            |
| PATCH  | `/tasks/{taskID}`                                            |
| DELETE | `/tasks/{taskID}`                                            |

---

# 🔍 Pagination

```http
GET /projects?page=1&limit=10
GET /projects/{id}/tasks?page=2&limit=5
```

Defaults:

* page = 1
* limit = 10 (max 100)

---

# ⚠️ Error Handling

All responses are JSON:

```json
{
  "error": "message"
}
```

### Status Codes

| Code | Meaning        |
| ---- | -------------- |
| 400  | Bad request    |
| 401  | Unauthorized   |
| 403  | Forbidden      |
| 404  | Not found      |
| 500  | Internal error |

---

# 🐳 Docker

### Run full system:

```bash
docker compose up --build
```

Includes:

* PostgreSQL container
* Backend container
* Automatic migrations
* Automatic seed

---

# 🧠 Key Design Decisions

### 1. Clean Architecture

Improves maintainability and testability

### 2. UUID for IDs

Better for distributed systems and security

### 3. pgx (No ORM)

Full control over SQL and performance

### 4. Go-based Seeding

Used instead of SQL to support:

* bcrypt hashing
* UUID generation

### 5. chi Router

Lightweight and idiomatic Go routing

### 6. Static Binary Build

Ensures compatibility with Alpine Docker

---

# ⚖️ Tradeoffs

* No pagination metadata (total count) to keep queries simple
* Minimal validation (can be extended)
* No caching layer

---

# 🚧 Future Improvements

* Add total count in pagination
* Add unit + integration tests
* Add request validation library
* Add Swagger/OpenAPI docs
* Implement RBAC
* Add CI/CD pipeline

---

# ✅ Features

✔ Authentication (JWT + bcrypt)
✔ Project CRUD (owner-based auth)
✔ Task CRUD with filters
✔ Pagination support
✔ PostgreSQL + migrations
✔ Dockerized setup
✔ Seed data
✔ Structured logging
✔ Clean architecture

---

# 🙌 Notes

This project focuses on:

* correctness
* clarity
* maintainability

rather than unnecessary complexity.

---
