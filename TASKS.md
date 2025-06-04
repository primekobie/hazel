# TODO: Project Management Backend (Go + Gin)

## ✅ Initial Setup
- [X] Initialize Git repository
- [X] Create project structure 
- [X] Initialize Go module (`go mod init`)
- [X] Set up configuration management (`viper` or env parsing)
- [X] Create `.env` for local development
- [X] Set up Dockerfile and docker-compose.yml
- [ ] Integrate PostgreSQL via Docker
- [X] Create Taskfile.yaml scripts for build and dev tasks

---

## 🔐 Authentication & Authorization (JWT)
- [X] User registration
- [ ] User login
- [X] Password hashing (bcrypt)
- [ ] JWT generation (access token)
- [ ] JWT middleware for route protection
- [ ] **Role-based middleware (Admin, Manager, User)** [Maybe]

---

## 🗂️ Core Models & Database Schema
- [X] User
- [ ] Workspace
- [ ] Project (belongs to Workspace)
- [ ] Task (belongs to Project)
- [ ] Comment (optional, belongs to Task)
- [ ] Relationships and constraints
- [ ] Auto migrations (`github.com/pressly/goose`)
- [ ] Database seeding for development

---

## 🔁 API Endpoints

### Auth
- [X] `POST /auth/register` – Register new user
- [X] `POST /auth/verify` – Verify user email address
- [ ] `POST /auth/login` – Authenticate and return JWT

### Users
- [ ] `GET /users/me` – Get current user profile
- [ ] `PATCH /users/me` – Update user profile
- [ ] `DELETE /users/me` – Delete account

### Workspaces
- [ ] `POST /workspaces` – Create workspace
- [ ] `GET /workspaces` – List all user workspaces
- [ ] `GET /workspaces/:id` – Get specific workspace
- [ ] `PATCH /workspaces/:id` – Update workspace
- [ ] `DELETE /workspaces/:id` – Delete workspace

### Projects
- [ ] `POST /workspaces/:workspaceId/projects` – Create project
- [ ] `GET /workspaces/:workspaceId/projects` – List projects in a workspace
- [ ] `GET /projects/:id` – Get project details
- [ ] `PATCH /projects/:id` – Update project
- [ ] `DELETE /projects/:id` – Delete project

### Tasks
- [ ] `POST /projects/:projectId/tasks` – Create task
- [ ] `GET /projects/:projectId/tasks` – List tasks in a project
- [ ] `GET /tasks/:id` – Get task
- [ ] `PATCH /tasks/:id` – Update task
- [ ] `DELETE /tasks/:id` – Delete task
- [ ] `PATCH /tasks/:id/assign` – Assign task to user
- [ ] `PATCH /tasks/:id/status` – Update task status (To Do, In Progress, Done)

### Comments (optional)
- [ ] `POST /tasks/:taskId/comments` – Add comment
- [ ] `GET /tasks/:taskId/comments` – Get all comments on task

---

## ⚙️ Middleware & Utilities
- [ ] JWT authentication middleware
- [ ] Error handling middleware
- [ ] Request logging (e.g., Gin's logger or `zap`)
- [ ] Rate limiting (e.g., `golang.org/x/time/rate`)
- [ ] Input validation (e.g., `go-playground/validator`)

---

## 🧪 Testing
- [ ] Unit tests for handlers
- [ ] Unit tests for services and database logic
- [ ] Integration tests with PostgreSQL
- [ ] Authentication & RBAC tests
- [ ] Mock external dependencies
- [ ] Test coverage tracking

---

## 📊 Documentation
- [ ] API documentation (Swagger / OpenAPI via `swaggo`)
- [ ] Postman collection or Insomnia setup
- [ ] README.md with setup, build, and API usage instructions

---

## 🚀 Deployment
- [ ] Production-ready Dockerfile
- [ ] `docker-compose.yml` for dev and prod
- [ ] Environment-based config (dev, staging, prod)
- [ ] Health check endpoint
- [ ] Deploy to platform (e.g., Render, AWS, Railway)
- [ ] Configure persistent PostgreSQL volume
- [ ] Setup CORS, HTTPS, etc.

---
