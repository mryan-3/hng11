package routes

import (
	"github.com/gofiber/fiber/v2"
	userControllers "github.com/mryan-3/hng11/stage2/controller"
	organisationControllers "github.com/mryan-3/hng11/stage2/controller"
	"github.com/mryan-3/hng11/stage2/middleware"
)


func SetUpRoutes(app *fiber.App) {
    api := app.Group("/api")
    user := api.Group("/users")

    app.Post("/auth/register", userControllers.CreateUser)
    app.Post("/auth/login", userControllers.LoginUser)

    // User organisation routes
    api.Get("/organisations", middleware.UserAuth, organisationControllers.GetUserOrganisations)
    api.Get("/users", userControllers.GetUsers)
    api.Get("/organisations/:orgId", middleware.UserAuth, organisationControllers.GetSingleOrganisation)
    api.Post("/organisations", middleware.UserAuth, organisationControllers.CreateOrganisation)
    api.Post("/organisations/:orgId/users", middleware.UserAuth, organisationControllers.AddUserToOrganisation)

	// User routes
    user.Get("/:id", middleware.UserAuth, userControllers.GetUser)

}
