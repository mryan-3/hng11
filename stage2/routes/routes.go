package routes

import (
	"github.com/gofiber/fiber/v2"
	userControllers "github.com/mryan-3/hng11/stage2/controller"
)


func SetUpRoutes(app *fiber.App) {
    api := app.Group("/api")
    user := api.Group("/users")

    app.Post("/auth/register", userControllers.CreateUser)
    app.Post("/auth/login", userControllers.LoginUser)

	// User routes
    user.Get("/:id", userControllers.GetUser)

}
