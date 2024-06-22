package handlers

import (
	"pratyushtiwary/sqs/queue"
	"pratyushtiwary/sqs/server"
)

// Handles the action of listing jobs
//
// Action: Jobs,
// returns list of jobs that are available,
// this includes queued and completed jobs
func JobsHandler(request server.Request, queue *queue.Queue) (*server.Response, error) {
	response := server.Response{}

	jobs, err := queue.Jobs()

	if err != nil {
		return nil, err
	}

	response.Data = map[string]any{
		"Jobs": jobs,
	}

	return &response, nil
}
