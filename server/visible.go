package server

const (
	MAX_VISIBLE_OBJS = 20
)
type VisibleRange struct {
	visibleObjs map[int][]*WorldObj
}

func newVisibleRange() *VisibleRange {
	return &VisibleRange{
		visibleObjs: make(map[int][]*WorldObj),
	}
}

func (v *VisibleRange) addVisibleObj(typ int, obj *WorldObj) {

}

func (v *VisibleRange) removeVisibleObj(typ int, objId int64) {

}

