package command

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleMultiCommand(client *client.Client, args []string) {
	if len(args) > 1 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	if client.InTransaction == true {
		client.WriteResponse(errors.New("client already in exec"))
		return
	}
	client.InTransaction = true
	client.WriteSimpleString("OK")
	return
}
