package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/primekobie/hazel/models"
	"github.com/primekobie/hazel/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateProject godoc
//	@Summary		Create project
//	@Description	Create a new project in a workspace
//	@Tags			projects
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			project	body		object	true	"Project info"
//	@Success		201		{object}	models.Project
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/projects [post]
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
	err = h.workspaces.CreateProject(c.Request.Context(), project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProject godoc
//	@Summary		Get project
//	@Description	Get a project by ID
//	@Tags			projects
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{object}	models.Project
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/projects/{id} [get]
func (h *Handler) GetProject(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	project, err := h.workspaces.GetProject(c.Request.Context(), id)
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

// UpdateProject godoc
//	@Summary		Update project
//	@Description	Update project details
//	@Tags			projects
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Project ID"
//	@Param			project	body		object	true	"Project update info"
//	@Success		200		{object}	models.Project
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/projects/{id} [patch]
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

	ws, err := h.workspaces.UpdateProject(c.Request.Context(), input)
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

// GetProjectsInWorkspace godoc
//	@Summary		Get projects in workspace
//	@Description	Get all projects for a workspace
//	@Tags			projects
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Workspace ID"
//	@Success		200	{array}		models.Project
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/workspaces/{id}/projects [get]
func (h *Handler) GetProjectsInWorkspace(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	projects, err := h.workspaces.GetProjectsForWorkspace(c.Request.Context(), id)
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

// DeleteProject godoc
//	@Summary		Delete project
//	@Description	Delete a project by ID
//	@Tags			projects
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/projects/{id} [delete]
func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get id param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	err = h.workspaces.DeleteProject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project successfully deleted"})
}
