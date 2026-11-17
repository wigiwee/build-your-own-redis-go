package command

import (
	"errors"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleBLPopCommand(client *client.Client, args []string) {
	if len(args) != 3 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	n, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		client.WriteResponse(err)
		return
	}
	ttl := time.Duration(n * float64(time.Second))
	output := h.store.BLPop(args[1], ttl)
	client.WriteResponse(output)
	return
}
