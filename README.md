# ABC User Management System

A full-stack user management application built with Vue.js and Golang, featuring event-driven architecture, Redis caching, and comprehensive CRUD operations.

## Features

- ✅ Complete user CRUD operations (Create, Read, Update, Delete)
- ✅ JWT-based authentication
- ✅ Admin default credentials on login page
- ✅ Search and pagination functionality
- ✅ Event-driven architecture with RabbitMQ
- ✅ Redis caching for improved performance
- ✅ Business rule validation (unique emails, age > 18)
- ✅ Clean architecture with repository pattern
- ✅ Dockerized deployment
- ✅ Unit tests for business logic

## Tech Stack

### Frontend
- Vue 3 with Composition API
- Element Plus UI Framework
- Vee-Validate for form validation
- Vuex for state management
- Vue Router for navigation
- Axios for API communication

### Backend
- Golang with Gin framework
- GORM for database ORM
- MySQL/MongoDB support
- Redis for caching
- RabbitMQ for event-driven architecture
- JWT for authentication
- Logrus for structured logging

---

## 🐳 Dockerized Deployment Guide

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop) and [Docker Compose](https://docs.docker.com/compose/) installed on your machine.

### 1. Clone the Repository

```bash
git clone <repository-url>
cd abc-user-management
```

### 2. Environment Configuration

- Copy `.env.example` to `.env` in the `backend/` directory and adjust values if needed (especially database passwords).

### 3. Build and Start All Services

```bash
docker-compose up -d --build
```

This will:
- Build the Go backend and Vue frontend images.
- Start MySQL, Redis, RabbitMQ, backend, and frontend containers.
- Run database migrations automatically.

### 4. Access the Application

- **Frontend:** [http://localhost:3000](http://localhost:3000)
- **Backend API:** [http://localhost:8085](http://localhost:8085)
- **RabbitMQ Management:** [http://localhost:15672](http://localhost:15672) (if enabled)

### 5. Stopping Services

```bash
docker-compose down
```

### 6. Logs & Troubleshooting

- View backend logs: `docker logs -f abc_backend`
- View MySQL logs: `docker logs -f abc_mysql`
- View frontend logs: `docker logs -f abc_frontend`

### 7. Running Tests

- Unit and integration tests can be run inside the backend container or locally with Go.
Run Everything (Unit + Integration):
backend 
NOTE: Mske sure to CD into backend

# Option 1: Using the script
its need to give run-tests.sh permission 
```bash
chmod +x run-tests.sh
./run-tests.sh
```
# Option 2: Using Make
```bash
make test-all
```
# Option 3: Just run make (default is test-all)
```bash
make test
Run Specific Test Types:
bash# Only unit tests
make test-unit
```

# Only integration tests
```bash
make test-integration
```

# With coverage report
```bash
make test-coverage
```

# Benchmark tests
```bash
make test-benchmark
```

# With race detection
```bash
make test-race
For CI/CD:
bash# Simple CI test with coverage threshold
chmod +x test-ci.sh
./test-ci.sh
```