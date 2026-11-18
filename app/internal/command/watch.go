package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleWatchCommand(client *client.Client, args []string) {
	if len(args) != 2 {
		client.WriteResponse("Invalid command")
		return
	}
	if client.InTransaction {
		client.WriteResponse(errors.New("ERR WATCH inside MULTI is not allowed"))
		return
	}

	for _, key := range args[1:] {
		client.WatchedKeys[key] = h.store.GetVersion(key)
	}
	client.WriteSimpleString("OK")
}
