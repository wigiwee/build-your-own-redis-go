package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleXAddCommand(writer *bufio.Writer, args []string) error {
	if len(args) < 5 {
		return errors.New("invalid command")
	}

	output, err := h.store.AddStream(args[1], args[2], args[3:])
	if err != nil {
		return err
	}
	bulkStr, err := utils.SerializeBulkString(output)
	writer.Write([]byte(bulkStr))
	writer.Flush()
	return nil
}
