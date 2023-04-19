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

type RobotEventProcessor struct {

}

func (p *RobotEventProcessor) OnExit(fromState string, args []interface{}) {
	
}

func (p *RobotEventProcessor) Action(action string, fromState string, toState string, args []interface{}) error {
	return nil
}

func (p *RobotEventProcessor) OnEnter(toState string, args []interface{}) {
	r := args[0].(*Robot)
	r.state = toState
}

func (p *RobotEventProcessor) OnActionFailure(action string, fromState string, toState string, args []interface{}, err error) {
	
}


