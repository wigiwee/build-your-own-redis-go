package command

import (
	"bufio"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func handlePingCommand(writer *bufio.Writer, args []string) error {
	pongBulkStr, err := utils.SerializeString("PONG")

	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("writing pong")
	_, err = writer.Write([]byte(pongBulkStr))
	writer.Flush()
	return err
}
