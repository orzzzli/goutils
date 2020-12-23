package server

import (
	"bufio"
	"errors"
	"net"
	"sync/atomic"
)

type CustomServer struct {
	Server
	handler  func(reader *bufio.Reader, connIndex int)
	connList map[int]net.Conn
	index    int32
}

func newCustomServer(name string, address string) *CustomServer {
	t := &CustomServer{}
	t.name = name
	t.address = address
	t.serviceMap = make(map[string]interface{}, 5)
	t.connList = make(map[int]net.Conn, 100)
	return t
}

func (c *CustomServer) SetHandler(reader func(reader *bufio.Reader, connIndex int)) {
	c.handler = reader
}

func (c *CustomServer) Serve() error {
	l, err := net.Listen("tcp", c.address)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		res := atomic.AddInt32(&c.index, 1)
		c.connList[int(res)] = conn
		go c.connHandler(conn, int(res))
	}
}

func (c *CustomServer) RegisterService(url string, handler ...interface{}) error {
	_, exist := c.serviceMap[url]
	if exist {
		return errors.New("service already register")
	}
	if len(handler) == 0 {
		return errors.New("handler must send")
	}
	t, ok := handler[0].(func(p interface{}) string)
	if !ok {
		return errors.New("handler type must be [ func(params interface{}) string ]")
	}
	c.serviceMap[url] = t
	return nil
}

func (c *CustomServer) CallService(url string, params interface{}) (string, error) {
	handler, exist := c.serviceMap[url]
	if !exist {
		return "", errors.New("service not exist")
	}
	h, _ := handler.(func(p interface{}) string)
	res := h(params)
	return res, nil
}

func (c *CustomServer) SendTo(index int, message string) {
	conn, ok := c.connList[index]
	if ok {
		conn.Write([]byte(message))
	}
}

func (c *CustomServer) BroadcastTo(message string, index ...int) {
	if len(index) == 0 {
		for _, conn := range c.connList {
			conn.Write([]byte(message))
		}
	} else {
		for _, i := range index {
			c.SendTo(i, message)
		}
	}
}

func (c *CustomServer) connHandler(conn net.Conn, index int) {
	bufReader := bufio.NewReader(conn)
	c.handler(bufReader, index)
}
