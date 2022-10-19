package router

import (
	"backend/app/controller"
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

	r := e.Group("/restricted")

	r.GET("/workspace/all", controller.FetchAllWorkSpace)
	r.POST("/workspace/create", controller.CreateWorkspace)
	r.DELETE("/workspace/delete", controller.DeleteWorkspace)
	r.POST("/workspace/invite", controller.CreateUser)
	r.DELETE("/workspace/remove", controller.DeleteUserFromWorkspace)
	r.GET("/workspace/member", controller.FetchAllUsersInWorkspace)
	r.POST("/workspace/role", controller.GrantRoleToUser)
	r.PUT("/workspace/edit", controller.ChangeWorkspaceAttributes)

	r.GET("/contribution", controller.FetchAllContributionInWorkspace)
	r.GET("/contribution/sent", controller.FetchAllContributionSent)
	r.GET("/contribution/received", controller.FetchAllContributionReceived)
	r.POST("/contribution", controller.CreateContribution)
	r.POST("/contribution/reaction", controller.SendReaction)
	r.PUT("/contribution", controller.EditContribution)
	r.DELETE("/contribution", controller.DeleteContribution)

	r.DELETE("/user", controller.DeleteUserFromWorkspace)
	r.POST("/user", controller.ChangeUserAttributes)

	e.Logger.Fatal(e.Start(":8080"))
}
