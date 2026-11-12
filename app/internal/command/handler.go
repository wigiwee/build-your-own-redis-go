package command

import (
	"bufio"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
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
		WriteError(writer, errors.New("no args received"))
	}

	switch strings.ToUpper(args[0]) {

	case "PING":
		err := h.handlePingCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "ECHO":
		err := h.handleEchoCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "GET":
		err := h.handleGetCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "SET":
		err := h.handleSetCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "RPUSH":
		err := h.handleRpushCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "LPUSH":
		err := h.handleLpushCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "LRANGE":
		err := h.handleLrangeCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "LLEN":
		err := h.handleLlenCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "LPOP":
		err := h.handleLPopCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "BLPOP":
		err := h.handleBLPopCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "TYPE":
		err := h.handleTypeCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "XADD":
		err := h.handleXAddCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "XRANGE":
		err := h.handleXRangeCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}
	case "XREAD":
		err := h.handleXReadCommand(writer, args)
		if err != nil {
			WriteError(writer, err)
		}

	default:
		WriteError(writer, errors.New("unsupported command"))

	}

}

func WriteError(writer *bufio.Writer, err error) {
	serErr, err := utils.Serialize(err)
	_, err = writer.Write([]byte(serErr))
	if err != nil {
		log.Println("error writing err", err)
	}

	err = writer.Flush()
	if err != nil {
		log.Println("error writing err", err)
	}
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
	GetType(key string) string
	AddStream(streamkey string, id string, args []string) (string, error)
	RangeStream(streamKey string, start string, stop string) ([]storage.StreamEntry, error)
	ReadStream(streamKey string, streamId string) ([]any, error)
}
