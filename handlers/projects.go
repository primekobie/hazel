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

func (h *Handler) CreateProject(c *gin.Context) {
	var input struct {
		WorkspaceId uuid.UUID   `json:"workspaceId" binding:"required,uuid"`
		Name        string      `json:"name" binding:"required"`
		Description string      `json:"description"`
		StartDate   models.Date `json:"startDate"`
		EndDate     models.Date `json:"endDate"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	project := &models.Project{
		Name:        input.Name,
		Description: input.Description,
		Workspace:   &models.Workspace{Id: input.WorkspaceId},
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	}
	err = h.wss.CreateProject(c.Request.Context(), project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *Handler) GetProject(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	project, err := h.wss.GetProject(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, project)

}

func (h *Handler) UpdateProject(c *gin.Context) {
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

	ws, err := h.wss.UpdateProject(c.Request.Context(), input)
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

	c.JSON(http.StatusOK, ws)
}

func (h *Handler) GetProjectsInWorkspace(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	projects, err := h.wss.GetProjectsForWorkspace(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	err = h.wss.DeleteProject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project successfully deleted"})
}
