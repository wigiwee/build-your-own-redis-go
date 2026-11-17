package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleLpushCommand(client *client.Client, args []string) {
	if !(len(args) > 2) {
		client.WriteResponse(errors.New("invalid command"))
		return
	}

	resInt := h.store.LPushList(args[1], args[2:])
	client.WriteResponse(resInt)
	return
}
