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

func (v *VisibleRange) addVisibleObj(obj *WorldObj) {
	// if v.pos.DistanceToSquared(obj.pos) > MAX_VISIBLE_DISTANCE * MAX_VISIBLE_DISTANCE {
	// 	return
	// }

	v.visibleObjs[obj.id] = obj
}

func (v *VisibleRange) removeVisibleObj(id int64) {
	delete(v.visibleObjs, id)
}
