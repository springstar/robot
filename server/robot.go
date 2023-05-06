package server

import (
	_ "net"
	"fmt"
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/msg"
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
	case "on_connection_established":
		r.on_connection_established()
	default:
		fmt.Println(action)	
	}
}

func (r *Robot) connect() {
	r.conn = core.NewWsConnection()
	err := r.conn.Connect(r.mgr.url)
	if err != nil {
		fmt.Print(err)
		r.fsm.trigger("connecting", "cfail", r)
	}

	r.fsm.trigger("connecting", "cok", r)
	
}

func (r *Robot) on_connection_established() {
	msg := msg.SerializeCSLogin(111, "robot", "123456", "", 1001, 1)
	r.conn.Write(msg)
}

func (r *Robot) loop() {

}

type RobotManager struct {
	robots map[int]*Robot
	url string
}

func newRobotManager(url string) *RobotManager {
	return &RobotManager{
		robots : make(map[int]*Robot),
		url : url,
	}
}

func (mgr *RobotManager)add(account int, robot *Robot) {
	mgr.robots[account] = robot
}




