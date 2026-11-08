package command

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleSetCommand(writer *bufio.Writer, args []string) error {
	if len(args) < 3 {
		return errors.New("invalid command")
	}
	ttl := 0 * time.Millisecond
	if len(args) > 3 {
		expValue, err := strconv.Atoi(args[4])
		if err != nil {
			return err
		}
		if strings.ToUpper(args[3]) == "EX" {
			ttl = time.Duration(expValue) * time.Second
		} else if strings.ToUpper(args[3]) == "PX" {
			ttl = time.Duration(expValue) * time.Millisecond
		} else {
			return errors.New("invalid exp value " + args[4])
		}
	}
	h.store.Set(args[1], args[2], ttl)
	writer.Write([]byte(utils.OK_BULK_STR))
	writer.Flush()
	return nil
}
