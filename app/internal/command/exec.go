package command

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	watchedKey := client.WatchedKey
	client.WatchedKey = ""
	client.QueuedCommands = nil
	client.InTransaction = false

	fmt.Println(commands)
	fmt.Println(watchedKey)

	originalWriter := client.Writer

	var byteBuffer bytes.Buffer
	tempWriter := bufio.NewWriter(&byteBuffer)

	client.SetWriter(tempWriter)

	for _, args := range commands {
		if watchedKey != "" {
			if strings.ToUpper(args[0]) == "SET" && args[1] == watchedKey {
				fmt.Println("in the condition")
				var emptArr []string
				client.Writer = originalWriter
				client.WriteResponse(emptArr)
				return
			}
		}
		h.executeCommand(client, args)
	}
	if err := tempWriter.Flush(); err != nil {
		client.SetWriter(originalWriter)
		client.WriteResponse(err)
		return
	}
	client.SetWriter(originalWriter)
	//why is the application not writing respose after this
	client.WriteRaw("*" + strconv.Itoa(len(commands)) + utils.CRLF + byteBuffer.String())
	return
}
