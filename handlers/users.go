package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/primekobie/hazel/models"
	"github.com/primekobie/hazel/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateUser godoc
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		object	true	"User registration info"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/register [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8,max=20"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.users.CreateUser(c.Request.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUser) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// VerifyUser godoc
//	@Summary		Verify user email
//	@Description	Verify a user's email with a code
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			verification	body		object	true	"Verification info"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/auth/verify [post]
func (h *Handler) VerifyUser(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.users.VerifyUser(c.Request.Context(), input.Code, input.Email)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// RequestVerification godoc
//	@Summary		Request verification email
//	@Description	Request a new verification code for a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			email	body		object	true	"User email"
//	@Success		202		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/verify/request [post]
func (h *Handler) RequestVerification(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := h.users.ResendVerificationEmail(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		} else if strings.Contains(err.Error(), "user already verified") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("new verification code has ben sent to '%s'", input.Email)})

}

// LoginUser godoc
//	@Summary		Login user
//	@Description	Authenticate user and return session tokens
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		object	true	"User credentials"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/auth/login [post]
func (h *Handler) LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	session, err := h.users.NewSession(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, services.ErrFailedOperation) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetUserAccessToken godoc
//	@Summary		Refresh access token
//	@Description	Get a new access token using a refresh token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			refreshToken	body		object	true	"Refresh token"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		400				{object}	map[string]string
//	@Failure		401				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/auth/access [post]
func (h *Handler) GetUserAccessToken(c *gin.Context) {

	var input struct {
		RefresToken string `json:"refreshToken" binding:"required,jwt"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	access, err := h.users.RefreshSession(c.Request.Context(), input.RefresToken)
	if err != nil {
		if errors.Is(err, services.ErrFailedOperation) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, access)

}

// GetUser godoc
//	@Summary		Get user by ID
//	@Description	Get user details by user ID
//	@Tags			users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	userId := c.Param("id")
	if err := validate.Var(userId, "uuid"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	user, err := h.users.FetchUser(c.Request.Context(), uuid.MustParse(userId))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserData godoc
//	@Summary		Update user data
//	@Description	Update user profile information
//	@Tags			users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		object	true	"User update info"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/users/profile [patch]
func (h *Handler) UpdateUserData(c *gin.Context) {
	var input map[string]any
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	idString, ok := c.Get("user_id")
	if !ok {
		slog.Error("failed to fetch user id from context")
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	id := uuid.MustParse(idString.(string))
	input["id"] = id

	user, err := h.users.UpdateUser(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrFailedOperation) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
//	@Summary		Delete user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	userId := c.Param("id")
	if err := validate.Var(userId, "uuid"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	err := h.users.DeleteUser(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
