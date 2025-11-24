package main

import (
	"github.com/primekobie/hazel/docs"
	"github.com/primekobie/hazel/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func (app *application) routes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/api/v1"

	open := router.Group("/api/v1")
	open.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "200",
			"message": "online",
		})
	})

	// users
	open.POST("/auth/register", app.handler.CreateUser)
	open.POST("/auth/login", app.handler.LoginUser)
	open.POST("/auth/access", app.handler.GetUserAccessToken)
	open.POST("/auth/verify", app.handler.VerifyUser)
	open.POST("/auth/verify/request", app.handler.RequestVerification)

	protected := open.Group("/")
	protected.Use(middlewares.Authentication())
	{
		//users
		protected.GET("/users/:id", app.handler.GetUser)
		protected.PATCH("/users/profile", app.handler.UpdateUserData)
		protected.DELETE("/users/:id", app.handler.DeleteUser)

		// workspaces
		protected.POST("/workspaces", app.handler.CreateWorkspace)
		protected.GET("/workspaces/:id", app.handler.GetWorkspace)
		protected.GET("/workspaces/me", app.handler.GetUserWorkspaces)
		protected.PATCH("/workspaces/:id", app.handler.UpdateWorkspace)
		protected.DELETE("/workspaces/:id", app.handler.DeleteWorkspace)
		protected.POST("/workspaces/:id/members", app.handler.AddWorkspaceMember)
		protected.GET("/workspaces/:id/members", app.handler.GetWorkspaceMembers)
		protected.DELETE("/workspaces/:id/members/:user_id", app.handler.DeleteWorkspaceMember)
		protected.GET("/workspaces/:id/projects", app.handler.GetProjectsInWorkspace)

		// projects
		protected.POST("/projects", app.handler.CreateProject)
		protected.GET("/projects/:id", app.handler.GetProject)
		protected.PATCH("/projects/:id", app.handler.UpdateProject)
		protected.DELETE("/projects/:id", app.handler.DeleteProject)
		protected.GET("/projects/:id/tasks", app.handler.GetProjectTasks)

		// Tasks
		protected.POST("/tasks", app.handler.CreateTask)
		protected.GET("/tasks/:id", app.handler.GetTask)
		protected.PATCH("/tasks/:id", app.handler.UpdateTask)
		protected.DELETE("/tasks/:id", app.handler.DeleteTask)
		protected.POST("/tasks/:id/assignments", app.handler.AssignTaskToUser)
		protected.GET("/tasks/:id/assignments", app.handler.GetAssignedUsers)
		protected.DELETE("/tasks/:id/assignments/:user_id", app.handler.RemoveAssignment)
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
