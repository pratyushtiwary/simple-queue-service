package handlers

import (
	"encoding/json"
	"pratyushtiwary/sqs/queue"
	"pratyushtiwary/sqs/server"
)

// Handles the action of listing jobs
//
// Action: add_job,
// adds a job to the queue
// used by producer to submit a job
func AddJob(request server.Request, queue *queue.Queue) (*server.Response, error) {
	data, err := json.Marshal(request.Data)

	response := server.Response{}

	if err != nil {
		return nil, err
	}

	job, err := queue.AddJob(data)

	if err != nil {
		return nil, err
	}

	response.Data = map[string]any{
		"Id": job.Id,
	}

	return &response, nil
}
