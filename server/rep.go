package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/config"
)

func isRep(mapSn int) bool {
	confMap := config.FindConfMap(mapSn)
	if confMap.Type == "rep" {
		return true
	}

	return false
}

func (r *Robot) sendEnterInstance(repSn int, questSn int) {
	request := msg.SerializeCSInstanceEnter(uint32(msg.MSG_CSInstanceEnter), int32(repSn), int32(questSn), false)
	r.sendPacket(request)
}