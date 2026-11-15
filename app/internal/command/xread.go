package command

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

type StreamReadResponse struct {
	Streams []storage.StreamReadResult
}

func (s StreamReadResponse) SerializeRESP() (string, error) {
	str, err := utils.Serialize(s.Streams)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (h *Handler) handleXReadCommand(writer *bufio.Writer, args []string) error {
	if len(args) < 4 {
		return errors.New("invalid command")
	}
	switch strings.ToUpper(args[1]) {
	case "STREAMS":
		if (len(args)-2)%2 != 0 {
			return errors.New("invalid command")
		}
		trueLen := len(args) - 2
		streamResultArr := []storage.StreamReadResult{}
		for i := 2; i < 2+(trueLen/2); i++ {
			j := i + trueLen/2

			output, err := h.store.ReadStream(args[i], args[j])
			if err != nil {
				return err
			}
			streamResultArr = append(streamResultArr, output)

		}
		outputStr, err := utils.Serialize(StreamReadResponse{Streams: streamResultArr})
		if err != nil {
			return err
		}
		writer.Write([]byte(outputStr))
		writer.Flush()
		return nil

	case "BLOCK":
		if strings.ToUpper(args[3]) == "STREAMS" {
			milliSec, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return utils.ErrInvalidCmd
			}
			var ttl time.Duration
			if milliSec == 0 {
				ttl = time.Duration(time.Duration(200000) * time.Hour)
			} else {
				ttl = time.Duration(time.Duration(milliSec) * time.Millisecond)
			}

			output, err := h.store.ReadStreamBlock(args[4], args[5], ttl)
			if err != nil {
				return err
			}
			if output.Entries == nil || len(output.Entries) == 0 {
				writer.Write([]byte("*-1\r\n"))
				writer.Flush()
				return nil
			}

			outputStr, err := utils.Serialize(StreamReadResponse{Streams: []storage.StreamReadResult{output}})
			writer.Write([]byte(outputStr))
			writer.Flush()

		} else {
			return utils.ErrInvalidCmd
		}

	default:
		return errors.New("invalid command")
	}
	return nil
}
