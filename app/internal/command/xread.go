package command

import (
	"bufio"
	"errors"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleXReadCommand(writer *bufio.Writer, args []string) error {
	if len(args) < 4 {
		return errors.New("invalid command")
	}

	//XREAD streams stream_key other_stream_key 0-0 0-1
	// 0      1			2		3				4    5

	outputArr := []any{}
	switch strings.ToUpper(args[1]) {
	case "STREAMS":
		if (len(args)-2)%2 != 0 {
			return errors.New("invalid command")
		}
		trueLen := len(args) - 2
		for i := 2; i < 2+(trueLen/2); i++ {
			j := i + trueLen/2

			output, err := h.store.ReadStream(args[i], args[j])
			if err != nil {
				return err
			}
			outputArr = append(outputArr, output)
		}
	default:
		return errors.New("invalid command")
	}
	outputStr, err := utils.Serialize(outputArr)
	if err != nil {
		return err
	}
	writer.Write([]byte(outputStr))
	writer.Flush()
	return nil
}
