package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handlePsyncCommand(client *client.Client, args []string) {

	client.WriteSimpleString("FULLRSYNC " + h.config.ReplicationConfig.Master_replid + " 0")
}
