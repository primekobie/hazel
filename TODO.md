# TODO: Project Management Backend (Go + Gin)

## Initial Setup
- [X] Initialize Git repository
- [X] Create project structure 
- [X] Initialize Go module (`go mod init`)
- [X] Set up Dockerfile and docker-compose.yml
- [X] Create Taskfile.yaml scripts for build and dev tasks

---

## Authentication & Authorization (JWT)
- [X] User registration
- [X] User login
- [X] Password hashing (bcrypt)
- [X] JWT generation (access token)
- [X] JWT middleware for route protection
- [ ] **Role-based middleware (Admin, Manager, User)** 

---

## Core Models & Database Schema
- [X] User
- [X] Workspace
- [X] Project (belongs to Workspace)
- [X] Task (belongs to Project)
- [X] Auto migrations (`github.com/pressly/goose`)

---

## API Endpoints

### Auth
- [X] `POST /auth/register` – Register new user
- [X] `POST /auth/verify` – Verify user email address
- [X] `POST /auth/login` – Authenticate and return JWT
- [X] `POST /auth/verify/request` - Request email verification code

### Users
- [X] `GET /users/:id` – Get user profile
- [X] `PATCH /users/me` – Update user profile
- [X] `DELETE /users/me` – Delete account

### Workspaces
- [X] `POST /workspaces` – Create workspace
- [X] `GET /workspaces` – List all user workspaces
- [X] `GET /workspaces/:id` – Get specific workspace
- [X] `PATCH /workspaces/:id` – Update workspace
- [X] `DELETE /workspaces/:id` – Delete workspace
- [X] `POST /workspaces/:id/members` - Join a workspace
- [X] `GET /workspaces/:id/members` - Get Workspace members
- [X] `DELETE /workspaces/:id/members` - Leave a workspace
- [X] `GET /workspaces/:id/projects` – List projects in a workspace

### Projects
- [X] `POST /projects` – Create project
- [X] `GET /projects/:id` – Get project details
- [X] `PATCH /projects/:id` – Update project
- [X] `DELETE /projects/:id` – Delete project
- [X] `GET /projects/:id/tasks` – List tasks in a project

### Tasks
- [X] `POST /projects/:projectId/tasks` – Create task
- [X] `GET /projects/:projectId/tasks` – List tasks in a project
- [X] `GET /tasks/:id` – Get task
- [X] `PATCH /tasks/:id` – Update task
- [X] `DELETE /tasks/:id` – Delete task
- [X] `POST /tasks/:id/assign` – Assign task to user
- [X] `DELETE /tasks/:id/assign` – Remove task assignment
- [X] `GET /tasks/:id/assign` – Get all users assigned to a task

---

## Middleware & Utilities
- [X] JWT authentication middleware
- [ ] Error handling middleware
- [ ] Request logging (Gin's logger or `zap`)
- [ ] Rate limiting (`golang.org/x/time/rate`)
- [X] Input validation (`go-playground/validator`)

---

## Testing
- [ ] Unit tests for handlers
- [ ] Unit tests for services and database logic
- [ ] Integration tests with PostgreSQL
- [ ] Authentication & RBAC tests
- [ ] Mock external dependencies
- [ ] Test coverage tracking

---

## Documentation
- [ ] API documentation (Swagger / OpenAPI via `swaggo`)
- [ ] Postman collection or Insomnia setup
- [ ] README.md with setup, build, and API usage instructions

---

## Deployment
- [ ] Production-ready Dockerfile
- [ ] `docker-compose.yml` for dev and prod
- [ ] Environment-based config (dev, staging, prod)
- [ ] Health check endpoint
- [ ] Deploy to platform (e.g., Render, AWS, Railway)
- [ ] Configure persistent PostgreSQL volume
- [ ] Setup CORS, HTTPS, etc.

---
