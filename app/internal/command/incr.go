package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleIncrCommand(client *client.Client, args []string) {
	if len(args) != 2 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	output, err := h.store.Incr(args[1])
	if err != nil {
		client.WriteResponse(err)
		return
	}
	client.WriteResponse(output)
	return
}
