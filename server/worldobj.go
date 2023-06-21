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
	
	typ WorldObjType
	pos *core.Vec2
}

func newWorldObj(id int64) *WorldObj {
	return &WorldObj{
		id: id,
	}
}