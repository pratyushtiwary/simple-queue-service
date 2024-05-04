package queue

import (
	"pratyushtiwary/sqs/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Queue struct {
	DB *gorm.DB
}

func Init() (*Queue, error) {
	db, err := gorm.Open(sqlite.Open("queue.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	models.Init(db)

	queue := new(Queue)

	queue.DB = db

	return queue, nil
}
