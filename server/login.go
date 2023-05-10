package server

import (
	"fmt"
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"
)

func (r *Robot) sendLoginRequest() {
	fmt.Println("send login")
	packet := msg.SerializeCSLogin(111, "robot", "123456", "", 1001, 1)
	r.sendPacket(packet)
}

func (r *Robot) querychars() {
	fmt.Println("query chars")
	packet := msg.SerializeCSQueryCharacters(1003, 1)
	r.sendPacket(packet)
}

func (r *Robot) handleLoginResult(packet *core.Packet) {
	msg := msg.ParseSCLoginResult(int32(packet.Type), packet.Data)
	code := msg.GetResultCode()
	if code == -1 {
		r.fsm.trigger("trylogin", "lfail", r)
		return
	}

	r.fsm.trigger("trylogin", "lok", r)

}

func (r *Robot) handleQueryCharacters(packet *core.Packet) {
	msg := msg.ParseSCQueryCharactersResult(int32(packet.Type), packet.Data)
	characters := msg.GetCharacters()
	if len(characters) == 0 {
		fmt.Println("create character")
	} else {

	}
}