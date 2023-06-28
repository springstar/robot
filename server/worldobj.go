package server

import (
	"github.com/springstar/robot/core"
)

type WorldObjType int32

const (
	WOT_HUMAN = iota		
	WOT_MONSTER  
	WOT_DROP  
	WOT_DOT  
	WOT_BULLET 
	WOT_NPC 
	WOT_PICK 
	WOT_TRIGGER
	WOT_COMMON
	WOT_PARTNER

)

type WorldObj struct {
	id int64
	sn int
	typ WorldObjType
	pos *core.Vec2
}

func newWorldObj(id int64, typ WorldObjType, pos *core.Vec2) *WorldObj {
	return &WorldObj{
		id: id,
		typ: typ,
		pos: pos,
	}
}

func (obj *WorldObj) getId() int64 {
	return obj.id
}

func (obj *WorldObj) getType() WorldObjType {
	return obj.typ
}

func (obj *WorldObj) getPos() *core.Vec2 {
	return obj.pos
}

type GatherObj struct {
	*WorldObj
}

type MonsterObj struct {
	*WorldObj
	curHp int64
	maxHp int64
}

func newGatherObj(id int64, typ WorldObjType, pos *core.Vec2) *GatherObj {
	return &GatherObj{
		WorldObj: newWorldObj(id, typ, pos),
	}
}

func newMonsterObj(id int64, typ WorldObjType, pos *core.Vec2) *MonsterObj {
	return &MonsterObj{
		WorldObj: newWorldObj(id, typ, pos),

	}
}

func (obj *MonsterObj) isDead() bool {
	return obj.curHp == 0
}
