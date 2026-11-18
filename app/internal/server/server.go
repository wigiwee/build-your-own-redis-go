package server

import (
	"net"
	"strconv"

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

func (s *Server) Start(port int) error {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
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
