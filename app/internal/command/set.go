package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleSetCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 3 {
		return errors.New("invalid command")
	}
	h.store.Set(args[1], args[2])
	writer.Write([]byte(utils.OK_BULK_STR))
	writer.Flush()
	return nil
}
