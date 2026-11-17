package command

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
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

func (h *Handler) handleXReadCommand(client *client.Client, args []string) {
	if len(args) < 4 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	switch strings.ToUpper(args[1]) {
	case "STREAMS":
		if (len(args)-2)%2 != 0 {
			client.WriteResponse(errors.New("invalid command"))
			return
		}
		trueLen := len(args) - 2
		streamResultArr := []storage.StreamReadResult{}
		for i := 2; i < 2+(trueLen/2); i++ {
			j := i + trueLen/2

			output, err := h.store.ReadStream(args[i], args[j])
			if err != nil {
				client.WriteResponse(err)
				return
			}
			streamResultArr = append(streamResultArr, output)

		}
		client.WriteResponse(StreamReadResponse{Streams: streamResultArr})
		return

	case "BLOCK":
		if strings.ToUpper(args[3]) == "STREAMS" {
			milliSec, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				client.WriteResponse(utils.ErrInvalidCmd)
				return
			}
			var ttl time.Duration
			if milliSec == 0 {
				ttl = time.Duration(time.Duration(200000) * time.Hour)
			} else {
				ttl = time.Duration(time.Duration(milliSec) * time.Millisecond)
			}

			output, err := h.store.ReadStreamBlock(args[4], args[5], ttl)
			if err != nil {
				client.WriteResponse(err)
				return
			}
			if output.Entries == nil || len(output.Entries) == 0 {
				var arr []string
				client.WriteResponse(arr)
				return
			}

			client.WriteResponse(StreamReadResponse{Streams: []storage.StreamReadResult{output}})

		} else {
			client.WriteResponse(utils.ErrInvalidCmd)
			return
		}

	default:
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	return
}
