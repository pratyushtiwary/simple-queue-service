package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"

	"pratyushtiwary/sqs/queue"
)

type RequestHandler = func(request Request, queue *queue.Queue) (*Response, error)

type ServerConfig struct {
	BufferSize int
	Timeout    int
	Host       string
	Port       string
	Verbose    bool
}

type ServerStatus string

const (
	SUCCESS ServerStatus = "success"
	ERROR   ServerStatus = "error"
)

type Server struct {
	Listener     net.Listener
	Config       ServerConfig
	Queue        *queue.Queue
	mapping      map[string]RequestHandler
	infoLogger   *log.Logger
	errorLogger  *log.Logger
	infoLogFile  *os.File
	errorLogFile *os.File
}

type Response struct {
	Status ServerStatus
	Data   map[string]any
}

type Request struct {
	Action string
	Data   map[string]any
}

func Listen(config ServerConfig, queue *queue.Queue) (*Server, error) {
	logsBasePath := "logs"

	_, err := os.Stat(logsBasePath)

	if err != nil {
		panic(err)
	}

	infoLogPath := path.Join(logsBasePath, "queue.log")
	errorLogPath := path.Join(logsBasePath, "queue.err")

	infoLogFile, err := os.OpenFile(infoLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		panic(err)
	}

	errorLogFile, err := os.OpenFile(errorLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		panic(err)
	}

	infoLogger := log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port))

	if err != nil {
		return nil, err
	}

	server := Server{
		Listener:     listener,
		Config:       config,
		Queue:        queue,
		infoLogger:   infoLogger,
		errorLogger:  errorLogger,
		infoLogFile:  infoLogFile,
		errorLogFile: errorLogFile,
	}

	return &server, nil
}

func (s *Server) Close() {
	s.infoLogFile.Close()
	s.errorLogFile.Close()
	s.Listener.Close()
}

func (s *Server) Log(status ServerStatus, message ...any) {
	if s.Config.Verbose {
		fmt.Println(status, message)
	}

	if status == ERROR {
		s.errorLogger.Println(message...)
	} else {
		s.infoLogger.Println(message...)
	}
}

func (s *Server) HandleRequest(conn net.Conn) {
	defer conn.Close()

	s.Log(SUCCESS, "Client Connected: ", conn.RemoteAddr().String())

	for {
		data, err := s.readData(conn)

		if err != nil {
			s.Log(ERROR, "Error while reading data: ", err)
			break
		}

		if len(data) == 0 {
			continue
		}

		s.Log(SUCCESS, "Data recevied from client, parsing data...")

		// if client sends BYE then break out of this loop and close connection
		if string(data) == "BYE" {
			s.Log(SUCCESS, "Received BYE from client")
			break
		}

		parsedData, err := s.parseData(data)

		if err != nil {
			s.Log(ERROR, err)
			break
		}

		handler, err := s.getHandler(*parsedData)

		if err != nil {
			s.Log(ERROR, err)
			break
		}

		s.Log(SUCCESS, "Handler for data recevied fetched successfully!")

		response, err := handler(*parsedData, s.Queue)

		if err != nil {
			response = &Response{
				Status: ERROR,
				Data: map[string]any{
					"Detail": err.Error(),
				},
			}
		} else {
			if response.Status == "" {
				response.Status = SUCCESS // defaults to success
			}
		}

		transformedResponse, err := s.transformResponse(response)

		if err != nil {
			fmt.Println(err)
		}

		s.Log(SUCCESS, "Request handling done, sending result back to client...")
		fmt.Fprint(conn, string(transformedResponse))
	}

	s.Log(SUCCESS, "Closing connection with the client: ", conn.RemoteAddr().String())
}

func (s *Server) transformResponse(response *Response) ([]byte, error) {
	result, err := json.Marshal(response)

	if err != nil {
		return nil, err
	}

	return result, nil
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
		s.mapping = make(map[string]RequestHandler)
	}
	s.mapping[action] = handler
}
