package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleXReadCommand(writer *bufio.Writer, args []string) error {
	if len(args) < 4 {
		return errors.New("invalid command")
	}

	output, err := h.store.ReadStream(args[2], args[3])
	if err != nil {
		return err
	}
	outputStr, err := utils.Serialize([]any{output})
	if err != nil {
		return err
	}
	writer.Write([]byte(outputStr))
	writer.Flush()
	return nil
}
