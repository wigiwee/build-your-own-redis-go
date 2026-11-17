package client

import (
	"bufio"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

type Client struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer

	InTransaction  bool
	QueuedCommands [][]string
	WatchedKey     string
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:           conn,
		Reader:         bufio.NewReader(conn),
		Writer:         bufio.NewWriter(conn),
		InTransaction:  false,
		QueuedCommands: make([][]string, 0),
	}
}

func (c *Client) ReadCommand() ([]string, error) {
	return utils.DeserializeArray(c.Reader)

}

func (c *Client) WriteResponse(value any) {
	respStr, err := utils.Serialize(value)
	if err != nil {
		log.Println("error writing to the client", err)
	}
	_, err = c.Writer.Write([]byte(respStr))
	if err != nil {
		log.Println("error writing to the client", err)
	}
	err = c.Writer.Flush()
	if err != nil {
		log.Println("error writing to the client", err)
	}
}

func (c *Client) WriteSimpleString(str string) {
	_, err := c.Writer.Write([]byte(utils.SerializeSimpleString(str)))
	if err != nil {
		log.Println("error writing to the client ", err)
	}
	err = c.Writer.Flush()
	if err != nil {
		log.Println("error writing to the client ", err)
	}
}

func (c *Client) SetWriter(writer *bufio.Writer) {
	c.Writer = writer
}

func (c *Client) WriteRaw(str string) {
	_, err := c.Writer.Write([]byte(str))
	if err != nil {
		log.Println("error writing to the client ", err)
	}
	err = c.Writer.Flush()
	if err != nil {
		log.Println("error writing to the client ", err)
	}
}
