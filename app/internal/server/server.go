package server

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
)

type Server struct {
	Handler command.Handler
}

func NewServer(handler command.Handler) Server {
	return Server{
		Handler: handler,
	}
}

func (s *Server) Start(addr string) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}
