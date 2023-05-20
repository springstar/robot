package server

import (
	_ "fmt"
	"github.com/springstar/robot/config"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFindConf(t *testing.T) {
	mgr := newJsonConfigManager()
	mgr.init("../config")

	assert.Greater(t, len(mgr.confs), 0)
	confAIDatas := config.GetAllConfAI()
	assert.Greater(t, len(confAIDatas), 0)
	conf := config.FindConfAI(102008)
	if conf == nil {
		t.Error("conf nil")
		return
	}
	assert.Equal(t, conf.Sn, 102008)
	
}