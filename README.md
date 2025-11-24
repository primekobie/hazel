# Hazel Project Management Backend

Hazel is a backend service for project management app, built with Go, using PostgreSQL as the database.

## Features

- User registration, authentication, and email verification
- Workspaces for organizing projects and users
- Projects and tasks management
- Role-based workspace memberships
- RESTful API endpoints
- JWT-based authentication

> **Check TODO.md to see all features**

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL
- [Goose](https://github.com/pressly/goose) for database migrations
- (Optional) [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- (Optional) [Taskfile](https://taskfile.dev/)

### Setup (Manual)

1. Clone the repository:
   ```sh
   git clone https://github.com/primekobie/hazel.git
   cd hazel
   ```

2. Copy `.env.example` to `.env` and fill in your environment variables.

3. Run database migrations:
   ```sh
   go install github.com/pressly/goose/v3/cmd/goose@latest
   goose -dir ./migrations postgres "$DB_URL" up
   ```

4. Start the server:
   ```sh
   go run ./cmd/server
   ```

### Alternate: Using Docker and Docker Compose

You can run the backend and database using Docker Compose:

```sh
docker compose up --build
```

This will build the Go backend and start both the application and a PostgreSQL database using the provided `docker-compose.yml` and `Dockerfile`.

### Alternate: Using Taskfile

If you have [Task](https://taskfile.dev/) installed, you can use the included `Taskfile.yaml` for common development tasks:

- Run the server:
  ```sh
  task run
  ```
- Run tests:
  ```sh
  task test
  ```
- Run database migrations:
  ```sh
  task up
  ```

## Documentation
> **http://localhost:8080/swagger/index.html**

![Swagger Documentation Screenshot](screenshot.png)

## Running Tests

```sh
go test ./...
```

## Run a development server

You can use [Air](https://github.com/cosmtrek/air) for live-reloading during development:

```sh
go install github.com/cosmtrek/air@latest
air
```

This will automatically reload the server on code changes using the `air` tool.

## Project Structure

- `handlers/` - HTTP route handlers
- `services/` - Business logic
- `models/` - Data models and interfaces
- `postgres/` - PostgreSQL implementations
- `migrations/` - Database schema migrations
