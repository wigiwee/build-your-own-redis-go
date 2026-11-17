package command

import (
	"errors"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleLrangeCommand(client *client.Client, args []string) {
	if len(args) != 4 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}

	start, err := strconv.Atoi(args[2])
	if err != nil {
		client.WriteResponse(err)
		return
	}
	stop, err := strconv.Atoi(args[3])
	if err != nil {
		client.WriteResponse(err)
		return
	}
	resArr := h.store.GetListRange(args[1], start, stop)

	client.WriteResponse(resArr)
	return
}
