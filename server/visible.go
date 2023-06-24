package server

import (
	_ "github.com/springstar/robot/core"
)

const (
	MAX_VISIBLE_OBJS = 20
	MAX_VISIBLE_DISTANCE = 10
)

type VisibleRange struct {
	*Robot
	visibleObjs map[int64]*WorldObj
}

func newVisibleRange(r *Robot) *VisibleRange {
	return &VisibleRange{
		Robot: r,
		visibleObjs: make(map[int64]*WorldObj),
	}
}

func (v *VisibleRange) addObj(obj *WorldObj) {
	// if v.pos.DistanceToSquared(obj.pos) > MAX_VISIBLE_DISTANCE * MAX_VISIBLE_DISTANCE {
	// 	return
	// }

	v.visibleObjs[obj.id] = obj
}

func (v *VisibleRange) removeObj(id int64) {
	delete(v.visibleObjs, id)
}

func (v *VisibleRange) findObj(id int64) *WorldObj {
	if obj, ok := v.visibleObjs[id]; ok {
		return obj
	}

	return nil
}

func (v *VisibleRange) findGatherObj(sn int) *WorldObj {
	for _, obj := range v.visibleObjs {
		if obj.typ != WOT_PICK {
			continue
		}

		if obj.sn == sn {
			return obj
		}
	}

	return nil
}

func (v *VisibleRange) findMonsterObj(sn int) *WorldObj {
	for _, obj := range v.visibleObjs {
		if obj.typ != WOT_MONSTER {
			continue
		}

		if obj.sn == sn {
			return obj
		}
	}

	return nil
}

