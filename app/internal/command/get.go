package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleGetCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid command")
	}
	output := h.store.Get(args[1])
	respOutput, err := utils.Serialize(output)
	if err != nil {
		return err
	}
	writer.Write([]byte(respOutput))
	writer.Flush()
	return nil
}
