package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MemberCode string    `gorm:"uniqueIndex;size:16" json:"member_code"`
	Email      string    `gorm:"uniqueIndex;not null" json:"email"`
	Password   string    `json:"-"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Phone      string    `json:"phone"`
	Birthday   time.Time `json:"birthday"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
