package command

import (
	"bufio"
	"errors"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
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
	outputStr, err := SerializeStreamEntries(output)
	if err != nil {
		return err
	}
	writer.Write([]byte(outputStr))
	writer.Flush()
	return nil
}

func SerializeStreamEntries(entries []storage.StreamEntry) (string, error) {
	if entries == nil {
		return "*-1\r\n", nil
	}

	var output strings.Builder

	output.WriteString("*")
	output.WriteString(strconv.Itoa(len(entries)))
	output.WriteString(utils.CRLF)

	for _, entry := range entries {
		// Each stream entry: [ID, [field, value, field, value...]]
		output.WriteString("*2")
		output.WriteString(utils.CRLF)

		id, err := utils.SerializeBulkString(entry.ID.ToString())
		if err != nil {
			return "", err
		}
		output.WriteString(id)

		output.WriteString("*")
		output.WriteString(strconv.Itoa(len(entry.Fields) * 2))
		output.WriteString(utils.CRLF)

		for field, value := range entry.Fields {
			serializedField, err := utils.SerializeBulkString(field)
			if err != nil {
				return "", err
			}

			serializedValue, err := utils.SerializeBulkString(value)
			if err != nil {
				return "", err
			}

			output.WriteString(serializedField)
			output.WriteString(serializedValue)
		}
	}

	return output.String(), nil
}
