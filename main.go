package main

import (
	"fmt"
	"os"
	"pratyushtiwary/sqs/handlers"
	"pratyushtiwary/sqs/queue"
	"pratyushtiwary/sqs/server"
	"strconv"
)

func main() {

	port, foundPort := os.LookupEnv("SQS_PORT")
	bufferSize, foundBufferSize := os.LookupEnv("SQS_BUFFER_SIZE")
	timeout, foundTimeout := os.LookupEnv("SQS_TIMEOUT") // in seconds
	verbose, foundVerbose := os.LookupEnv("SQS_VERBOSE") // 0 -> false, any other value -> true
	host := "localhost"

	if !foundPort {
		port = "4500"
	}

	if !foundBufferSize {
		bufferSize = "4096"
	}

	if !foundTimeout {
		timeout = "5"
	}

	if !foundVerbose {
		verbose = "0"
	}

	fmt.Println("Starting Simple Queue Service...")
	fmt.Println("Configuration:")
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Buffer Size: %s\n", bufferSize)
	fmt.Printf("Verbose: %s\n", verbose)
	fmt.Printf("Timeout (in seconds): %s\n", timeout)

	bufferSizeInt, err := strconv.Atoi(bufferSize)

	if err != nil {
		panic(err)
	}

	timeoutInt, err := strconv.Atoi(timeout)

	if err != nil {
		panic(err)
	}

	verboseBool, err := strconv.ParseBool(verbose)

	if err != nil {
		panic(err)
	}

	config := server.ServerConfig{
		BufferSize: bufferSizeInt,
		Timeout:    timeoutInt,
		Host:       host,
		Port:       port,
		Verbose:    verboseBool,
	}

	q, err := queue.Init()

	if err != nil {
		panic(err)
	}

	s, err := server.Listen(config, q)

	if err != nil {
		panic(err)
	}

	defer s.Close()

	fmt.Println("Server started successfully!")

	s.SetHandler("jobs", handlers.JobsHandler)
	s.SetHandler("add_job", handlers.AddJob)
	s.SetHandler("get_job", handlers.GetJob)
	s.SetHandler("job", handlers.Job)

	for {
		// Accept incoming connections
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go s.HandleRequest(conn)
	}

}
