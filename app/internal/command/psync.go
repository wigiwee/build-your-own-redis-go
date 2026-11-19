package command

import (
	"encoding/hex"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handlePsyncCommand(client *client.Client, args []string) {

	client.WriteSimpleString("FULLRESYNC " + h.config.ReplicationConfig.Master_replid + " 0")

	rdbByte, err := hex.DecodeString(utils.EMPTY_RDB_HEX)
	if err != nil {
		client.WriteResponse(err)
		return
	}
	client.WriteRaw("$" + strconv.Itoa(len(rdbByte)) + utils.CRLF + string(rdbByte))

}
