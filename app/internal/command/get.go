package command

import (
	"bufio"
	"errors"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleGetCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid command")
	}
	value, doesExist := h.store.Get(args[1])
	if !doesExist {
		writer.Write([]byte(utils.NULL_BULK_STR))
	} else {
		bulkStr, err := utils.SerializeBulkString(value)
		if err != nil {
			log.Printf("error serializing bulk str ", value)
		}
		writer.Write([]byte(bulkStr))
	}

	writer.Flush()
	return nil
}
