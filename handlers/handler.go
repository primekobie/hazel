package handlers

import (
	"github.com/freekobie/hazel/services"
	"github.com/go-playground/validator/v10"
)

// TODO: remove global variable
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Handler struct {
	users      *services.UserService
	workspaces *services.WorkspaceService
}

func NewHandler(us *services.UserService, wks *services.WorkspaceService) *Handler {
	return &Handler{
		users:      us,
		workspaces: wks,
	}
}
