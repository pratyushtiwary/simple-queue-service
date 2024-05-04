package models

import "gorm.io/gorm"

type Detail struct {
	gorm.Model
	Id   string
	Data string
}
