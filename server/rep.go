package server

import (
	"github.com/springstar/robot/config"
)

func isRep(mapSn int) bool {
	confMap := config.FindConfMap(mapSn)
	if confMap.Type == "rep" {
		return true
	}

	return false
}