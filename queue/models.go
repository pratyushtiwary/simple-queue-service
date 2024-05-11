package queue

import (
	"database/sql/driver"
	"time"
)

// 1:1 representation of tables in db
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

type Detail struct {
	Id   string
	Data string
}

type Job struct {
	Id        string
	Status    Status
	CreatedAt time.Time
	Detail    *Detail
}
