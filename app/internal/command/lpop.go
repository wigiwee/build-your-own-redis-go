package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleLPopCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid command")
	}

	pop := h.store.LPop(args[1])
	if len(pop) == 0 {
		writer.Write([]byte(utils.NULL_BULK_STR))
	} else {
		bulkstr, err := utils.SerializeBulkString(pop)
		if err != nil {
			return err
		}
		writer.Write([]byte(bulkstr))

	}
	writer.Flush()
	return nil
}
