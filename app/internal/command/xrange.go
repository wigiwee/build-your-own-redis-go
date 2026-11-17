package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleXRangeCommand(client *client.Client, args []string) {
	if len(args) < 4 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}

	output, err := h.store.RangeStream(args[1], args[2], args[3])
	if err != nil {
		client.WriteResponse(err)
		return
	}
	client.WriteResponse(output)
	return
}
