package command

import (
	"bufio"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleTypeCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid command")
	}

	keyType := h.store.GetType(args[1])
	simpleStr, err := utils.SerializeString(keyType)
	if err != nil {
		return err
	}
	writer.Write([]byte(simpleStr))
	writer.Flush()
	return nil
}
