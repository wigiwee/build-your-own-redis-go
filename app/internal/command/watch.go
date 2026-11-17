package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleWatchCommand(client *client.Client, args []string) {
	if len(args) != 2 {
		client.WriteResponse("Invalid command")
		return
	}

	client.WriteSimpleString("OK")
}
