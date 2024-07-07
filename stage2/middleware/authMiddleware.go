package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mryan-3/hng11/stage2/database"
	"github.com/mryan-3/hng11/stage2/models"
	"github.com/mryan-3/hng11/stage2/utils"
)

// Allow only authenticated user
func UserAuth(c *fiber.Ctx) error {
	// Get jwt from cookie
	userJwt := c.Cookies("user")

	if userJwt == "" {
		return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	// Verify jwt
	userId, isValid, err := utils.VerifyJwtToken(userJwt)

	if err != nil || !isValid {
		return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	var user models.User
	result := database.DB.Db.First(&user, "user_id=?", userId["user_id"])

	if result.Error != nil {
		return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})

	}

	// Set user id in context
	c.Locals("userId", user.UserID.String())

	return c.Next()

}

