package command

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleInfoCommand(client *client.Client, args []string) {

	switch strings.ToUpper(args[1]) {
	case "REPLICATION":
		client.WriteResponse("role:" + h.config.ReplicationConfig.Role)
	}
}
