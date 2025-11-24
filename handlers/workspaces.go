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

// CreateWorkspace godoc
//	@Summary		Create workspace
//	@Description	Create a new workspace
//	@Tags			workspaces
//	@Accept			json
//	@Produce		json
//	@Param			workspace	body		object	true	"Workspace info"
//	@Success		201			{object}	models.Workspace
//	@Failure		400			{object}	map[string]string
//	@Router			/workspaces [post]
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

	err = h.workspaces.NewWorkspace(c.Request.Context(), ws)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ws)
}

// GetWorkspace godoc
//	@Summary		Get workspace
//	@Description	Get a workspace by ID
//	@Tags			workspaces
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Workspace ID"
//	@Success		200	{object}	models.Workspace
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/workspaces/{id} [get]
func (h *Handler) GetWorkspace(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	ws, err := h.workspaces.GetWorkspace(c.Request.Context(), id)
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

// GetUserWorkspaces godoc
//	@Summary		Get my workspaces
//	@Description	Get all workspaces where the authenticated user has membership
//	@Tags			workspaces
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.Workspace
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/workspaces/me [get]
func (h *Handler) GetUserWorkspaces(c *gin.Context) {
	idStr, _ := c.Get("user_id")

	workspaces, err := h.workspaces.GetUserWorkspaces(c.Request.Context(), uuid.MustParse(idStr.(string)))
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

// UpdateWorkspace godoc
//	@Summary		Update workspace
//	@Description	Update workspace details
//	@Tags			workspaces
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"Workspace ID"
//	@Param			workspace	body		object	true	"Workspace update info"
//	@Success		200			{object}	models.Workspace
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/workspaces/{id} [patch]
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

	ws, err := h.workspaces.UpdateWorkspace(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, ws)
}

// DeleteWorkspace godoc
//	@Summary		Delete workspace
//	@Description	Delete a workspace by ID
//	@Tags			workspaces
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Workspace ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/workspaces/{id} [delete]
func (h *Handler) DeleteWorkspace(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	err = h.workspaces.DeleteWorkspace(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workspace successfully deleted"})
}

// AddWorkspaceMember godoc
//	@Summary		Add workspace member
//	@Description	Add a member to a workspace
//	@Tags			workspaces
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Workspace ID"
//	@Param			member	body		object	true	"Member info"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/workspaces/{id}/members [post]
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

	err = h.workspaces.AddWorkspaceMember(c.Request.Context(), id, input.UserId, input.Role)
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

// GetWorkspaceMembers godoc
//	@Summary		Get workspace members
//	@Description	Get all members of a workspace
//	@Tags			workspaces
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Workspace ID"
//	@Success		200	{array}		models.User
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/workspaces/{id}/members [get]
func (h *Handler) GetWorkspaceMembers(c *gin.Context) {

	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	users, err := h.workspaces.GetWorkspaceMembers(c.Request.Context(), id)
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

// DeleteWorkspaceMember godoc
//	@Summary		Remove workspace member
//	@Description	Remove a member from a workspace
//	@Tags			workspaces
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id		path		string	true	"Workspace ID"
//	@Param			user_id	path		string	true	"Member ID"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/workspaces/{id}/members/{user_id} [delete]
func (h *Handler) DeleteWorkspaceMember(c *gin.Context) {
	id, err := getUUIDparam(c, "id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	memberId, err := getUUIDparam(c, "user_id")
	if err != nil {
		slog.Error("failed to get uuid param", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	err = h.workspaces.DeleteWorkspaceMember(c.Request.Context(), id, memberId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workspace member successfully deleted"})
}
