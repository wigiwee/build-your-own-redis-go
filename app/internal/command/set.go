package command

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/client"
)

func (h *Handler) handleSetCommand(client *client.Client, args []string) {
	if len(args) < 3 {
		client.WriteResponse(errors.New("invalid command"))
		return
	}
	ttl := 0 * time.Millisecond
	if len(args) > 4 {
		expValue, err := strconv.Atoi(args[4])
		if err != nil {
			client.WriteResponse(err)
			return
		}
		if strings.ToUpper(args[3]) == "EX" {
			ttl = time.Duration(expValue) * time.Second
		} else if strings.ToUpper(args[3]) == "PX" {
			ttl = time.Duration(expValue) * time.Millisecond
		} else {
			client.WriteResponse(errors.New("invalid exp value " + args[4]))
		}
	}
	h.store.Set(args[1], args[2], ttl)
	client.WriteSimpleString("OK")
	return
}
