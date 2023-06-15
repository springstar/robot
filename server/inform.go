package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"
)

func (r *Robot) handleInform(packet *core.Packet) {
	inform := msg.ParseSCInformMsg(int32(msg.MSG_SCInformMsg), packet.Data)
	core.Info(inform.GetContent())
}