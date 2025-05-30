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
	us *services.UserService
}

func NewHandler(us *services.UserService) *Handler {
	return &Handler{
		us: us,
	}
}
