package modules

import (
	"github.com/springstar/robot/msg"
	"fmt"
	"github.com/springstar/robot/core"
)

type RobotLoginSession struct {
	core.IDispatcher
}

func NewLoginSession(dispatcher core.IDispatcher) *RobotLoginSession {
	return &RobotLoginSession{
		dispatcher,
	}
}

func (s *RobotLoginSession) Init() {
	s.Register(112, s)
}

func (s *RobotLoginSession) Fini() {

}

func (s *RobotLoginSession) HandleMessage(packet *core.Packet) {
	switch packet.Type {
	case 112:
		s.handleLoginResult(packet)
	}
}

func (s *RobotLoginSession) handleLoginResult(packet *core.Packet) {
	msg := msg.ParseSCLoginResult(int32(packet.Type), packet.Data)
	fmt.Println(msg.GetResultCode())
	fmt.Println(msg.GetResultReason())
}