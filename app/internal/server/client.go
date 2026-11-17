package server

import (
	"errors"
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	client := client.NewClient(conn)

	for {
		args, err := client.ReadCommand()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("error reading from client, connection closed")
				return
			}
			log.Println("error reading from client")
			return
		}
		log.Println("processing command ", args)

		s.Handler.HandleCommand(client, args)
	}
}
