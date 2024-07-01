package main

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/gofiber/fiber/v2"
	"github.com/jpiontek/go-ip-api"
	"os"
)

var apiKey = os.Getenv("OWM_API_KEY")

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		client := goip.NewClient()
		location, err := client.GetLocation()
		if err != nil {
			return c.Status(500).SendString("Error getting location: " + err.Error())
		}

		if location.City == "" {
			return c.Status(500).SendString("City not found")
		}

		fmt.Println("City:", location.City)
		fmt.Println("API Key:", apiKey)

		weather, err := owm.NewCurrent("C", "EN", apiKey)
		if err != nil {
			return c.Status(500).SendString("Error initializing weather: " + err.Error())
		}

		weatherData := weather.CurrentByName(location.City)

		fmt.Println("Weather:", weatherData)
		return c.SendString(fmt.Sprintf("Weather in %s: %v", location.City, weatherData))
	})

	app.Listen(":5000")
}
