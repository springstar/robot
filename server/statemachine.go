package server

import (
	"github.com/smallnest/gofsm"
)

var (

	

)

func defaultTransitions() []fsm.Transition {
	transitions := []fsm.Transition{
		fsm.Transition{From: "waitForConnect", Event: "connect", To: "connected", Action: "connect"},
		fsm.Transition{From: "connected", Event: "done", To: "finished", Action: "connect"},
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
	return nil
}

func (p *RobotEventProcessor) OnEnter(toState string, args []interface{}) {
	r := args[0].(*Robot)
	r.state = toState
}

func (p *RobotEventProcessor) OnActionFailure(action string, fromState string, toState string, args []interface{}, err error) {
	
}