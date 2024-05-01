package main

import (
	"fmt"
	"net"
	"os"
	"pratyushtiwary/sqs/server"
)

func main() {

	port, foundPort := os.LookupEnv("SQS_PORT")
	host := "localhost"

	if !foundPort {
		port = "4500"
	}

	fmt.Println("Starting Simple Queue Service...")
	fmt.Println("Configuration:")
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("Port: %s\n", port)

	s, err := server.Listen(host, port)

	if err != nil {
		panic(err)
	}

	defer s.Close()

	fmt.Println("Server started successfully!")

	for {
		// Accept incoming connections
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

		// Handle client connection in a goroutine
		go handleClient(conn)
	}

}

// move to queue package and let it handle the requests
func handleClient(conn net.Conn) {
	defer conn.Close()
}
