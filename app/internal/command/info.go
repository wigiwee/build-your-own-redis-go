package command

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleInfoCommand(client *client.Client, args []string) {

	switch strings.ToUpper(args[1]) {
	case "REPLICATION":
		var sb strings.Builder
		sb.WriteString("role:" + h.config.ReplicationConfig.Role + "\n")
		sb.WriteString("master_replid:" + h.config.ReplicationConfig.Master_replid + "\n")
		sb.WriteString("master_repl_offset:" + strconv.Itoa(h.config.ReplicationConfig.Master_repl_offset) + "\n")
		client.WriteResponse(sb.String())
	}
}
