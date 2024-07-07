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
	"github.com/stretchr/testify/assert"
)

func cleanupDatabase(t *testing.T) {
	// Drop all tables
	err := database.DB.Db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error
	if err != nil {
		t.Fatalf("Error cleaning up database: %v", err)
	}

	// Remigrate all tables
	database.MigrateDatabase(database.DB.Db)
}

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
	defer cleanupDatabase(t)

	t.Run("Should Register User Successfully with Default Organisation", func(t *testing.T) {
		reqBody := map[string]string{
			"firstName": "Lia",
			"lastName":  "Doe",
			"email":     "lia@example.com",
			"password":  "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// Print out the response status code
		t.Logf("Response Status Code: %d", resp.StatusCode)

		// Read and log the response body
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Response Body: %s", string(body))

		// Attempt to parse the JSON response
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)

		status, ok := result["status"].(string)
		assert.True(t, ok)
		assert.Equal(t, "success", status)

		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok)

		// Check user details
		user, ok := data["user"].(map[string]interface{})
		assert.True(t, ok)

		assert.Equal(t, reqBody["firstName"], user["firstName"])
		assert.Equal(t, reqBody["lastName"], user["lastName"])
		assert.Equal(t, reqBody["email"], user["email"])
		assert.NotNil(t, data["accessToken"])

		organisations, ok := user["organisations"].([]interface{})
		assert.True(t, ok)

		// Assuming the first organization is the default one
		if len(organisations) > 0 {
			org, ok := organisations[0].(map[string]interface{})
			assert.True(t, ok)

			expectedOrgName := fmt.Sprintf("%s's Organisation", reqBody["firstName"])
			assert.Equal(t, expectedOrgName, org["name"])
		} else {
			t.Fatal("User should have at least one associated organisation")
		}
	})

	t.Run("Should Log the user in successfully", func(t *testing.T) {
		// First, register a user
		registerReqBody := map[string]string{
			"firstName": "Jan",
			"lastName":  "Doe",
			"email":     "jan@example.com",
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

	t.Run("Should Fail If Required Fields Are Missing", func(t *testing.T) {
		requiredFields := []string{"firstName", "lastName", "email", "password"}
		for _, field := range requiredFields {
			reqBody := map[string]string{
				"firstName": "John",
				"lastName":  "Doe",
				"email":     "john@example.com",
				"password":  "password123",
			}
			delete(reqBody, field)
			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

			var result map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&result)
			assert.NoError(t, err)

			errors, ok := result["errors"].([]interface{})
			assert.True(t, ok, "Expected 'errors' field in response to be a slice")

			errorFound := false
			for _, err := range errors {
				errObj, ok := err.(map[string]interface{})
				if !ok {
					continue
				}
				if errMsg, exists := errObj["message"].(string); exists && errMsg == "Field is required" {
					if fieldName, exists := errObj["field"].(string); exists && fieldName == field {
						errorFound = true
						break
					}
				}
			}
			assert.True(t, errorFound, "Expected error message not found for field: "+field)
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

