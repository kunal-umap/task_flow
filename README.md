# TaskFlow — Backend (Go + PostgreSQL)

A minimal but production-style task management backend built with Go, PostgreSQL, and Docker.

This service supports user authentication, project management, and task tracking with proper authorization, clean architecture, and containerized deployment.

---

# 🚀 Tech Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Auth:** JWT + bcrypt
* **Routing:** chi
* **Database Driver:** pgx
* **Migrations:** golang-migrate
* **Containerization:** Docker + Docker Compose

---

# 🧠 Architecture Overview

This project follows a **clean layered architecture**:

```
Handler → Service → Repository → Database
```

### Layers

* **Handler (HTTP Layer)**
  Handles request/response, validation, and status codes

* **Service (Business Logic)**
  Contains core logic, validation, and authorization

* **Repository (Data Access)**
  Handles SQL queries using pgx

* **Models**
  Shared structs representing database entities

---

# 🗂️ Project Structure

```
backend/
├── cmd/server          # Entry point
├── internal/
│   ├── config         # Env config
│   ├── db             # DB connection
│   ├── models         # Data models
│   ├── repository     # DB queries
│   ├── service        # Business logic
│   ├── handler        # HTTP handlers
│   ├── middleware     # JWT auth middleware
│   └── utils          # bcrypt + JWT
├── migrations         # SQL migrations
├── seed               # Seed script (Go)
├── Dockerfile
```

---

# ⚙️ Running Locally

### 1. Clone repo

```bash
git clone <your-repo-url>
cd taskflow
```

---

### 2. Setup environment

```bash
cp .env.example .env
```

---

### 3. Run everything

```bash
docker compose up --build
```

---

### 4. API available at

```
http://localhost:8080
```

---

# 🔁 Migrations

* Migrations are automatically executed on startup using `golang-migrate`
* No manual step required

---

# 🌱 Seed Data

Seed runs automatically on startup.

### Test credentials:

```
Email:    test@example.com
Password: password123
```

---

# 🔐 Authentication

* JWT-based authentication
* Token expiry: 24 hours
* Protected routes require:

```
Authorization: Bearer <token>
```

---

# 📡 API Endpoints

## Auth

| Method | Endpoint         |
| ------ | ---------------- |
| POST   | `/auth/register` |
| POST   | `/auth/login`    |

---

## Projects

| Method | Endpoint                |
| ------ | ----------------------- |
| GET    | `/projects`             |
| POST   | `/projects`             |
| GET    | `/projects/{projectID}` |
| PATCH  | `/projects/{projectID}` |
| DELETE | `/projects/{projectID}` |

---

## Tasks

| Method | Endpoint                      |
| ------ | ----------------------------- |
| GET    | `/projects/{projectID}/tasks` |
| POST   | `/projects/{projectID}/tasks` |
| GET    | `/tasks/{taskID}`             |
| PATCH  | `/tasks/{taskID}`             |
| DELETE | `/tasks/{taskID}`             |

---

## Filters

```
GET /projects/{projectID}/tasks?status=todo
GET /projects/{projectID}/tasks?assignee=<uuid>
```

---

# ⚠️ Error Handling

Standardized responses:

```json
{
  "error": "validation failed"
}
```

| Code | Meaning      |
| ---- | ------------ |
| 400  | Bad request  |
| 401  | Unauthorized |
| 403  | Forbidden    |
| 404  | Not found    |

---

# 🐳 Docker Setup

* Multi-stage Dockerfile for backend
* PostgreSQL runs in container
* Everything starts with one command:

```bash
docker compose up
```

---

# 🧠 Key Design Decisions

### 1. Layered Architecture

Separates concerns for maintainability and testability

### 2. UUID Usage

Used UUIDs for all primary keys for scalability and distributed safety

### 3. pgx instead of ORM

Avoided ORM to maintain control over queries and performance

### 4. Go-based Seeding

Used Go seed script instead of raw SQL for:

* bcrypt hashing
* UUID control
* flexibility

### 5. chi Router

Lightweight and idiomatic routing with middleware support

---

# ⚖️ Tradeoffs

* Used query params initially before switching to chi (incremental approach)
* No pagination implemented to keep scope focused
* Minimal validation (could be expanded using validator library)

---

# 🚧 What I Would Improve With More Time

* Add pagination (`?page=&limit=`)
* Add unit + integration tests
* Implement structured logging (slog)
* Add request validation library
* Improve error response consistency
* Add role-based access (RBAC)
* Add Swagger/OpenAPI docs
* Implement CI/CD pipeline

---

# ✅ Features Completed

✔ Authentication (JWT + bcrypt)
✔ Project CRUD with authorization
✔ Task CRUD with filters
✔ PostgreSQL with migrations
✔ Dockerized environment
✔ Seed data
✔ Clean architecture

---

# 🙌 Final Notes

This project focuses on:

* correctness
* clarity
* maintainability

rather than over-engineering.

---
