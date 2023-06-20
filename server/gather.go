package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"
)

func (r *Robot) gatherFirst(id int64) {
	request := msg.SerializeCSGatherFirst(uint32(msg.MSG_CSGatherFirst), id, 0)
	r.sendPacket(request)
}

func (r *Robot) gatherSecond(id int64) {
	request := msg.SerializeCSGatherSecond(uint32(msg.MSG_CSGatherSecond), id)
	r.sendPacket(request)
}

func (r *Robot) HandleGatherFirst(packet *core.Packet) {

}

func (r *Robot) HandleGatherSecond(packet *core.Packet) {

}