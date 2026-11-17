package command

import (
	"errors"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleLPopCommand(client *client.Client, args []string) {
	if len(args) > 3 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	n := 1
	if len(args) != 2 {
		var err error
		n, err = strconv.Atoi(args[2])
		if err != nil {
			client.WriteResponse(err)
			return
		}
	}
	pop := h.store.LPop(args[1], n)
	if len(pop) == 0 {
		client.WriteResponse(nil)
	} else if len(pop) == 1 {
		client.WriteResponse(pop[0])
	} else {
		client.WriteResponse(pop)
	}
	return
}
