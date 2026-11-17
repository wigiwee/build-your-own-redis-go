package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleLlenCommand(client *client.Client, args []string) {
	if len(args) != 2 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}

	resInt := h.store.GetListLen(args[1])
	client.WriteResponse(resInt)
	return
}
