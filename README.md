# STUB PROJECT for Go


---
# Task Management System (Go, Gin, PostgreSQL)

A simple microservice for managing tasks, built in Go using the GinGonic framework and PostgreSQL. Supports CRUD operations, pagination, and status-based filtering.

## Problem Breakdown & Design Decisions

- **Goal**: Build a Task Management System with create, read, update, delete, pagination, and filtering features.
- **Microservice Architecture**: Clear separation of concerns—controllers, repositories, models, and drivers.
- **Single Responsibility Principle**: Each package (controller, repo, model, driver) handles a distinct responsibility.
- **API Design**: RESTful endpoints for all operations.
- **Scalability**: Stateless service; can be scaled horizontally by running multiple instances behind a load balancer.
- **Database**: PostgreSQL for persistent storage, using pgx driver for database interactions. Use row level locking and transactions to handle concurrent updates.
- **Pagination**: Implemented using cursor-based pagination for efficient data retrieval.
- **Soft Delete**: Tasks are marked as deleted instead of being removed from the database, allowing for recovery and audit trails.
- **Extensibility**: Easy to add new microservices with REST for communication.

## Frameworks & Tools
- Go (with Go modules)
- GinGonic (REST API framework)
- PostgreSQL (database)
- pgx (Postgres driver)
- golang-migrate (database migrations)

## Setup
- Clone the repository:
 https://github.com/ankeshkmr2010/GoStubproject/tree/ankesh/alleDemo
- Make sure to switch to the correct branch: ankesh/alleDemo
- Have postgres running on your machine or use a Docker container (docker compose file provided if needed).
- Configure environment variables:
  - Edit .env.local with your DB credentials and rename to .env.
- Start the service:
  - ``` go run main.go ```
  - (Migrations run automatically on startup.)
    The API will be available at http://localhost:8080.

  
## API Endpoints

### Create Task

- **Endpoint:** `POST /task/create`
- **Request:**
    ```json
    {
      "name": "Sample Task",
      "description": "Details",
      "status": "Pending",
      "priority": 1,
      "created_by": "uuid",
      "task_data": {},
      "requested_at": "2024-06-01T12:00:00Z"
    }
    ```
- **Response:**
    ```json
    {
      "id": "uuid",
      "name": "...",
      "description": "...",
      "status": "Pending",
      "priority": 1,
      "created_by": "uuid",
      "created_at": "...",
      "updated_at": "...",
      "serial_number": 1,
      "requested_at": "...",
      "is_deleted": false
    }
    ```

---

### Get Task by ID

- **Endpoint:** `GET /task/:id`
- **Response:** Task details as above.

---

### Update Task

- **Endpoint:** `PUT /task/update`
- **Request:** Same as create, with `id` field.

---

### Delete Task

- **Endpoint:** `DELETE /task/delete/:id`
- **Response:**
    ```json
    { "message": "Task deleted successfully" }
    ```

---

### List All Tasks

- **Endpoint:** `GET /task/list`
- **Response:**
    ```json
    {
      "tasks": [
        {
            "id": "484c0428-6609-472f-8be1-5f3dd9b204c0",
            "name": "test-task-1",
            "description": "test task by ankesh kmr",
            "status": "STARTED",
            "priority": 1,
            "created_by": "00000000-0000-0000-0000-000000000000",
            "task_data": "{\"a\":123}",
            "created_at": "2025-07-31T04:39:01.538787+05:30",
            "updated_at": "2025-07-31T04:39:01.538787+05:30",
            "requested_at": "0001-01-01T05:53:28+05:53"
        },
        {
            "id": "f43d75ef-83b2-493b-a652-5f4ad2c74349",
            "name": "test-task-2",
            "description": "test task by ankesh kmr",
            "status": "STARTED",
            "priority": 1,
            "created_by": "00000000-0000-0000-0000-000000000000",
            "task_data": "{\"a\":123}",
            "created_at": "2025-07-31T04:39:54.096812+05:30",
            "updated_at": "2025-07-31T04:39:54.096812+05:30",
            "requested_at": "0001-01-01T05:53:28+05:53"
        }
      ]
    }
    ```

---

### List Tasks (Paginated)

- **Endpoint:** `POST /task/list_paginated?status={status}`
- **Request:**
    ```json
     {
      "cursor": {
        "created_at": "2025-07-31T04:40:28.584057+05:30",
        "serial_number": 4
      }
    }
    ```
- **Response:**
    ```json
    {
      "tasks": [ ... ],
      "cursor": { "created_at": "...", "serial_number": 11 }
    }
    ```

---
## Microservices Concepts Demonstrated

- **Separation of Concerns:** Controllers, repositories, models, and drivers are decoupled.
- **Scalability:** Service is stateless and horizontally scalable.
- **Extensibility:** New microservices (e.g., User Service) can be added; communication via REST.
- **API Design:** Consistent RESTful endpoints.

---

## Future Enhancements:
- **Add Tests:** Implement unit tests and funcitonal tests. Leverage test-containers.
- **Authentication & Authorization:** Implement JWT-based security.
- **AccessControl:** Allow access control for who gets to update or delete tasks.
- **Maintain Task Versions:** Track changes to tasks over time.
- **Task Dependencies:** Allow tasks to depend on other tasks.
- **Advanced Filtering:** Add more complex query capabilities.
- **Task Notifications:** Implement a notification system for task updates.