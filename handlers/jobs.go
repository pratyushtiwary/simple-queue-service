package handlers

import (
	"fmt"
	"pratyushtiwary/sqs/server"
)

// Handles the action of listing jobs
//
// Action: Jobs,
// returns list of jobs that are available,
// this includes queued and completed jobs
func JobsHandler(request server.Request) {
	fmt.Printf("Action: %s\n", request.Action)
}
