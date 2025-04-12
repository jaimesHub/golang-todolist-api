# Setup and Deployment Guide

## Prerequisites

- Docker and Docker Compose
- Git (for cloning the repository)

## Local Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/user/golang-todo-app.git
cd golang-todo-app
```

### 2. Configure Environment Variables

Copy the example environment file and modify it as needed:

```bash
cp .env.example .env
```

Edit the `.env` file to set your configuration:

```
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo_app
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24

# Logging Configuration
LOG_LEVEL=info
LOG_FILE=logs/app.log
```

### 3. Run with Docker Compose

The easiest way to run the application is using Docker Compose:

```bash
docker-compose up -d
```

This will start the following services:
- API server on port 8080
- PostgreSQL database on port 5432
- Redis on port 6379

### 4. Run Database Migrations

The migrations should run automatically when the application starts, but you can also run them manually:

```bash
docker-compose exec app ./scripts/run_migrations.sh
```

### 5. Access the API

The API will be available at http://localhost:8080

## Manual Setup (Without Docker)

### 1. Install Dependencies

Make sure you have Go 1.21 or later installed:

```bash
go version
```

Install PostgreSQL and Redis:

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install postgresql redis-server

# macOS with Homebrew
brew install postgresql redis
```

### 2. Configure Environment Variables

Copy and configure the `.env` file as described above.

### 3. Install Go Dependencies

```bash
go mod download
```

### 4. Run Database Migrations

```bash
./scripts/run_migrations.sh
```

### 5. Build and Run the Application

```bash
go build -o app ./cmd/api
./app
```

## Production Deployment

For production deployment, consider the following additional steps:

### 1. Use a Secure JWT Secret

Generate a secure random string for the JWT secret:

```bash
openssl rand -hex 32
```

Update the `JWT_SECRET` in your environment variables.

### 2. Configure HTTPS

In production, you should use HTTPS. You can:
- Use a reverse proxy like Nginx with Let's Encrypt certificates
- Configure your cloud provider's load balancer with SSL/TLS

### 3. Set Up Database Backups

Configure regular backups for your PostgreSQL database.

### 4. Monitoring and Logging

Consider setting up:
- Prometheus and Grafana for monitoring
- ELK stack or similar for centralized logging

### 5. Scaling

For higher load:
- Scale the API server horizontally behind a load balancer
- Consider using PostgreSQL read replicas
- Set up Redis clustering

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Check that PostgreSQL is running
   - Verify database credentials in `.env`
   - Ensure the database exists

2. **Redis Connection Errors**
   - Check that Redis is running
   - Verify Redis connection details in `.env`

3. **API Server Won't Start**
   - Check logs for errors
   - Verify port 8080 is not in use
   - Ensure all environment variables are set correctly

### Viewing Logs

With Docker Compose:
```bash
docker-compose logs -f app
```

Without Docker:
```bash
tail -f logs/app.log
```

## Updating the Application

### With Docker Compose

```bash
git pull
docker-compose down
docker-compose build
docker-compose up -d
```

### Without Docker

```bash
git pull
go mod download
go build -o app ./cmd/api
./app
```
