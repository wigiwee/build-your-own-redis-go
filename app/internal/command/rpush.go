package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleRpushCommand(writer *bufio.Writer, args []string) error {
	if !(len(args) > 2) {
		return errors.New("invalid command")
	}
	resInt := h.store.RPushList(args[1], args[2:])
	writer.Write([]byte(utils.SerializeInt(resInt)))
	writer.Flush()
	return nil
}
