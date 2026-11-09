package command

import (
	"bufio"
	"errors"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleBLPopCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 3 {
		return errors.New("invalid command")
	}
	n, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	ttl := time.Duration(n) * time.Second

	output := h.store.BLPop(args[1], ttl)
	respArr, err := utils.SerializeArray(output)
	if err != nil {
		return err
	}
	writer.Write([]byte(respArr))
	writer.Flush()
	return nil
}
