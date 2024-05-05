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

	fmt.Println("Starting Simple Queue Service...")
	fmt.Println("Configuration:")
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Buffer Size: %s\n", bufferSize)
	fmt.Printf("Timeout (in seconds): %s\n", timeout)

	bufferSizeInt, err := strconv.Atoi(bufferSize)

	if err != nil {
		panic(err)
	}

	timeoutInt, err := strconv.Atoi(timeout)

	if err != nil {
		panic(err)
	}

	config := server.ServerConfig{
		BufferSize: bufferSizeInt,
		Timeout:    timeoutInt,
		Host:       host,
		Port:       port,
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

	for {
		// Accept incoming connections
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

		// Handle client connection in a goroutine
		go s.HandleRequest(conn)
	}

}
