package server

import (
	"fmt"
)

type Robot struct {
	account *Account
}

func newRobot() *Robot {
	return &Robot{}
}

func (robot *Robot) startup() {
	robot.account = serv.accountMgr.alloc()
	fmt.Println(robot.account)
}