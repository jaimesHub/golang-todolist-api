# System Architecture Documentation

## Overview

The Golang TODO List application is a modern, containerized web application built with a microservices-inspired architecture. It provides user management, authentication, task management, and activity logging functionality with a focus on performance, security, and scalability.

## Architecture Diagram

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│  Client         │────▶│  API Server     │────▶│  PostgreSQL     │
│  (Web/Mobile)   │     │  (Golang)       │     │  Database       │
│                 │◀────│                 │◀────│                 │
└─────────────────┘     └────────┬────────┘     └─────────────────┘
                                 │
                                 │
                        ┌────────▼────────┐
                        │                 │
                        │  Redis          │
                        │  (Queue/Cache)  │
                        │                 │
                        └─────────────────┘
```

## Components

### 1. API Server (Golang)

The core of the application is a RESTful API server built with Golang and the Gin web framework. It handles HTTP requests, processes business logic, and interacts with the database and Redis.

Key components:
- **Handlers**: Process HTTP requests and responses
- **Services**: Implement business logic
- **Middleware**: Handle cross-cutting concerns like authentication, logging, and CORS
- **Models**: Define data structures and database schema
- **Config**: Manage application configuration
- **Logger**: Provide structured logging
- **Monitoring**: Implement health checks and metrics

### 2. PostgreSQL Database

PostgreSQL serves as the primary data store for the application, storing user information, tasks, and activity logs.

Key features:
- **GORM ORM**: Used for database operations and migrations
- **Supabase Integration**: Provides additional database management capabilities
- **Migrations**: Automated schema management

### 3. Redis

Redis is used for two primary purposes:
- **Queue Processing**: Handling asynchronous tasks
- **Caching**: (Potential future use)

### 4. Docker Containerization

The application is containerized using Docker and orchestrated with Docker Compose, making it easy to deploy and scale.

Components:
- **Application Container**: Runs the Golang API server
- **PostgreSQL Container**: Runs the database
- **Redis Container**: Runs Redis for queue processing

## Data Flow

1. **Authentication Flow**:
   - Client sends credentials to `/api/v1/auth/login`
   - Server validates credentials and generates JWT token
   - Client stores token and includes it in subsequent requests
   - Server validates token for protected endpoints

2. **Task Management Flow**:
   - Authenticated users can create, read, update, and delete tasks
   - Each operation is logged in the activity log
   - Task operations trigger events that can be processed asynchronously

3. **Queue Processing Flow**:
   - Events are added to Redis queues
   - Worker processes consume events from queues
   - Workers perform actions like sending notifications or updating statistics

## Security

1. **Authentication**: JWT-based authentication with token expiration and refresh
2. **Authorization**: Role-based access control for protected resources
3. **Password Security**: Passwords are hashed using bcrypt
4. **HTTPS**: All communication should be over HTTPS in production
5. **CORS**: Configured to allow specific origins

## Scalability

The application is designed to be scalable:

1. **Horizontal Scaling**: API servers can be scaled horizontally behind a load balancer
2. **Database Scaling**: PostgreSQL can be scaled with read replicas
3. **Queue Scaling**: Redis can be clustered for higher throughput
4. **Containerization**: Docker makes it easy to deploy multiple instances

## Monitoring and Logging

1. **Structured Logging**: JSON-formatted logs with contextual information
2. **Health Checks**: Endpoints to verify service health
3. **Metrics**: Basic application metrics for monitoring
4. **Activity Tracking**: User activities are logged for audit purposes

## Deployment

The application can be deployed using Docker Compose:

```bash
docker-compose up -d
```

This will start the API server, PostgreSQL, and Redis containers with the appropriate configuration.

## Future Enhancements

1. **Caching Layer**: Implement Redis caching for frequently accessed data
2. **Full-text Search**: Add search capabilities for tasks
3. **Notifications**: Implement email and push notifications
4. **API Rate Limiting**: Protect against abuse
5. **Advanced Monitoring**: Integrate with Prometheus and Grafana
