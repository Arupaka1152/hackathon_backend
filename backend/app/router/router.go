package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func Init() {
	e := echo.New()

	e.POST("/signup", controller.Signup)
	e.POST("/login", controller.Login)

	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}))

	api := e.Group("/api")

	api.GET("/workspace/all", controller.FetchAllWorkSpace)
	api.POST("/workspace/create", controller.CreateWorkspace)
	api.DELETE("/workspace/delete", controller.DeleteWorkspace)
	api.POST("/workspace/invite", controller.CreateUser)
	api.DELETE("/workspace/remove", controller.DeleteUserFromWorkspace)
	api.GET("/workspace/member", controller.FetchAllUsersInWorkspace)
	api.POST("/workspace/role", contoller.GrantRoleToUser)

	api.GET("/contribution/all", controller.FetchAllContributionInWorkspace)
	api.GET("/contribution/sent", controller.FetchAllContributionSent)
	api.GET("/contribution/received", controller.FetchAllContributionReceived)
	api.POST("/contribution", controller.CreateContribution)
	api.POST("/contribution/reaction", controller.SendReaction)
	api.PUT("/contribution", controller.EditContribution)
	api.DELETE("/contribution", contoller.DeleteContribution)

	api.DELETE("/user", controller.DeleteUserFromWorkspace)
	api.POST("/user/edit/name", controller.ChangeUserName)
	api.POST("/user/edit/avatar", controller.ChangeUserAvatarUrl)

	e.Logger.Fatal(e.Start(":8080"))
}
