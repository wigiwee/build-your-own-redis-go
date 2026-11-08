package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		args, err := utils.DeserializeArray(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("error reading from client, connection closed")
				return
			}
			log.Println("error reading from client")
			return
		}
		writer := bufio.NewWriter(conn)

		log.Println("processing command ", args)
		s.Handler.HandleCommand(writer, args)
	}
}
