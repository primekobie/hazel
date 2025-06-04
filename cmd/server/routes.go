package main

import (
	"github.com/gin-gonic/gin"
)

func (s *application) routes() *gin.Engine {
	router := gin.Default()

	// users
	router.POST("auth/register", s.h.CreateUser)
	router.POST("auth/login", s.h.LoginUser)
	router.POST("/auth/verify", s.h.VerifyUser)
	router.GET("/users/:id", s.h.GetUser)
	router.DELETE("/users/:id", s.h.DeleteUser)

	// users
	router.POST("/workspaces", s.h.CreateWorkspace)
	return router
}
