package server

import (
	"strconv"
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"

	"github.com/springstar/robot/pb"
	
)


type Character struct {
	humanId int64
	name string
	profession int32
	level int32
	sex int32
	combat int32
	equips []int32
	fashionSn int32
	roleSn int32
	soul int32
	curHp int64
	curMp int32
	curExp int64
	speed int32
	mapSn int32
	pos *core.Vec2
	dir *core.Vec2

}

func newCharacter() *Character {
	return &Character{
	}
}

func (c *Character) setName(name string) {
	c.name = name
}

func (c *Character) onCreate(humanId int64, fashionSn int32) {
	c.humanId = humanId
	c.fashionSn = fashionSn
}

func (c *Character) onLoad(d *pb.DCharacter) {
	c.humanId = d.GetId()
	c.profession = d.GetProfession()
	c.equips = d.GetEquip()
	c.level = d.GetLevel()
	c.combat = d.GetCombat()
	c.name = d.GetName()
	c.roleSn = d.GetRoleSn()
	c.soul = d.GetSoul()
	c.sex = d.GetSex()
}

func (c *Character) onInit(human *pb.DHuman, stage *pb.DInitDataStage) {
	c.curHp = human.HpCur
	c.curMp = human.MpCur
	c.curExp = human.ExpCur
	c.combat = human.Combat
	c.speed = human.Prop.Speed
	c.mapSn = stage.GetSn()

	// core.Info("speed is ", c.speed)

	c.pos = core.NewVec2(stage.PosNow.GetX(), stage.PosNow.GetY())
	c.dir = core.NewVec2(stage.DirNow.GetX(), stage.DirNow.GetY())
}

func (r *Robot) awakeSoul() {
	souls := []int{11, 12, 13}
	idx := core.GenRandomInt(len(souls))
	msg := msg.SerializeCSSoulAwaken(msg.MSG_CSSoulAwaken, strconv.Itoa(souls[idx]))
	r.sendPacket(msg)
}

func (r *Robot) handleSoulAwaken(packet *core.Packet) {
	awakeMsg := msg.ParseSCSoulAwaken(int32(msg.MSG_SCSoulAwaken), packet.Data)
	r.profession = awakeMsg.GetProfession()
	r.roleSn = awakeMsg.GetRoleSn()
	r.soul = awakeMsg.GetSoul()
	
	
	request := msg.SerializeCSInstanceLeave(msg.MSG_CSInstanceLeave)
	r.sendPacket(request)
}

func (r *Robot) sendUpgradeSoulRequest() {
	request := msg.SerializeCSIncSoulPointRequest(uint32(msg.MSG_CSIncSoulPointRequest))
	r.sendPacket(request)
}

func (r *Robot) sendGM(s string) {
	request := msg.SerializeCSInformChat(uint32(msg.MSG_CSInformChat), 1, s, r.humanId)
	r.sendPacket(request)
}