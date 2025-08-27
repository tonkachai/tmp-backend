package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"tmp-backend/db"
	"tmp-backend/models"
)

type transferPayload struct {
	MemberCode string `json:"member_code"`
	Amount     int64  `json:"amount"`
	Memo       string `json:"memo"`
}

// Transfer to a user by member code
func Transfer(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var p transferPayload
	if err := c.BodyParser(&p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	var to models.User
	if err := db.DB.First(&to, "member_code = ?", p.MemberCode).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "recipient not found")
	}

	// In a real app we'd check balances, create ledger entries, etc. Here we just record a transfer.
	t := models.Transfer{
		FromID:    uid.(uint),
		ToID:      to.ID,
		Amount:    p.Amount,
		Memo:      p.Memo,
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&t).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create transfer")
	}
	return c.Status(fiber.StatusCreated).JSON(t)
}

// RecentContacts returns a list of recent contacts the user transferred to or received from
func RecentContacts(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var transfers []models.Transfer
	if err := db.DB.Where("from_id = ? OR to_id = ?", uid, uid).Order("created_at desc").Limit(20).Find(&transfers).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transfers")
	}

	// build a set of contact IDs preserving order
	contacts := []models.User{}
	seen := map[uint]bool{}
	for _, tr := range transfers {
		var contactID uint
		if tr.FromID == uid.(uint) {
			contactID = tr.ToID
		} else {
			contactID = tr.FromID
		}
		if contactID == 0 || seen[contactID] {
			continue
		}
		var u models.User
		if err := db.DB.First(&u, "id = ?", contactID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load contact")
		}
		u.Password = ""
		contacts = append(contacts, u)
		seen[contactID] = true
	}
	return c.JSON(contacts)
}

// SearchUserByMemberCode finds a user by member code
func SearchUserByMemberCode(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing query")
	}
	var u models.User
	if err := db.DB.First(&u, "member_code = ?", q).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	u.Password = ""
	return c.JSON(u)
}
