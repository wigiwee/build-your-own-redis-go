package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleGetCommand(client *client.Client, args []string) {
	if len(args) != 2 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	output := h.store.Get(args[1])
	client.WriteResponse(output)
	return
}
