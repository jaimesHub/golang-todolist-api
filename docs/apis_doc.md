# API Documentation for Golang TODO List Application

## Authentication

### Register a new user
- **URL**: `/api/v1/auth/register`
- **Method**: `POST`
- **Auth required**: No
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }
  ```
- **Success Response**:
  - **Code**: 201 Created
  - **Content**:
    ```json
    {
      "message": "User registered successfully",
      "user": {
        "id": "uuid-string",
        "email": "user@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "created_at": "2025-04-11T16:00:00Z"
      }
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request
  - **Content**:
    ```json
    {
      "error": "User with this email already exists"
    }
    ```

### Login
- **URL**: `/api/v1/auth/login`
- **Method**: `POST`
- **Auth required**: No
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Login successful",
      "token": "jwt-token-string",
      "user": {
        "id": "uuid-string",
        "email": "user@example.com",
        "first_name": "John",
        "last_name": "Doe"
      }
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**:
    ```json
    {
      "error": "Invalid email or password"
    }
    ```

### Refresh Token
- **URL**: `/api/v1/auth/refresh`
- **Method**: `POST`
- **Auth required**: No
- **Request Body**:
  ```json
  {
    "token": "current-jwt-token"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Token refreshed successfully",
      "token": "new-jwt-token-string"
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**:
    ```json
    {
      "error": "Invalid or expired token"
    }
    ```

## User Management

### Get User Profile
- **URL**: `/api/v1/users/me`
- **Method**: `GET`
- **Auth required**: Yes (JWT token in Authorization header)
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "user": {
        "id": "uuid-string",
        "email": "user@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "is_active": true,
        "created_at": "2025-04-11T16:00:00Z",
        "updated_at": "2025-04-11T16:00:00Z"
      }
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**:
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### Update User Profile
- **URL**: `/api/v1/users/me`
- **Method**: `PUT`
- **Auth required**: Yes (JWT token in Authorization header)
- **Request Body**:
  ```json
  {
    "first_name": "John",
    "last_name": "Smith"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Profile updated successfully",
      "user": {
        "id": "uuid-string",
        "email": "user@example.com",
        "first_name": "John",
        "last_name": "Smith",
        "updated_at": "2025-04-11T16:30:00Z"
      }
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**:
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### Get User Activities
- **URL**: `/api/v1/users/activities?limit=10&offset=0`
- **Method**: `GET`
- **Auth required**: Yes (JWT token in Authorization header)
- **Query Parameters**:
  - `limit` (optional): Number of activities to return (default: 10)
  - `offset` (optional): Offset for pagination (default: 0)
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "activities": [
        {
          "id": "uuid-string",
          "user_id": "uuid-string",
          "action": "create",
          "entity": "task",
          "entity_id": "uuid-string",
          "details": "Task created: Example Task",
          "created_at": "2025-04-11T16:30:00Z"
        }
      ],
      "pagination": {
        "limit": 10,
        "offset": 0,
        "count": 1
      }
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**:
    ```json
    {
      "error": "Unauthorized"
    }
    ```

## Task Management

### Create Task
- **URL**: `/api/v1/tasks`
- **Method**: `POST`
- **Auth required**: Yes (JWT token in Authorization header)
- **Request Body**:
  ```json
  {
    "title": "Example Task",
    "description": "This is an example task",
    "priority": 1,
    "due_date": "2025-04-15T16:00:00Z"
  }
  ```
- **Success Response**:
  - **Code**: 201 Created
  - **Content**:
    ```json
    {
      "message": "Task created successfully",
      "task": {
        "id": "uuid-string",
        "title": "Example Task",
        "description": "This is an example task",
        "status": "pending",
        "priority": 1,
        "due_date": "2025-04-15T16:00:00Z",
        "user_id": "uuid-string",
        "created_at": "2025-04-11T16:30:00Z",
        "updated_at": "2025-04-11T16:30:00Z"
      }
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**:
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### List Tasks
- **URL**: `/api/v1/tasks?status=pending&priority=1&limit=10&offset=0`
- **Method**: `GET`
- **Auth required**: Yes (JWT token in Authorization header)
- **Query Parameters**:
  - `status` (optional): Filter by status (pending, in_progress, completed)
  - `priority` (optional): Filter by priority (0: low, 1: medium, 2: high)
  - `limit` (optional): Number of tasks to return (default: 10)
  - `offset` (optional): Offset for pagination (default: 0)
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "tasks": [
        {
          "id": "uuid-string",
          "title": "Example Task",
          "description": "This is an example task",
          "status": "pending",
          "priority": 1,
          "due_date": "2025-04-15T16:00:00Z",
          "user_id": "uuid-string",
          "created_at": "2025-04-11T16:30:00Z",
          "updated_at": "2025-04-11T16:30:00Z"
        }
      ],
      "pagination": {
        "total": 1,
        "limit": 10,
        "offset": 0
      }
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**:
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### Get Task by ID
- **URL**: `/api/v1/tasks/:id`
- **Method**: `GET`
- **Auth required**: Yes (JWT token in Authorization header)
- **URL Parameters**:
  - `id`: UUID of the task
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "task": {
        "id": "uuid-string",
        "title": "Example Task",
        "description": "This is an example task",
        "status": "pending",
        "priority": 1,
        "due_date": "2025-04-15T16:00:00Z",
        "user_id": "uuid-string",
        "created_at": "2025-04-11T16:30:00Z",
        "updated_at": "2025-04-11T16:30:00Z"
      }
    }
    ```
- **Error Response**:
  - **Code**: 404 Not Found
  - **Content**:
    ```json
    {
      "error": "Task not found"
    }
    ```

### Update Task
- **URL**: `/api/v1/tasks/:id`
- **Method**: `PUT`
- **Auth required**: Yes (JWT token in Authorization header)
- **URL Parameters**:
  - `id`: UUID of the task
- **Request Body**:
  ```json
  {
    "title": "Updated Task",
    "description": "This is an updated task",
    "status": "in_progress",
    "priority": 2,
    "due_date": "2025-04-16T16:00:00Z"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Task updated successfully",
      "task": {
        "id": "uuid-string",
        "title": "Updated Task",
        "description": "This is an updated task",
        "status": "in_progress",
        "priority": 2,
        "due_date": "2025-04-16T16:00:00Z",
        "user_id": "uuid-string",
        "created_at": "2025-04-11T16:30:00Z",
        "updated_at": "2025-04-11T16:45:00Z"
      }
    }
    ```
- **Error Response**:
  - **Code**: 404 Not Found
  - **Content**:
    ```json
    {
      "error": "Task not found"
    }
    ```

### Delete Task
- **URL**: `/api/v1/tasks/:id`
- **Method**: `DELETE`
- **Auth required**: Yes (JWT token in Authorization header)
- **URL Parameters**:
  - `id`: UUID of the task
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "message": "Task deleted successfully"
    }
    ```
- **Error Response**:
  - **Code**: 404 Not Found
  - **Content**:
    ```json
    {
      "error": "Task not found"
    }
    ```

## Health Check and Monitoring

### Health Check
- **URL**: `/health`
- **Method**: `GET`
- **Auth required**: No
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "status": "ok",
      "message": "Service is running",
      "time": "2025-04-11T16:45:00Z"
    }
    ```

### Metrics
- **URL**: `/metrics`
- **Method**: `GET`
- **Auth required**: No
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "status": "ok",
      "message": "Metrics endpoint"
    }
    ```
