package core

import (
	"fmt"
	"time"
	"context"
	"nhooyr.io/websocket"
)

type WebSockConn struct {

}

func NewWsConnection() *WebSockConn {
	return &WebSockConn{

	}
}

func (conn *WebSockConn) Connect(addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	
	c, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		fmt.Print("connect err ", err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	return nil
}

