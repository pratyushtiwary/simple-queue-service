package handlers

import (
	"fmt"
	"pratyushtiwary/sqs/queue"
	"pratyushtiwary/sqs/server"
)

// Handles the action of listing jobs
//
// Action: Jobs,
// returns list of jobs that are available,
// this includes queued and completed jobs
func JobsHandler(request server.Request, queue *queue.Queue) (*server.Response, error) {
	fmt.Printf("Action: %s\n", request.Action)
	response := server.Response{}

	response.Data = map[string]any{
		"a": 1,
	}

	response.Status = "success"

	return &response, nil
}
