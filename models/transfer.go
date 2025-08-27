package models

import "time"

type Transfer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FromID    uint      `json:"from_id"`
	ToID      uint      `json:"to_id"`
	Amount    int64     `json:"amount"`
	Memo      string    `json:"memo"`
	CreatedAt time.Time `json:"created_at"`
}
