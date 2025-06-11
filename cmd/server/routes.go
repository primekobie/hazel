package main

import (
	"github.com/freekobie/hazel/middlewares"
	"github.com/gin-gonic/gin"
)

func (s *application) routes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// users
	router.POST("/auth/register", s.h.CreateUser)
	router.POST("/auth/login", s.h.LoginUser)
	router.POST("/auth/access", s.h.GetUserAccessToken)
	router.POST("/auth/verify", s.h.VerifyUser)
	router.POST("/auth/verify/request", s.h.RequestVerification)

	authorized := router.Group("/")
	authorized.Use(middlewares.Authentication())
	{
		//users
		authorized.GET("/users/:id", s.h.GetUser)
		authorized.PATCH("/users/profile", s.h.UpdateUserData)
		authorized.DELETE("/users/:id", s.h.DeleteUser)

		// workspaces
		authorized.POST("/workspaces", s.h.CreateWorkspace)
		authorized.GET("/workspaces/:id", s.h.GetWorkspace)
		authorized.GET("/workspaces/me", s.h.GetUserWorkspaces)
		authorized.PATCH("/workspaces/:id", s.h.UpdateWorkspace)
		authorized.GET("/workspaces/:id/members")

		// projects

		// Tasks

		// Comments
	}

	return router
}
