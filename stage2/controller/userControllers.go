package controller

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mryan-3/hng11/stage2/database"
	"github.com/mryan-3/hng11/stage2/models"
	"github.com/mryan-3/hng11/stage2/utils"
	"github.com/mryan-3/hng11/stage2/validation"
	"golang.org/x/crypto/bcrypt"
)

// Create User
// route POST /auth/register
func CreateUser(c *fiber.Ctx) error {
	type ReqBody struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
		Phone     string `json:"phone"`
	}

	fmt.Println("Creating user ...")

	body := new(ReqBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":     "Bad Request",
			"statusCode": http.StatusBadRequest,
			"message":    "Registration unsuccessful",
		})
	}

	validationErrors := validation.ValidateStruct(body)

	if len(validationErrors) > 0 {
		response := fiber.Map{"errors": validationErrors}
		return c.Status(http.StatusUnprocessableEntity).JSON(response)
	}

	// hash password
	hashedPassword, hashingError := utils.CreateHashFromText(body.Password, 10)

	if hashingError != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "error",
			"message": "An error occcured while hashing password",
		})
	}

	orgName := body.FirstName + "'s" + " Organisation"

	// Create an organisation
	org := models.Organisation{
		Name: orgName,
	} // Create organization and associate user
	if err := database.DB.Db.Create(&org).Error; err != nil {
        fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":     "Bad Request",
			"statusCode": http.StatusBadRequest,
			"message":    "Failed to create org",
		})
	}

	// Create user
	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  hashedPassword,
		Phone:     body.Phone,
	}

	// Create user and organisation
	result := database.DB.Db.Create(&user)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"status":     "Bad Request",
				"statusCode": http.StatusBadRequest,
				"message":    "Registration unsuccessful",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON("An error occurred while creating user")
	}

	// Add user to organisation
	database.DB.Db.Model(&org).Association("Users").Append(&user)

	// Add organisation to user
	database.DB.Db.Model(&user).Association("Organisations").Append(&org)

	// Generate token
	token, err := utils.SignJwtToken(user.UserID.String())

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "error",
			"message": "An error occurred while generating token!",
		})

	}

	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "user"
	cookie.Value = token
	cookie.HTTPOnly = utils.Check(os.Getenv("APP_ENV") == "prod", true, false)
	cookie.SameSite = utils.Check(os.Getenv("APP_ENV") == "prod", "strict", "None")
	cookie.Secure = utils.Check(os.Getenv("APP_ENV") == "prod", true, false)
	cookie.Expires = time.Now().Add(24 * time.Hour * 7)

	// Set cookie
	c.Cookie(cookie)

	response := fiber.Map{
		"status":  "success",
		"message": "Regstration successful",
		"data": fiber.Map{
			"accessToken": token,
			"user": fiber.Map{
				"userId":    user.UserID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"phone":     user.Phone,
			},
		},
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// Log in a user
// route POST /auth/login
func LoginUser(c *fiber.Ctx) error {
	type ReqBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	body := new(ReqBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": http.StatusBadRequest,
		})
	}

	validationErrors := validation.ValidateStruct(body)

	if len(validationErrors) > 0 {
		response := fiber.Map{"errors": validationErrors}
		return c.Status(http.StatusUnprocessableEntity).JSON(response)
	}

	var user models.User

	result := database.DB.Db.First(&user, "email = ?", body.Email)
    fmt.Println(result.Error)

    if result.Error != nil {
        return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
            "status":     "Bad request",
            "message":    "Authentication failed",
            "statusCode": http.StatusUnauthorized,
        })
    }

	// Compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": http.StatusUnauthorized,
		})

	}

	// Generate JWT token
	token, err := utils.SignJwtToken(user.UserID.String())

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "error",
			"message": "An error occurred while generating token!",
		})

	}

	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "user"
	cookie.Value = token
	cookie.HTTPOnly = utils.Check(os.Getenv("APP_ENV") == "prod", true, false)
	cookie.SameSite = utils.Check(os.Getenv("APP_ENV") == "prod", "strict", "None")
	cookie.Secure = utils.Check(os.Getenv("APP_ENV") == "prod", true, false)
	cookie.Expires = time.Now().Add(24 * time.Hour * 7)

	// Set cookie
	c.Cookie(cookie)

	response := fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data": fiber.Map{
			"accessToken": token,
			"user": fiber.Map{
				"userId":    user.UserID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"phone":     user.Phone,
			},
		},
	}

	return c.Status(http.StatusOK).JSON(response)
}

// Get a user
// route GET /api/users/:id
func GetUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	if userId == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":  "error",
			"message": "Missing Param",
		})
	}

	var user models.User

	if err := database.DB.Db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "User not found",
		})
	}

	response := fiber.Map{
		"status":  "success",
		"message": "User found",
		"data": fiber.Map{
			"user": fiber.Map{
				"userId":    user.UserID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"phone":     user.Phone,
			},
		},
	}

	return c.Status(http.StatusOK).JSON(response)
}
