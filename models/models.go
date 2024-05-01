package models

import "gorm.io/gorm"

func Init(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&Job{})
}
