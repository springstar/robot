package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/config"
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

func (sm *SceneManager) getPath(sn int32) []*core.Vec2 {
	if path, ok := sm.pathpoints[sn]; ok {
		return path
	}

	return []*core.Vec2{}

}