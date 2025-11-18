package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	// Using SQLite for simplicity in demo
	db, err := gorm.Open(sqlite.Open("students.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
