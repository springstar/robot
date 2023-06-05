package server

import (
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
	c.mapSn = stage.GetSn()
	c.pos = core.NewVec2(stage.PosNow.GetX(), stage.PosNow.GetY())
	c.dir = core.NewVec2(stage.DirNow.GetX(), stage.DirNow.GetY())
}