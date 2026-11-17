package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleMultiCommand(writer *bufio.Writer, args []string) error {
	if len(args) > 1 {
		return errors.New("invalid command")
	}
	writer.Write([]byte(utils.OK_BULK_STR))
	writer.Flush()
	return nil
}
