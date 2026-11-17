package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleEchoCommand(client *client.Client, args []string) {
	if len(args) != 2 {
		client.WriteResponse(errors.New("invalid input"))
		return
	}
	client.WriteResponse(args[1])
	return
}
