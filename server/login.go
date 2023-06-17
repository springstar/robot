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
	// fmt.Println("send login")
	packet := msg.SerializeCSLogin(msg.MSG_CSLogin, r.account.name, "123456", "", 1001, 1)
	r.sendPacket(packet)
}

func (r *Robot) querychars() {
	// fmt.Println("query chars")
	packet := msg.SerializeCSQueryCharacters(msg.MSG_CSQueryCharacters, 1)
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

func (r *Robot) handleLoginQueue(packet *core.Packet) {
	fmt.Println("login queue")
}

func (r *Robot) handleQueryCharacters(packet *core.Packet) {
	msg := msg.ParseSCQueryCharactersResult(int32(packet.Type), packet.Data)
	characters := msg.GetCharacters()
	if len(characters) == 0 {
		r.fsm.trigger(r.fsm.state, "create", r)
	} else {
		idx := core.GenRandomInt(len(characters))
		char := characters[idx]
		r.onLoad(char)
		r.fsm.trigger(r.fsm.state, "clogin", r)	
	}
}

func (r *Robot) sendCharacterLogin() {
	packet := msg.SerializeCSCharacterLogin(msg.MSG_CSCharacterLogin, r.humanId)
	r.sendPacket(packet)	
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

	name := serv.getNameManager().randomGenName(conf.RoleSex)
	r.setName(name)
	
	msg := msg.SerializeCSCharacterCreate(msg.MSG_CSCharacterCreate, name, int32(soulInt), false, avatars, 1)
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

	r.onCreate(msg.GetHumanId(), msg.GetFashionSn())
	r.fsm.trigger(r.fsm.state, "creatok", r)
	queueStat(STAT_CREATE_ROLES, 1)

}

func (r *Robot)handleCharacterLogin(packet *core.Packet) {
	msg := msg.ParseSCCharacterLoginResult(int32(packet.Type), packet.Data)
	if msg.GetResultCode() == 0 {
		r.fsm.trigger(r.fsm.state, "cloginok", r)
	} else {
		fmt.Println("character login failed")
		r.fsm.trigger(r.fsm.state, "cloginfail", r)
		queueStat(STAT_LOGIN_FAILS, int32(1))
	}
}

func (r *Robot)waitForInit() {
	// fmt.Println("waitForInit")
}

func (r *Robot)handleInitData(packet *core.Packet) {
	// fmt.Println("recv init data")
	msg := msg.ParseSCInitData(int32(packet.Type), packet.Data)
	r.onInit(msg.Human, msg.Stage)
	executor := r.findExecutor("quest").(*RobotQuestExecutor)
	executor.initQuest(msg.QuestInfo)
	r.fsm.trigger(r.fsm.state, "initok", r)
}

func (r *Robot)handleKick(packet *core.Packet) {
	msg := msg.ParseSCHumanKick(int32(packet.Type), packet.Data)
	fmt.Println(msg.GetReason())
}
