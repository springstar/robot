package server

import (
	"github.com/smallnest/gofsm"
)

type Robot struct {
	mgr *RobotManager
	account *Account
	fsm *fsm.StateMachine
	state string
}

func newRobot(account *Account, robotMgr *RobotManager, fsm *fsm.StateMachine) *Robot {
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

func (robot *Robot) startup() {

}

func (robot *Robot) initFsm(transitions []fsm.Transition) {
	delegate := &fsm.DefaultDelegate{P: &RobotEventProcessor{}}

	robot.fsm = fsm.NewStateMachine(delegate, transitions...)
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




