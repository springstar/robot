package server

import (
	"fmt"
	_ "github.com/springstar/robot/config"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFindConf(t *testing.T) {
	mgr := NewJsonConfigManager()
	mgr.loadConfFile("../config")
	assert.Greater(t, len(mgr.confs), 0)
	conf := mgr.findConf("ConfAI", 301301)
	fmt.Println(conf)
	// assert.Equal(t, conf.Sn, 301301)
	
}