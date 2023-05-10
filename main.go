package main

import (
	"github.com/springstar/robot/server"
	"strings"
	"fmt"
	"flag"
)

var (
	tag *string
)

func init() {
	tag = flag.String("tag", "server", "start as server")

}

func initMsgDesciptors() {
	g := newDescriptorGen()
	g.parse("msg/protocol")
}


func startServer() {
	initMsgDesciptors()
	serv := server.NewServer()
	serv.Init()
	serv.Run()
}

func startClient() {
	fmt.Println("start client")
}

func main() {
	flag.Parse()
	if strings.Compare(*tag, "server") == 0 {
		startServer()
	} else if strings.Compare(*tag, "client") == 0 {
		startClient()
	}
}