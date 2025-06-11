package handlers

import (
	"math"
	"strconv"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	service ports.UserService
}

func NewHttpUserHandler(userservice ports.UserService) *HttpUserHandler {
	return &HttpUserHandler{service: userservice}
}

func (h *HttpUserHandler) RegisterUser(c *fiber.Ctx) error {
	user := new(domain.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.RegisterUser(user); err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *HttpUserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(domain.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.UpdateUser(user, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user updated successfully",
	})
}

func (h *HttpUserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delete Account Successfully..",
	})
}

func (h *HttpUserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetUserById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": user})
}

func (h *HttpUserHandler) FindUsers(c *fiber.Ctx) error {
	user, err := h.service.FindUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	var result []domain.ResUserWaithID
	for _, u := range user {
		result = append(result, domain.ResUserWaithID{
			ID:         u.ID,
			FullName:   u.FirstName + " " + u.LastName,
			EmployeeID: u.EmployeeID,
			Email:      u.Email,
			Status:     string(u.Status),
			UpdatedAt:  u.UpdatedAt,
		})
	}
	return c.JSON(result)
}

func (h *HttpUserHandler) GetPaginationUsers(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	r, err := h.service.GetPaginationUsers(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(r)
}

func (h *HttpUserHandler) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	seacrch := c.Query("search", "")
	status := c.Query("status", "")

	filter := domain.UserFilter{
		Search: seacrch,
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	users, total, err := h.service.GetUsers(&filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	var result []domain.ResUSerNoID
	for _, u := range users {
		result = append(result, domain.ResUSerNoID{
			FullName:   u.FirstName + " " + u.LastName,
			EmployeeID: u.EmployeeID,
			Email:      u.Email,
			Status:     string(u.Status),
			UpdatedAt:  u.UpdatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"total_page": int(math.Ceil(float64(total) / float64(limit))),
		"data":       result,
	})
}

// >uploade profile user
func (h *HttpUserHandler) UploadProfilePicture(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required as query parameter or path parameter"})
	}

	file, err := c.FormFile("profile_image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to get file from form:" + err.Error()})
	}

	// อ่านไฟล์เป็น byte array
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file content:" + err.Error()})
	}

	buf := make([]byte, file.Size)
	_, err = fileContent.Read(buf)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file content:" + err.Error()})
	}

	profilePicURL, err := h.service.UploadProfilePicture(id, buf, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload profile picture:" + err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "profile picture uploaded successfully",
		"profilePicUrl": profilePicURL,
	})
}
