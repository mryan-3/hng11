package controller

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mryan-3/hng11/stage2/database"
	"github.com/mryan-3/hng11/stage2/models"
	"github.com/mryan-3/hng11/stage2/utils"
	"github.com/mryan-3/hng11/stage2/validation"
)

// Create User
// route POST /api/v1/user/auth/register
func CreateUser(c *fiber.Ctx) error {
	type ReqBody struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Phone     string `json:"phone"`
	}

	body := new(ReqBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	validationErrors := validation.ValidateUser(body)

	if len(validationErrors) > 0 {
		response := fiber.Map{"errors": validationErrors}
		return c.Status(http.StatusUnprocessableEntity).JSON(response)
	}

	// hash password
	hashedPassword, hashingError := utils.CreateHashFromText(body.Password, 10)

	if hashingError != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "error",
			"message": "An error occurred while creating account",
		})
	}

	orgName := body.FirstName + "'s" + " Organisation"

	// Create an organisation
	org := models.Organisation{
		Name: orgName,
	}
	// Create organization and associate user
	if err := database.DB.Db.Create(&org).Error; err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to create organization")
	}

	// Create user
	user := models.User{
		FirstName:     body.FirstName,
		LastName:      body.LastName,
		Email:         body.Email,
		Password:      hashedPassword,
		Phone:         body.Phone,
		Organisations: []*models.Organisation{&org},
	}

	// Create user and organisation
	result := database.DB.Db.Create(&user)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"status":  "error",
				"message": "User already exists",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON("An error occurred while creating account")
	}


    // Generate token
    token, err := utils.SignJwtToken(user.UserID.String())

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "error",
			"message": "An error occurred while creating account",
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
        "data":  fiber.Map{
            "accessToken": token,
            "user": fiber.Map{
               "userId": user.UserID,
               "firstName": user.FirstName,
               "lastName": user.LastName,
               "email": user.Email,
               "phone": user.Phone,
            },
        },
    }

    return c.Status(http.StatusCreated).JSON(response)
}
