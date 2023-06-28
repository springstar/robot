package server

import (
	"github.com/springstar/robot/core"
)

const (
	MAX_VISIBLE_OBJS = 20
	MAX_VISIBLE_DISTANCE = 10
)

type iVisible interface {
	getPos() *core.Vec2
	getId() int64
	getType() WorldObjType
}

type VisibleRange struct {
	*Robot
	visibleObjs map[int64]iVisible
}

func newVisibleRange(r *Robot) *VisibleRange {
	return &VisibleRange{
		Robot: r,
		visibleObjs: make(map[int64]iVisible),
	}
}

func (v *VisibleRange) addObj(obj iVisible) {
	// if v.pos.DistanceToSquared(obj.pos) > MAX_VISIBLE_DISTANCE * MAX_VISIBLE_DISTANCE {
	// 	return
	// }

	v.visibleObjs[obj.getId()] = obj
}

func (v *VisibleRange) removeObj(id int64) {
	delete(v.visibleObjs, id)
}

func (v *VisibleRange) findObj(id int64) iVisible {
	if obj, ok := v.visibleObjs[id]; ok {
		return obj
	}

	return nil
}

func (v *VisibleRange) findGatherObj(sn int) *GatherObj {
	for _, obj := range v.visibleObjs {
		if obj.getType() != WOT_PICK {
			continue
		}
		
		gatherObj := (obj).(*GatherObj)
		if gatherObj.sn == sn {
			return gatherObj
		}
	}

	return nil
}

func (v *VisibleRange) findMonsterObj(sn int) *MonsterObj {
	for _, obj := range v.visibleObjs {
		if obj.getType() != WOT_MONSTER {
			continue
		}

		monsterObj := (obj).(*MonsterObj)
		if monsterObj.sn == sn {
			return monsterObj
		}
	}

	return nil
}

func createGather(id int64, typ WorldObjType, pos *core.Vec2, sn int) iVisible {
	gatherObj := newGatherObj(id, typ, pos)
	gatherObj.sn = sn
	return gatherObj
}

func createMonster(id int64, typ WorldObjType, pos *core.Vec2, sn int, curHp, maxHp int64) iVisible {
	monsterObj := newMonsterObj(id, typ, pos)
	monsterObj.sn = sn
	monsterObj.curHp = curHp
	monsterObj.maxHp = maxHp
	return monsterObj
}