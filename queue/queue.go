package queue

import (
	"pratyushtiwary/sqs/models"

	"github.com/google/uuid"
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

func (q *Queue) AddJob(data string) (*models.Job, error) {
	id := uuid.New().String()
	detailId := uuid.New().String()
	detail := models.Detail{
		Id:   detailId,
		Data: data,
	}
	job := models.Job{
		Id:       id,
		Status:   models.PENDING,
		DetailID: detail.Id,
		Detail:   detail,
	}

	result := q.DB.Create(&detail)

	if result.Error != nil {
		return nil, result.Error
	}

	result = q.DB.Create(&job)

	if result.Error != nil {
		return nil, result.Error
	}

	return &job, nil
}
