package command

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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
	fmt.Println("i was here")
	commands := client.QueuedCommands
	client.QueuedCommands = nil
	client.InTransaction = false

	originalWriter := client.Writer

	var byteBuffer bytes.Buffer
	tempWriter := bufio.NewWriter(&byteBuffer)

	client.SetWriter(tempWriter)

	for _, args := range commands {
		fmt.Println("did go in loop")
		h.executeCommand(client, args)
	}
	fmt.Println("Loop over")
	if err := tempWriter.Flush(); err != nil {
		fmt.Println("flushed and setting orignal writer")
		client.SetWriter(originalWriter)
		client.WriteResponse(err)
		return
	}
	client.SetWriter(originalWriter)
	fmt.Println("writing final response")
	//why is the application not writing respose after this
	client.WriteRaw("*" + strconv.Itoa(len(commands)) + utils.CRLF + byteBuffer.String())
	return
}
