package server

import (
	"net"
)

type Server struct {
}

func NewServer() Server {
	return Server{}
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
