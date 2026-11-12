package command

import (
	"bufio"
	"errors"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleLPopCommand(writer *bufio.Writer, args []string) error {
	if len(args) > 3 {
		return errors.New("invalid command")
	}
	n := 1
	if len(args) != 2 {
		var err error
		n, err = strconv.Atoi(args[2])
		if err != nil {
			return err
		}
	}
	pop := h.store.LPop(args[1], n)
	if len(pop) == 0 {
		writer.Write([]byte(utils.NULL_BULK_STR))
	} else if len(pop) == 1 {
		bulkstr, err := utils.Serialize(pop[0])
		if err != nil {
			return err
		}
		writer.Write([]byte(bulkstr))

	} else {
		respArr, err := utils.Serialize(pop)
		if err != nil {
			return err
		}
		writer.Write([]byte(respArr))

	}
	writer.Flush()
	return nil
}
