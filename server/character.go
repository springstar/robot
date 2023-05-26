package server

import (
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

func (c *Character) onInit(d *pb.DHuman) {
	c.curHp = d.HpCur
	c.curMp = d.MpCur
	c.curExp = d.ExpCur
	c.combat = d.Combat
}