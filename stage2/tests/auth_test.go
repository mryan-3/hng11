package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mryan-3/hng11/stage2/controller"
	"github.com/mryan-3/hng11/stage2/database"
	"github.com/mryan-3/hng11/stage2/models"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	godotenv.Load("../.env.test")

	database.ConnectTestDb() // Make sure this connects to a test database
	app.Post("/auth/register", controller.CreateUser)
	app.Post("/auth/login", controller.LoginUser)
	return app
}

func TestRegisterEndpoint(t *testing.T) {
	app := setupTestApp()

	t.Run("Should Register User Successfully with Default Organisation", func(t *testing.T) {
		reqBody := map[string]string{
			"firstName": "Jill",
			"lastName":  "Doe",
			"email":     "rill@example.com",
			"password":  "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		// Read and log the response body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("Response Body: %s", string(body))

		// Attempt to parse the JSON response
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)

		fmt.Printf("Decoded Response Body: %+v\n", result)

		if err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		data := result["data"].(map[string]interface{})
		userData := data["user"].(map[string]interface{})

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "Jill", userData["firstName"])
		assert.Equal(t, "Doe", userData["lastName"])
		assert.Equal(t, "rill@example.com", userData["email"])
		assert.NotNil(t, data["accessToken"])

		var org models.Organisation
		database.DB.Db.Preload("Users").Where("name = ?", "Jill's Organisation").First(&org)

		assert.Equal(t, "Jill's Organisation", org.Name)
		assert.Equal(t, 1, len(org.Users))
		assert.Equal(t, "rill@example.com", org.Users[0].Email)
	})

	t.Run("Should Log the user in successfully", func(t *testing.T) {
		// First, register a user
		registerReqBody := map[string]string{
			"firstName": "Jan",
			"lastName":  "Doe",
			"email":     "jin@example.com",
			"password":  "password123",
		}
		jsonRegisterBody, _ := json.Marshal(registerReqBody)
		registerReq := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonRegisterBody))
		registerReq.Header.Set("Content-Type", "application/json")

		app.Test(registerReq)

		// Now, try to log in
		loginReqBody := map[string]string{
			"email":    "jan@example.com",
			"password": "password123",
		}
		jsonLoginBody, _ := json.Marshal(loginReqBody)
		loginReq := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonLoginBody))
		loginReq.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(loginReq)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}

		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "success", result["status"])

		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok, "Data field is not a map[string]interface{}")

		user, ok := data["user"].(map[string]interface{})
		assert.True(t, ok, "User field is not a map[string]interface{}")

		assert.NotEmpty(t, user["userId"])
		assert.NotEmpty(t, user["firstName"])
		assert.NotEmpty(t, user["lastName"])
		assert.NotEmpty(t, data["accessToken"])

	})

	t.Run("It Should Fail If Required Fields Are Missing", func(t *testing.T) {

		testCases := []struct {
			name       string
			user       map[string]string
			missingKey string
		}{
			{
				name: "Missing FirstName",
				user: map[string]string{
					"lastName": "Doe", "email": "john@example.com", "password": "password123", "phone": "1234567890",
				},
				missingKey: "firstName",
			},
			{
				name: "Missing LastName",
				user: map[string]string{
					"firstName": "John", "email": "john@example.com", "password": "password123", "phone": "1234567890",
				},
				missingKey: "lastName",
			},
			{
				name: "Missing Email",
				user: map[string]string{
					"firstName": "John", "lastName": "Doe", "password": "password123", "phone": "1234567890",
				},
				missingKey: "email",
			},
			{
				name: "Missing Password",
				user: map[string]string{
					"firstName": "John", "lastName": "Doe", "email": "john@example.com", "phone": "1234567890",
				},
				missingKey: "password",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.user)

				req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)
				if err != nil {
					t.Fatalf("Failed to perform request: %v", err)
				}

				assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read response body: %v", err)
				}
				resp.Body.Close()
				fmt.Printf("Raw Response Body (%s): %s\n", tc.name, string(respBody))

				if len(respBody) == 0 {
					t.Fatalf("Response body is empty")
				}

				var responseBody map[string]interface{}
				err = json.Unmarshal(respBody, &responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				fmt.Printf("Decoded Response Body (%s): %+v\n", tc.name, responseBody)

				errors, ok := responseBody["errors"].([]interface{})
				if !ok {
					t.Fatalf("Response body 'errors' field is not of expected type: %v", responseBody)
				}

				for _, err := range errors {
					errorMap := err.(map[string]interface{})
					field, fieldExists := errorMap["field"].(string)
					message, messageExists := errorMap["message"].(string)
					assert.True(t, fieldExists)
					assert.True(t, messageExists)
					if fieldExists && messageExists {
						switch field {
						case tc.missingKey:
							assert.Equal(t, fmt.Sprintf("%s is required", field), message)
						}
					}
				}
			})
		}
	})

	t.Run("Should Fail if there's Duplicate Email", func(t *testing.T) {
		reqBody := map[string]string{
			"firstName": "John",
			"lastName":  "Doe",
			"email":     "duplicate@example.com",
			"password":  "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		// Register the first user
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req)

		// Try to register another user with the same email
		req = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "Bad Request", result["status"])
		assert.Equal(t, "Registration unsuccessful", result["message"])

	})

}

