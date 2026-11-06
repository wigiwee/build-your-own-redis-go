package server

import (
	"log"
	"net"
)

type server struct {
}

func NewServer() server {
	return server{}
}

func (s *server) Start(addr string) error {

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

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		conn.Write([]byte("+PONG\r\n"))
		log.Println("PONGed it")
	}

}
