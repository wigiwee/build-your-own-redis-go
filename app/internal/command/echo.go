package command

import (
	"bufio"
	"errors"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleEchoCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid input")
	}
	bulkstr, err := utils.SerializeBulkString(args[1])
	if err != nil {
		return err
	}
	log.Println("writing ", args[1])
	_, err = writer.Write([]byte(bulkstr))
	writer.Flush()
	return err
}
