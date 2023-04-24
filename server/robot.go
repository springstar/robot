package server

import (
	_ "net"
	"fmt"
	"github.com/springstar/robot/core"
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
	r.conn.Connect(serv.cfg.Url)

	
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




