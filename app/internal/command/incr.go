package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleIncrCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid command")
	}
	output, err := h.store.Incr(args[1])
	if err != nil {
		return err
	}
	outputStr, err := utils.Serialize(output)
	if err != nil {
		return err
	}
	writer.Write([]byte(outputStr))
	writer.Flush()
	return nil
}
