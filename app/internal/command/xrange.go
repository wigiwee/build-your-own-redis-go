package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleXRangeCommand(writer *bufio.Writer, args []string) error {
	if len(args) < 4 {
		return errors.New("invalid command")
	}

	output, err := h.store.RangeStream(args[1], args[2], args[3])
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
