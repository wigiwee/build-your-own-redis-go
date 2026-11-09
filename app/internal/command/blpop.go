package command

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (h *Handler) handleBLPopCommand(writer *bufio.Writer, args []string) error {
	if len(args) != 3 {
		return errors.New("invalid command")
	}
	n, err := strconv.ParseFloat(args[2], 64)
	fmt.Println(n)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ttl := time.Duration(n * float64(time.Second))
	fmt.Println("calling func")
	output := h.store.BLPop(args[1], ttl)
	fmt.Println("got res ", output)
	respArr, err := utils.SerializeArray(output)
	if err != nil {
		return err
	}
	writer.Write([]byte(respArr))
	writer.Flush()
	return nil
}
