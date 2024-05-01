package server

import (
	"fmt"
	"net"
)

type Server struct {
	Listener net.Listener
}

func Listen(host string, port string) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		return nil, err
	}

	server := Server{Listener: listener}

	return &server, nil
}

func (s *Server) Close() {
	s.Listener.Close()
}
