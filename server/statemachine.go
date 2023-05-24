package server

import (
	"github.com/smallnest/gofsm"
	_ "fmt"
)

var (

	

)

type RobotFsm struct {
	fsm *fsm.StateMachine

	state string
}

func newFsm() *RobotFsm {
	return &RobotFsm{
		fsm: newDefaultStateMachine(defaultTransitions()),
	}
}

func (fsm *RobotFsm) trigger(currentState string, event string, args interface{}) {
	fsm.fsm.Trigger(currentState, event, args)
}

func defaultTransitions() []fsm.Transition {
	transitions := []fsm.Transition{
		fsm.Transition{From: "entry", Event: "connect", To: "connecting", Action: "connect"},
		fsm.Transition{From: "connecting", Event: "cok", To: "connected", Action: "on_connection_established"},
		fsm.Transition{From: "connecting", Event: "cfail", To: "retry", Action: "retry"},
		fsm.Transition{From: "connected", Event: "login", To: "trylogin", Action: "login"},
		fsm.Transition{From: "trylogin", Event: "lok", To: "querychars", Action: "querychars"},
		fsm.Transition{From: "trylogin", Event: "lfail", To: "disconnected", Action: "disconnected"},
		fsm.Transition{From: "querychars", Event: "create", To: "createchar", Action: "createchar"},
		fsm.Transition{From: "querychars", Event: "clogin", To: "clogin", Action: "sendCharacterLogin"},
		fsm.Transition{From: "createchar", Event: "creatok", To: "creatok", Action: "sendCharacterLogin"},
		fsm.Transition{From: "clogin", Event: "cloginfail", To: "disconnected", Action: "disconnected"},
		fsm.Transition{From: "creatok", Event: "cloginfail", To: "disconnected", Action: "disconnected"},
		fsm.Transition{From: "clogin", Event: "cloginok", To: "waitinit", Action: "waitInitData"},
		fsm.Transition{From: "creatok", Event: "cloginok", To: "waitinit", Action: "waitInitData"},
	}

	return transitions
}
 
func newDefaultStateMachine(transitions []fsm.Transition) *fsm.StateMachine {
	delegate := &fsm.DefaultDelegate{P: &RobotEventProcessor{}}
	fsm := fsm.NewStateMachine(delegate, transitions...)
	return fsm

}

type RobotEventProcessor struct {

}

func (p *RobotEventProcessor) OnExit(fromState string, args []interface{}) {
	
}

func (p *RobotEventProcessor) Action(action string, fromState string, toState string, args []interface{}) error {
	r := args[0].(*Robot)
	r.doAction(action)
	return nil
}

func (p *RobotEventProcessor) OnEnter(toState string, args []interface{}) {
	r := args[0].(*Robot)
	r.fsm.state = toState
}

func (p *RobotEventProcessor) OnActionFailure(action string, fromState string, toState string, args []interface{}, err error) {
	
}