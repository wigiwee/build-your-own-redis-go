package command

import (
	"errors"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) HandleCommand(client *client.Client, args []string) {

	if len(args) == 0 {
		client.WriteResponse(errors.New("no args received"))
	}

	switch strings.ToUpper(args[0]) {

	case "MULTI":
		h.handleMultiCommand(client, args)
		return

	case "EXEC":
		h.handleExecCommand(client, args)
		return

	case "DISCARD":
		h.handleDiscardCommand(client, args)
		return

	case "WATCH":
		h.handleWatchCommand(client, args)
		return
	}

	if client.InTransaction {
		client.QueuedCommands = append(client.QueuedCommands, args)
		client.WriteSimpleString("QUEUED")
		return
	}
	h.executeCommand(client, args)

}

func (h *Handler) executeCommand(client *client.Client, args []string) {

	switch strings.ToUpper(args[0]) {
	case "PING":
		h.handlePingCommand(client, args)

	case "ECHO":
		h.handleEchoCommand(client, args)

	case "GET":
		h.handleGetCommand(client, args)

	case "SET":
		h.handleSetCommand(client, args)

	case "RPUSH":
		h.handleRpushCommand(client, args)

	case "LPUSH":
		h.handleLpushCommand(client, args)

	case "LRANGE":
		h.handleLrangeCommand(client, args)

	case "LLEN":
		h.handleLlenCommand(client, args)

	case "LPOP":
		h.handleLPopCommand(client, args)

	case "BLPOP":
		h.handleBLPopCommand(client, args)

	case "TYPE":
		h.handleTypeCommand(client, args)

	case "XADD":
		h.handleXAddCommand(client, args)

	case "XRANGE":
		h.handleXRangeCommand(client, args)

	case "XREAD":
		h.handleXReadCommand(client, args)

	case "INCR":
		h.handleIncrCommand(client, args)

	default:
		client.WriteResponse(errors.New("unsupported command"))
	}
}

type Store interface {
	Set(key string, value string, ttl time.Duration)
	Get(key string) any
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
	ReadStream(streamKey string, streamId string) (storage.StreamReadResult, error)
	ReadStreamBlock(streamKey string, streamId string, ttl time.Duration) (storage.StreamReadResult, error)
	Incr(key string) (uint64, error)
	GetVersion(key string) uint64
}
