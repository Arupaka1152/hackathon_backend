package router

import (
	"backend/app/controller"
	"backend/app/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() {
	g := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins: []string{
			"*",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"authentication",
			"workspace_id",
		},
	}
	g.Use(cors.New(corsConfig))

	g.POST("/signup", controller.Signup)
	g.POST("/login", controller.Login)

	g.POST("/workspace", controller.CreateWorkspace)
	g.GET("/workspace", controller.FetchAllWorkSpaces)

	api := g.Group("/api")

	api.Use(cors.New(corsConfig))
	api.Use(middleware.AuthMiddleware())

	api.GET("/workspace/member", controller.FetchAllUsersInWorkspace)
	api.POST("/workspace/invite", controller.CreateUser)
	api.POST("/workspace/role", controller.GrantRoleToUser)
	api.POST("/workspace/remove", controller.RemoveUserFromWorkspace)
	api.PUT("/workspace", controller.ChangeWorkspaceAttributes)
	api.DELETE("/workspace", controller.DeleteWorkspace)

	api.GET("/contribution", controller.FetchAllContributionInWorkspace)
	api.GET("/contribution/sent", controller.FetchAllContributionSent)
	api.GET("/contribution/received", controller.FetchAllContributionReceived)
	api.POST("/contribution", controller.CreateContribution)
	api.POST("/contribution/reaction", controller.SendReaction)
	api.PUT("/contribution", controller.EditContribution)
	api.DELETE("/contribution", controller.DeleteContribution)

	api.GET("/me", controller.FetchUserInfo)
	api.PUT("/user", controller.ChangeUserAttributes)
	api.DELETE("/user", controller.DeleteUser)

	g.Run(":8080")
}
