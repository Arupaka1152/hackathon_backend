package router

import (
	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()

	e.POST("/signup", controller.Signup)
	e.POST("/login", controller.Login)

}
