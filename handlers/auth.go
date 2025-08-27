package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"tmp-backend/db"
	"tmp-backend/models"
	"tmp-backend/utils"
)

type registerPayload struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Birthday  string `json:"birthday"` // ISO date
}

func Register(c *fiber.Ctx) error {
	var p registerPayload
	if err := c.BodyParser(&p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to hash")
	}

	// parse birthday
	b, _ := time.Parse("2006-01-02", p.Birthday)

	user := models.User{
		Email:     p.Email,
		Password:  string(hashed),
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Phone:     p.Phone,
		Birthday:  b,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		if err == gorm.ErrDuplicatedKey || err.Error() == "UNIQUE constraint failed: users.email" {
			return fiber.NewError(fiber.StatusBadRequest, "email already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
	}

	// hide password
	user.Password = ""
	return c.Status(fiber.StatusCreated).JSON(user)
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var p loginPayload
	if err := c.BodyParser(&p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	var user models.User
	if err := db.DB.First(&user, "email = ?", p.Email).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(fiber.Map{"token": token})
}

func Me(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", uid).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	user.Password = ""
	return c.JSON(user)
}
