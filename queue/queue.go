package queue

import (
	"database/sql"
	"errors"
	"log"
	"pratyushtiwary/sqs/queue/queries"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Queue struct {
	DB *sql.DB
}

func Init() (*Queue, error) {
	db, err := sql.Open("sqlite3", "./queue.db")
	if err != nil {
		log.Fatal(err)
	}

	queue := Queue{
		DB: db,
	}

	// let's create tables

	// first we'll create details, cause jobs is dependent on it
	_, err = db.Exec(queries.DETAIL_CREATE_QUERY)
	if err != nil {
		log.Printf("%q: %s\n", err, queries.DETAIL_CREATE_QUERY)
	}

	// let's create jobs table
	_, err = db.Exec(queries.JOB_CREATE_QUERY)
	if err != nil {
		log.Printf("%q: %s\n", err, queries.JOB_CREATE_QUERY)
	}

	return &queue, nil
}

func (q *Queue) AddJob(data []byte) (*Job, error) {
	id := uuid.New().String()
	detailId := uuid.New().String()

	// start a transaction
	tx, err := q.DB.Begin()

	defer tx.Commit()

	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(queries.INSERT_DETAIL_QUERY)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(
		detailId,
		data,
	)

	stmt.Close()

	if err != nil {
		return nil, err
	}

	stmt, err = tx.Prepare(queries.INSERT_JOB_QUERY)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(
		id,
		detailId,
	)

	stmt.Close()

	if err != nil {
		return nil, err
	}

	detail := Detail{
		Id:   detailId,
		Data: string(data),
	}

	job := Job{
		Id:        id,
		Status:    PENDING,
		Detail:    &detail,
		CreatedAt: time.Now(),
	}

	return &job, nil
}

func (q *Queue) getDetail(detailId string) (*Detail, error) {
	detail := new(Detail)
	var data []byte

	rows, err := q.DB.Query(queries.SELECT_DETAIL_QUERY, detailId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	valid := rows.Next()

	if !valid {
		return nil, errors.New("failed to find details for job")
	}

	err = rows.Scan(&data)

	if err != nil {
		return nil, err
	}

	detail.Data = string(data)
	detail.Id = detailId

	return detail, nil
}

func (q *Queue) GetJob(id string) (*Job, error) {
	job := new(Job)

	job.Id = id

	rows, err := q.DB.Query(queries.SELECT_JOB_QUERY, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	valid := rows.Next()

	if !valid {
		return nil, errors.New("invalid job id")
	}

	var status Status
	var createdAt time.Time
	var detailId string

	err = rows.Scan(&status, &detailId, &createdAt)

	if err != nil {
		return nil, err
	}

	job.Status = status
	job.CreatedAt = createdAt

	detail, err := q.getDetail(detailId)

	if err != nil {
		return nil, err
	}

	job.Detail = detail

	return job, nil
}

func (q *Queue) Job() (*Job, error) {
	job := new(Job)

	rows, err := q.DB.Query(queries.SELECT_FIRST_QUEUED_JOB_QUERY, PENDING)
	if err != nil {
		log.Fatal(err)
	}
	valid := rows.Next()

	if !valid {
		return nil, errors.New("no pending jobs found")
	}

	var status Status = SUBMITTED
	var id string
	var detailId string
	var createdAt time.Time

	err = rows.Scan(&id, &detailId, &createdAt)

	if err != nil {
		return nil, err
	}

	job.Status = status
	job.CreatedAt = createdAt
	job.Id = id

	detail, err := q.getDetail(detailId)

	if err != nil {
		return nil, err
	}

	job.Detail = detail

	rows.Close()

	// update the job's status

	// start a transaction
	tx, err := q.DB.Begin()

	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(queries.UPDATE_JOB_STATUS_QUERY)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(
		status,
		id,
	)

	stmt.Close()

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return job, nil
}
