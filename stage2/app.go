package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mryan-3/hng11/stage2/database"
	"github.com/mryan-3/hng11/stage2/initializer"
	"github.com/mryan-3/hng11/stage2/routes"
)



func init() {
	fmt.Println("Initializing the server ...")
	database.ConnectDb()
}

func main (){

	fmt.Println("Starting the server ...")
    app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CLIENT_FRONTEND_URL") + "," + os.Getenv("ADMIN_FRONTEND_URL"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))


    routes.SetUpRoutes(app)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Server is online.")
	})

	// Handle invalid routes
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})


    port := os.Getenv("PORT")
	app.Listen(":"+port)
}
