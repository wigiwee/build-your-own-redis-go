package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleTypeCommand(client *client.Client, args []string) {
	if len(args) != 2 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}

	keyType := h.store.GetType(args[1])
	client.WriteSimpleString(keyType)
	return
}
