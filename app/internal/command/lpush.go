package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleLpushCommand(writer *bufio.Writer, args []string) error {
	if !(len(args) > 2) {
		return errors.New("invalid command")
	}
	resInt := h.store.LPushList(args[1], args[2:])
	serInt, err := utils.Serialize(resInt)
	if err != nil {
		return err
	}
	writer.Write([]byte(serInt))
	writer.Flush()
	return nil
}
