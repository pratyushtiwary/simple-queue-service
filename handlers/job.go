package handlers

import (
	"pratyushtiwary/sqs/queue"
	"pratyushtiwary/sqs/server"
)

// Handles the action of returning first job from the queue
//
// Action: job,
// fetches detail for the first queued job which is still pending
// the job returned becomes submitted and is assumed to be processed
func Job(request server.Request, queue *queue.Queue) (*server.Response, error) {
	response := new(server.Response)

	job, err := queue.Job()

	if err != nil {
		return nil, err
	}

	response.Data = map[string]any{
		"Id":     job.Id,
		"Status": job.Status,
		"Data":   job.Detail.Data,
	}

	return response, nil
}
