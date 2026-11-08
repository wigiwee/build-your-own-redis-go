package command

import (
	"bufio"
	"errors"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleLrangeCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 4 {
		return errors.New("invalid command")
	}
	start, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	stop, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}

	resArr := h.store.GetListRange(args[1], start, stop)
	resArrStr, err := utils.SerializeArray(resArr)
	if err != nil {
		return err
	}
	writer.Write([]byte(resArrStr))
	writer.Flush()
	return nil
}
