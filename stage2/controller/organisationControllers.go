package controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mryan-3/hng11/stage2/database"
	"github.com/mryan-3/hng11/stage2/models"
	"github.com/mryan-3/hng11/stage2/validation"
)

// Get a users organisations
// route GET /api/organisations
func GetUserOrganisations(c *fiber.Ctx) error {

	type OrganizationResponse struct {
		OrgID       string `json:"orgId"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	userIdString := c.Locals("userId").(string)

	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":  "error",
			"message": "Error parsing the doctor ID",
		})
	}

	var user models.User
	if err := database.DB.Db.Preload("Organisations").First(&user, "user_id = ?", userId).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "User not found",
		})
	}

	var organizationsResponse []OrganizationResponse

	for _, org := range user.Organisations {
		organizationsResponse = append(organizationsResponse, OrganizationResponse{
			OrgID:       org.ID.String(),
			Name:        org.Name,
			Description: org.Description,
		})
	}
	response := fiber.Map{
		"status":  "success",
		"message": "Organisations found",
		"data": fiber.Map{
			"organisations": organizationsResponse,
		},
	}

	return c.Status(http.StatusOK).JSON(response)

}

// Get a single organisation
// route GET /api/organisations/:orgId
func GetSingleOrganisation(c *fiber.Ctx) error {
	type OrganizationResponse struct {
		OrgID       string `json:"orgId"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}
	orgId := c.Params("orgId")

	if orgId == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":  "error",
			"message": "Missing Param",
		})
	}

	var org models.Organisation
	if err := database.DB.Db.Where("id = ?", orgId).First(&org).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "Organisation not found",
		})
	}

	organizationResponse := OrganizationResponse{
		OrgID:       org.ID.String(),
		Name:        org.Name,
		Description: org.Description,
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Organisation found",
		"data":    organizationResponse,
	}

	return c.Status(http.StatusOK).JSON(response)
}

// Create an organisation
// route POST /api/organisations
func CreateOrganisation(c *fiber.Ctx) error {
	type ReqBody struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" `
	}

	// Find the user who created the organisation
	userId := c.Locals("userId").(string)
	var user models.User
	if err := database.DB.Db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "User not found",
		})
	}

	body := new(ReqBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":     "Bad request",
			"message":    "Client error",
			"statusCode": http.StatusBadRequest,
		})
	}

	validationErrors := validation.ValidateStruct(body)

	if len(validationErrors) > 0 {
		response := fiber.Map{"errors": validationErrors}
		return c.Status(http.StatusUnprocessableEntity).JSON(response)
	}

	org := models.Organisation{
		Name:        body.Name,
		Description: body.Description,
	}

	if err := database.DB.Db.Create(&org).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":     "Bad request",
			"message":    "Client error",
			"statusCode": http.StatusBadRequest,
		})
	}

	// Add user to organisation
	database.DB.Db.Model(&user).Association("Organisations").Append(&org)

	fmt.Println("User added to organisation", user.Organisations)

	// Add organisation to user
	database.DB.Db.Model(&org).Association("Users").Append(&user)

	fmt.Println("Organisation added to user", org.Users)

	response := fiber.Map{
		"status":  "success",
		"message": "Organisation created successfully",
		"data": fiber.Map{
			"orgId":       org.ID.String(),
			"name":        org.Name,
			"description": org.Description,
		},
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// Add a user to a particular organisation
// route POST /api/organisations/:orgId/users
func AddUserToOrganisation(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	if orgId == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":  "error",
			"message": "Missing Param",
		})
	}

	type ReqBody struct {
		UserId string `json:"userId" `
	}

	body := new(ReqBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"status":     "Bad request",
			"message":    "Client error",
			"statusCode": http.StatusBadRequest,
		})
	}

	userId := body.UserId

	var org models.Organisation
	if err := database.DB.Db.Where("id = ?", orgId).First(&org).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "Organisation not found",
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

	// Add user to organisation
	database.DB.Db.Model(&user).Association("Organisations").Append(&org)

	// Add organisation to user
	database.DB.Db.Model(&org).Association("Users").Append(&user)

	response := fiber.Map{
		"status":  "success",
		"message": "User added to organisation successfully",
	}

	return c.Status(http.StatusOK).JSON(response)
}
