package command

import (
	"bufio"
	"log"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func HandleCommand(writer *bufio.Writer, args []string) {

	if len(args) == 0 {
		log.Fatal("empty args received")
	}

	switch strings.ToUpper(args[0]) {

	case "PING":
		err := handlePingCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "ECHO":
		err := handleEchoCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	}

}

func WriteError(writer *bufio.Writer, errorStr string) {
	_, err := writer.Write([]byte(utils.BasicError(errorStr)))
	writer.Flush()
	log.Println(err)
}
