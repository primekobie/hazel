package handlers

import (
	"errors"
	"net/http"

	"github.com/freekobie/hazel/models"
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
	idStr := c.Param("id")

	if err := validate.Var(idStr, "uuid"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	ws, err := h.wss.GetWorkspace(c.Request.Context(), uuid.MustParse(idStr))
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
	idStr := c.Param("id")
	if err := validate.Var(idStr, "uuid"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	var input map[string]string

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input["id"] = idStr

	ws, err := h.wss.UpdateWorkspace(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, ws)
}
