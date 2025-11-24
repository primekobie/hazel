package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/primekobie/hazel/models"
	"github.com/primekobie/hazel/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTask godoc
//	@Summary		Create task
//	@Description	Create a new task in a project
//	@Security		BearerAuth
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		object	true	"Task info"
//	@Success		201		{object}	models.Task
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/tasks [post]
func (h *Handler) CreateTask(c *gin.Context) {
	var input struct {
		ProjectId   uuid.UUID           `json:"projectId" binding:"required,uuid"`
		Title       string              `json:"title" binding:"required"`
		Description string              `json:"description"`
		Due         time.Time           `json:"due"`
		Priority    models.TaskPriority `json:"priority"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	task := &models.Task{
		Title:       input.Title,
		Description: input.Description,
		Project:     &models.Project{Id: input.ProjectId},
		Due:         input.Due,
		Priority:    input.Priority,
	}
	err = h.workspaces.CreateTask(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTask godoc
//	@Summary		Get task
//	@Description	Get a task by ID
//	@Security		BearerAuth
//	@Tags			tasks
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	models.Task
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/tasks/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	task, err := h.workspaces.GetTask(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, task)

}

// UpdateTask godoc
//	@Summary		Update task
//	@Description	Update task details
//	@Security		BearerAuth
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Task ID"
//	@Param			task	body		object	true	"Task update info"
//	@Success		200		{object}	models.Task
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/tasks/{id} [patch]
func (h *Handler) UpdateTask(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	var input map[string]any

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input["id"] = id

	task, err := h.workspaces.UpdateTask(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		} else if errors.Is(err, services.ErrInvalidDateFormat) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetProjectTasks godoc
//	@Summary		Get project tasks
//	@Description	Get all tasks for a project
//	@Security		BearerAuth
//	@Tags			tasks
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{array}		models.Task
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/projects/{id}/tasks [get]
func (h *Handler) GetProjectTasks(c *gin.Context) {

	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	tasks, err := h.workspaces.GetProjectTasks(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// DeleteTask godoc
//	@Summary		Delete task
//	@Description	Delete a task by ID
//	@Tags			tasks
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	err = h.workspaces.DeleteTask(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project successfully deleted"})
}

// AssignTaskToUser godoc
//	@Summary		Assign task to user
//	@Description	Assign a task to a user
//	@Tags			tasks
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"Task ID"
//	@Param			assignment	body		object	true	"Assignment info"
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	map[string]string
//	@Failure		422			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/tasks/{id}/assignments [post]
func (h *Handler) AssignTaskToUser(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	var input struct {
		UserId uuid.UUID `json:"userId" validate:"required,uuid"`
	}

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.workspaces.AssignTaskToUser(c.Request.Context(), id, input.UserId)
	if err != nil {
		if errors.Is(err, services.ErrFailedOperation) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
			return
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task successfully assigned to user"})

}

// RemoveAssignment godoc
//	@Summary		Remove task assignment
//	@Description	Remove a user's assignment from a task
//	@Tags			tasks
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id		path		string	true	"Task ID"
//	@Param			user_id	path		string	true	"User ID"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/tasks/{id}/assignments/{user_id} [delete]
func (h *Handler) RemoveAssignment(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	userId, err := getUUIDparam(c, "user_id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	err = h.workspaces.UnassignTask(c.Request.Context(), id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task assignment successfully removed"})
}

// GetAssignedUsers godoc
//	@Summary		Get assigned users
//	@Description	Get all users assigned to a task
//	@Tags			tasks
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{array}		models.User
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/tasks/{id}/assignments [get]
func (h *Handler) GetAssignedUsers(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	users, err := h.workspaces.GetAssignedUsers(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
