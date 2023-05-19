package server

import (
	"fmt"
	"strings"
	"github.com/springstar/robot/config"
	"github.com/springstar/robot/msg"
	_"github.com/springstar/robot/pb"
	"github.com/springstar/robot/core"
)

func (r *Robot) sendLoginRequest() {
	fmt.Println("send login")
	packet := msg.SerializeCSLogin(111, r.account.name, "123456", "", 1001, 1)
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
		r.fsm.trigger(r.fsm.state, "create", r)
	} else {

	}
}

func (r *Robot) createChar() {
	fmt.Println("create char ", r.account.name)
	conf := randomRole()
	souls := strings.Split(conf.MatchingSoul, ",")
	soul := souls[core.GenRandomInt(len(souls))]
	soulInt, err := core.Str2Int(soul)
	if (err != nil) {
		fmt.Println(err)
		return
	}

	confModel := config.FindConfCharacterModel(conf.ModelSn)
	if confModel == nil {
		return
	}

	avatars, err := core.Str2Int32Slice(confModel.PresetAvatar)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	msg := msg.SerializeCSCharacterCreate(1005, r.account.name, int32(soulInt), false, avatars, 1)
	r.sendPacket(msg)
}

func randomRole() *config.ConfCharacterRole {
	datas := config.GetAllConfCharacterRole()
	n := len(datas)
	rnd := core.GenRandomInt(n)
	conf := datas[rnd]
	return conf	
	
}

func (r *Robot)handleCreateResult(packet *core.Packet) {
	msg := msg.ParseSCCharacterCreateResult(int32(packet.Type), packet.Data)
	if msg.GetResultCode() == -1 {
		fmt.Printf("%s create character failed, %s", r.account.name, msg.GetResultReason())
		return
	}

	fmt.Println("create ok")
}
