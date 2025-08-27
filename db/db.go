package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"tmp-backend/models"
)

var DB *gorm.DB

func Init() {
	// place sqlite file next to the binary / project root
	dbPath := filepath.Join(".", "data.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// ensure file exists when GORM opens it
		f, _ := os.Create(dbPath)
		f.Close()
	}

	d, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	DB = d

	// Migrate the schema
	if err := DB.AutoMigrate(&models.User{}, &models.Transfer{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// fill missing member codes (simple sequential code)
	var users []models.User
	if err := DB.Find(&users).Error; err == nil {
		for i, u := range users {
			if u.MemberCode == "" {
				code := fmt.Sprintf("LBK%05d", u.ID)
				u.MemberCode = code
				DB.Model(&u).Update("member_code", u.MemberCode)
			}
			_ = i
		}
	}
}
