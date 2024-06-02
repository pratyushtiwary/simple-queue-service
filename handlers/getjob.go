package handlers

import (
	"pratyushtiwary/sqs/queue"
	"pratyushtiwary/sqs/server"
)

// Handles the action of listing specified job's detail
//
// Action: get_job,
// fetches detail for the job with provided id
// can be used by consumer to check job status
func GetJob(request server.Request, queue *queue.Queue) (*server.Response, error) {
	response := server.Response{}

	id := request.Data["Id"].(string)

	job, err := queue.GetJob(id)

	if err != nil {
		return nil, err
	}

	response.Data = map[string]any{
		"Detail": job,
	}

	return &response, nil
}
