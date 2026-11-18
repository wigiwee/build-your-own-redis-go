package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleDiscardCommand(client *client.Client, args []string) {

	if len(args) > 1 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}

	if client.InTransaction == false {
		client.WriteResponse(errors.New("ERR DISCARD without MULTI"))
		return
	}

	client.InTransaction = false
	client.QueuedCommands = nil
	client.WatchedKeys = make(map[string]uint64)
	client.WriteSimpleString("OK")
	return
}
