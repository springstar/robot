package server

import (
	_ "net"
	"fmt"
	"context"
	"time"
	"github.com/springstar/robot/core"
	"nhooyr.io/websocket"
	_ "github.com/gobwas/ws"
)

type Robot struct {
	conn core.NetConnection
	mgr *RobotManager
	account *Account
	fsm *RobotFsm
}

func newRobot(account *Account, robotMgr *RobotManager, fsm *RobotFsm) *Robot {
	r := &Robot{
		mgr : robotMgr,
		account : account,
		fsm : fsm,
	}

	if (r.account != nil) {
		robotMgr.add(r.account.id, r)		
	}

	return r
}

func (r *Robot) startup() {
	r.fsm.trigger("entry", "connect", r)
}

func (r *Robot) doAction(action string) {
	switch action {
	case "connect":
		r.connect()
	default:
		fmt.Println(action)	
	}
}

func (r *Robot) connect() {
	r.conn = core.NewWsConnection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	
	c, _, err := websocket.Dial(ctx, serv.cfg.Url, nil)
	if err != nil {
		fmt.Print("connect err ", err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	
}

type RobotManager struct {
	robots map[int]*Robot
}

func newRobotManager() *RobotManager {
	return &RobotManager{
		robots : make(map[int]*Robot),
	}
}

func (mgr *RobotManager)add(account int, robot *Robot) {
	mgr.robots[account] = robot
}




