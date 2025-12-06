package mister

import (
	"bufio"
	"net"
	"sync"
	"time"
)

type Client struct {
	addr             *net.TCPAddr
	name             string
	readWriteTimeout time.Duration
	dialer           *net.Dialer
	lock             sync.Mutex
	isConnected      bool
	isClosed         bool
	c                *net.TCPConn
	r                *bufio.Reader
}

func NewClient(addr *net.TCPAddr, name string, timeout time.Duration) (c *Client) {
	c = &Client{
		addr:             addr,
		name:             name,
		readWriteTimeout: timeout,
		dialer:           &net.Dialer{Timeout: timeout},
	}

	return
}

func (c *Client) Connect() (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.isClosed = false

	var conn net.Conn
	netAddr := net.Addr(c.addr)
	conn, err = c.dialer.Dial("tcp", netAddr.String())
	if err != nil {
		c.isConnected = false
		return
	}
	c.c = conn.(*net.TCPConn)

	c.r = bufio.NewReaderSize(c.c, 4096)
	c.isConnected = true
	return
}

func (c *Client) Close() (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.isClosed = true
	c.isConnected = false
	err = c.c.Close()
	return
}

func (c *Client) IsClosed() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.isClosed
}
