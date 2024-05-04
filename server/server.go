package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type RequestHandler = func(request Request)

type ServerConfig struct {
	BufferSize int
	Timeout    int
	Host       string
	Port       string
}

type Server struct {
	Listener net.Listener
	Config   ServerConfig
	mapping  map[string]RequestHandler
}

type Request struct {
	Action string
	Data   map[string]any
}

func Listen(config ServerConfig) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port))

	if err != nil {
		return nil, err
	}

	server := Server{Listener: listener, Config: config}

	return &server, nil
}

func (s *Server) Close() {
	s.Listener.Close()
}

func (s *Server) HandleRequest(conn net.Conn) {
	defer conn.Close()

	data, err := s.readData(conn)

	if err != nil {
		fmt.Printf("Error while reading data: %s\n", err)
	}

	parsedData, err := s.parseData(data)

	if err != nil {
		panic(err)
	}

	handler, err := s.getHandler(*parsedData)

	if err != nil {
		panic(err)
	}

	handler(*parsedData)
}

// Reads data from passed connection,
// returns bytearray, it can be empty if no data found
//
// Size of bytearray = Config.BufferSize
//
// Reading timeout occurs after Config.Timeout seconds
func (s *Server) readData(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, s.Config.BufferSize)
	size := 0

	for {
		err := conn.SetReadDeadline(time.Now().Add(time.Duration(s.Config.Timeout) * time.Second))
		if err != nil {
			return nil, err
		}

		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF && !errors.Is(err, os.ErrDeadlineExceeded) {
				return nil, err
			}
			break
		}
		if n == 0 {
			break
		}

		// Process and use the data (here, we'll just print it)
		size = n
	}

	fmt.Printf("Received: %s\n", buffer[:size])

	return buffer[:size], nil
}

// Parses data and returns a map
//
// Data is expected to be a bytearray in JSON format
func (s *Server) parseData(data []byte) (*Request, error) {
	var d Request
	err := json.Unmarshal(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

// Returns RequestHandler if found
func (s *Server) getHandler(request Request) (RequestHandler, error) {
	action := request.Action

	handler, ok := s.mapping[action]

	if ok {
		return handler, nil
	}
	err := fmt.Errorf("no handler configured for %s action", action)
	return nil, err
}

func (s *Server) SetHandler(action string, handler RequestHandler) {
	if s.mapping == nil {
		s.mapping = make(map[string]func(request Request))
	}
	s.mapping[action] = handler
}
