package router

import (
	"backend/app/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() {
	g := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
		},
	}
	g.Use(cors.New(corsConfig))

	g.POST("/signup", controller.Signup)
	g.POST("/login", controller.Login)

	api := g.Group("/api")

	api.GET("/workspace", controller.FetchAllWorkSpaces)
	api.POST("/workspace", controller.CreateWorkspace)
	api.DELETE("/workspace", controller.DeleteWorkspace)
	api.POST("/workspace/invite", controller.CreateUser)
	api.GET("/workspace/member", controller.FetchAllUsersInWorkspace)
	api.POST("/workspace/role", controller.GrantRoleToUser)
	api.PUT("/workspace", controller.ChangeWorkspaceAttributes)

	api.GET("/contribution", controller.FetchAllContributionInWorkspace)
	api.GET("/contribution/sent", controller.FetchAllContributionSent)
	api.POST("/contribution/received", controller.FetchAllContributionReceived)
	api.POST("/contribution", controller.CreateContribution)
	api.POST("/contribution/reaction", controller.SendReaction)
	api.PUT("/contribution", controller.EditContribution)
	api.DELETE("/contribution", controller.DeleteContribution)

	api.DELETE("/user", controller.DeleteUser)
	api.POST("/user", controller.ChangeUserAttributes)

	g.Run(":8000")
}
