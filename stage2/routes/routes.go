package routes

import (
	"github.com/gofiber/fiber/v2"
	userControllers "github.com/mryan-3/hng11/stage2/controller"
)


func SetUpRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Version 1
	v1 := api.Group("/v1")

	// User routes
	user := v1.Group("/user")
    user.Post("/auth/register", userControllers.CreateUser)
    user.Post("/auth/login", userControllers.LoginUser)

}
