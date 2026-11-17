package command

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleExecCommand(client *client.Client, args []string) {
	if len(args) > 1 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	if !client.InTransaction {
		client.WriteResponse(errors.New("ERR EXEC without MULTI"))
		return
	}

	commands := client.QueuedCommands
	client.QueuedCommands = nil
	client.InTransaction = false
	originalWriter := client.Writer

	var byteBuffer bytes.Buffer
	tempWriter := bufio.NewWriter(&byteBuffer)

	client.SetWriter(tempWriter)

	for _, args := range commands {
		h.HandleCommand(client, args)
	}
	client.SetWriter(originalWriter)
	client.WriteRaw("*" + strconv.Itoa(len(commands)) + utils.CRLF + byteBuffer.String())
	return
}
