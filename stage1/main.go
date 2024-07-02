package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
    "os"
    "github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
	"github.com/jpiontek/go-ip-api"
)


type WeatherData struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

func getWeatherData(city string) (WeatherData, error) {
	q := city

    apiKey := os.Getenv("WEATHER_API_KEY")
    if apiKey == "" {
        return WeatherData{}, fmt.Errorf("WEATHER_API_KEY environment variable not set")
    }

	url := "http://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&q=" + q + "&days=1&aqi=no&alerts=no"

	resp, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	var data WeatherData

	err = json.Unmarshal(body, &data)
	if err != nil {
		return WeatherData{}, err
	}

	return data, nil
}

func main() {
	app := fiber.New()
    godotenv.Load()
    port := os.Getenv("PORT")

	app.Get("/api/hello", func(c *fiber.Ctx) error {
		visitorName := c.Query("visitor_name")

		clientIp := c.IP()

		location, err := goip.NewClient().GetLocation()
		if err != nil {
			return c.Status(500).SendString("Error getting location: " + err.Error())
		}

		weather, _ := getWeatherData(location.City)
		temperature := weather.Current.TempC

		greeting := fmt.Sprintf("Hello, %s! The temperature is %.2f°C in %s.", visitorName, temperature, location.City)

		response := fiber.Map{
			"client_ip": clientIp,
			"location":  location.City,
			"greeting":  greeting,
		}

		return c.JSON(response)

	})

    app.Listen(":" + port)
}

