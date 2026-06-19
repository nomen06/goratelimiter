package redis

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	mu     sync.Mutex
}

func (c *Client) Do(args []string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := len(args)
	input := "*"
	input += strconv.Itoa(n)
	input += "\r"
	input += "\n"
	for comp := range args {
		len2 := len(args[comp])
		input += "$"
		input += strconv.Itoa(len2)
		input += "\r"
		input += "\n"
		input += args[comp]
		input += "\r"
		input += "\n"
	}
	fmt.Fprint(c.conn, input)
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return response, nil
}
func Newclient(addr string) (*Client, error) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	connReader := bufio.NewReader(con)
	client := Client{conn: con, reader: connReader}

	return &client, nil
}
