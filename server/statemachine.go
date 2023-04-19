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