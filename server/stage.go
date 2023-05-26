package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/msg"
	"fmt"
)

func (r *Robot) enterStage() {
	fmt.Println("enterStage")
	packet := msg.SerializeCSStageEnter(msg.MSG_CSStageEnter)	
	r.sendPacket(packet)
}

func (r *Robot) handleEnterStage(packet *core.Packet) {
	fmt.Println("recv ")
	msg := msg.ParseSCStageEnterResult(int32(packet.Type), packet.Data)
	stageObjs := msg.GetObj()
	for _, obj := range stageObjs {
		fmt.Println(obj.GetPos())
		fmt.Println(obj.GetName())
	}
}