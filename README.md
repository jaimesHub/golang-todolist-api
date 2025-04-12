# Golang TODO List Application

A simple, yet powerful TODO List application built with Golang, featuring user management, authentication/authorization, and task management.

## Features

- User Management
  - Registration and authentication
  - Profile management
  - Activity logging

- Task Management
  - Create, read, update, delete tasks
  - Task status tracking
  - Task prioritization

- Authentication/Authorization
  - JWT-based authentication
  - Role-based access control

- Technical Features
  - PostgreSQL database with Supabase
  - GORM ORM for database operations
  - Redis for queue processing
  - Structured logging and monitoring
  - Containerized with Docker and Docker Compose

## Project Structure

```
golang-todo-app/
├── cmd/
│   └── api/
│       └── main.go         # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── database/           # Database connection and migrations
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Database models
│   ├── routes/             # API route definitions
│   └── services/           # Business logic services
│       └── redis/          # Redis client and operations
├── pkg/                    # Reusable packages
├── migrations/             # Database migration files
├── config/                 # Configuration files
├── scripts/                # Utility scripts
└── docs/                   # Documentation
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Redis
- Docker and Docker Compose (for containerized deployment)

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo_app
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24

# Logging
LOG_LEVEL=info
LOG_FILE=
```

### Running Locally

1. Clone the repository
2. Install dependencies: `go mod download`
3. Run the application: `go run cmd/api/main.go`

### Running with Docker

```bash
docker-compose up -d
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Users

- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user profile
- `GET /api/v1/users/activities` - Get user activities

### Tasks

- `POST /api/v1/tasks` - Create a new task
- `GET /api/v1/tasks` - List all tasks
- `GET /api/v1/tasks/:id` - Get task by ID
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task

## License

This project is licensed under the MIT License.
