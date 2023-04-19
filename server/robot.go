package server

import (
	"github.com/smallnest/gofsm"
)

type Robot struct {
	account *Account
	fsm *fsm.StateMachine
	state string
}

func newRobot() *Robot {
	return &Robot{}
}

func (robot *Robot) startup(accountMgr *AccountManager, robotMgr *RobotManager, fsm *fsm.StateMachine) {
	robot.account = accountMgr.alloc()
	if (robot.account != nil) {
		robotMgr.add(robot.account.id, robot)		
	}

	robot.fsm = fsm

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




