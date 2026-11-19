package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleReplConfCommand(client *client.Client, args []string) {
	// if len(args) != 3 {
	// 	client.WriteResponse(errors.New("Invalid command"))
	// 	return
	// }

	client.WriteSimpleString("OK")
}
