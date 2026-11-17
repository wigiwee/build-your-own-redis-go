package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handlePingCommand(client *client.Client, args []string) {

	client.WriteSimpleString("PONG")
	return
}
