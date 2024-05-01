package models

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type Status string

const (
	PENDING   Status = "pending"
	SUBMITTED Status = "submitted"
)

func (st *Status) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		*st = Status(b)
	}
	return nil
}

func (st Status) Value() (driver.Value, error) {
	return string(st), nil
}

type Job struct {
	gorm.Model
	Id     string
	Status Status
}
