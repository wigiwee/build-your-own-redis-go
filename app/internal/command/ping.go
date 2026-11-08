package command

import (
	"bufio"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handlePingCommand(writer *bufio.Writer, args []string) error {
	pongBulkStr, err := utils.SerializeString("PONG")

	if err != nil {
		log.Println(err)
		log.Println("failed here 1")
		return err
	}
	log.Println("writing pong")
	_, err = writer.Write([]byte(pongBulkStr))
	if err != nil {
		log.Println("failed here 2")
		return err
	}
	writer.Flush()
	return nil
}
