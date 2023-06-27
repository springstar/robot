package server

import (
	_ "fmt"
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/config"
)

type SceneEventType int32

const (
	SEVT_FLUSHMONSTER = 5
)


type SceneManager struct {
	pathpoints map[int32][]*core.Vec2
}

func newSceneManager() *SceneManager {
	return &SceneManager{
		pathpoints: make(map[int32][]*core.Vec2),
	}
}

func (sm *SceneManager) init() {
	confPathPoints := config.GetAllConfPathPoint()
	for _, conf := range confPathPoints {
		var points []*core.Vec2
		strPoints := []string{conf.Position1, conf.Position2, conf.Position3, conf.Position4, conf.Position5, conf.Position6, conf.Position7,
			conf.Position8, conf.Position9, conf.Position10 }	
	
		for _, s := range strPoints {
			pts := core.Str2Float32Slice(s)
			pt := core.NewVec2(pts[0], pts[1])
			points = append(points, pt)
		}

		sm.pathpoints[int32(conf.Sn)] = points
	}
}

func (sm *SceneManager) getPoint(sn int32, num int) *core.Vec2 {
	if num < 1 || num > 10 {
		return nil
	}

	path := sm.getPath(sn)
	if len(path) < 10 {
		return nil
	}

	p := path[num - 1]
	return p
}

func (sm *SceneManager) getPath(sn int32) []*core.Vec2 {
	if path, ok := sm.pathpoints[sn]; ok {
		return path
	}

	return []*core.Vec2{}

}