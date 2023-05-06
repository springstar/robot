package core

import (

	"fmt"
	"time"
	"context"
	"nhooyr.io/websocket"
)

type WebSockConn struct {
	c *websocket.Conn
}

func NewWsConnection() *WebSockConn {
	return &WebSockConn{
	}
}

func (wsc *WebSockConn) Connect(addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	var err error
	wsc.c, _, err = websocket.Dial(ctx, addr, nil)
	if err != nil {
		fmt.Print("connect err ", err)
	}

	return nil
}

func (wsc *WebSockConn) Write(p []byte) (n int, err error) {
	n = len(p)
	err = wsc.c.Write(context.Background(), websocket.MessageBinary, p)
	if (err != nil) {
		fmt.Println(err)
		return 0, err
	}
	return n, nil
}
