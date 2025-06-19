package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/freekobie/hazel/models"
	"github.com/freekobie/hazel/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateWorkspace(c *gin.Context) {
	var input struct {
		Name        string    `json:"name" binding:"required"`
		Description string    `json:"description"`
		UserID      uuid.UUID `json:"userId" binding:"required,uuid"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ws := &models.Workspace{
		Name:        input.Name,
		Description: input.Description,
		User:        &models.User{Id: input.UserID},
	}

	err = h.wss.NewWorkspace(c.Request.Context(), ws)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ws)
}

func (h *Handler) GetWorkspace(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	ws, err := h.wss.GetWorkspace(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, ws)
}

// Get all the workspaces where a user has membership.
func (h *Handler) GetUserWorkspaces(c *gin.Context) {
	idStr, _ := c.Get("user_id")

	workspaces, err := h.wss.GetUserWorkspaces(c.Request.Context(), uuid.MustParse(idStr.(string)))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

func (h *Handler) UpdateWorkspace(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	var input map[string]string

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input["id"] = id.String()

	ws, err := h.wss.UpdateWorkspace(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, ws)
}

func (h *Handler) DeleteWorkspace(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	err = h.wss.DeleteWorkspace(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workspace successfully deleted"})
}

func (h *Handler) AddWorkspaceMember(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	var input struct {
		UserId uuid.UUID `json:"userId" validate:"required,uuid"`
		Role   string    `json:"role" validate:"required"`
	}

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.wss.AddWorkspaceMember(c.Request.Context(), id, input.UserId, input.Role)
	if err != nil {
		if errors.Is(err, services.ErrFailedOperation) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
			return
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workspace member successfully added"})
}

func (h *Handler) GetWorkspaceMembers(c *gin.Context) {

	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	users, err := h.wss.GetWorkspaceMembers(c.Request.Context(), id)
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

func (h *Handler) DeleteWorkspaceMember(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	memberId, err := getUUIDparam(c, "member_id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	err = h.wss.DeleteWorkspaceMember(c.Request.Context(), id, memberId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workspace member successfully deleted"})
}
