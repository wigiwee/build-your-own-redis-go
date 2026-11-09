package command

import (
	"bufio"
	"log"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) HandleCommand(writer *bufio.Writer, args []string) {

	if len(args) == 0 {
		log.Fatal("empty args received")
	}

	switch strings.ToUpper(args[0]) {

	case "PING":
		err := h.handlePingCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "ECHO":
		err := h.handleEchoCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "GET":
		err := h.handleGetCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "SET":
		err := h.handleSetCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "RPUSH":
		err := h.handleRpushCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "LPUSH":
		err := h.handleLpushCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "LRANGE":
		err := h.handleLrangeCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "LLEN":
		err := h.handleLlenCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "LPOP":
		err := h.handleLPopCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}
	case "BLPOP":
		err := h.handleBLPopCommand(writer, args)
		if err != nil {
			WriteError(writer, err.Error())
		}

	default:
		WriteError(writer, "Unknown Command")

	}

}

func WriteError(writer *bufio.Writer, errorStr string) {
	_, err := writer.Write([]byte(utils.BasicError(errorStr)))
	writer.Flush()
	log.Println(err)
}

type Store interface {
	Set(key, value string, ttl time.Duration)
	Get(key string) (string, bool)
	Delete(key string) bool
	RPushList(listName string, values []string) int
	LPushList(listName string, values []string) int
	GetListLen(listName string) int
	LPop(listName string, n int) []string
	BLPop(listName string, ttl time.Duration) []string
	GetListRange(listName string, start, stop int) []string
}
