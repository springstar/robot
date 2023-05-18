package server

import (
	_ "fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"github.com/springstar/robot/config"
)

type JsonConfig struct {
	name string
	// struct slice, every map is a struct, like ConfAI, ConItem etc
	data []map[string]interface{}
}

func NewJsonConfig(name string) *JsonConfig {
	return &JsonConfig{
		name : name,
	}
}

func (c *JsonConfig) loadData(bytes []byte) error {
	if err := json.Unmarshal(bytes, &c.data); err != nil {
		panic(err)
	}

	for _, v :=range c.data {
		config.LoadConf(c.name, v)
	}
	return nil
}

type JsonConfigManager struct {
	confs map[string]*JsonConfig
}

func newJsonConfigManager() *JsonConfigManager {
	return &JsonConfigManager{
		confs: make(map[string]*JsonConfig),
	}
}

func (m *JsonConfigManager) init(path string) {
	config.InitLoaders()
	m.loadConf(path)
}

func (m *JsonConfigManager) loadConf(path string) {
	files, err := ioutil.ReadDir(path)
	if (err != nil) {
		log.Fatal(err)
	}

	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if strings.Compare(ext, ".json") != 0 {
			continue
		}

		fname := filepath.Join(path, f.Name())
		content, err := ioutil.ReadFile(fname)
		if (err != nil) {
			log.Fatal(err)
		}

		k := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		conf := NewJsonConfig(k)
		conf.loadData(content)
		m.addConf(k, conf)
	}
}

func (m *JsonConfigManager) addConf(name string, conf *JsonConfig) {
	if _, ok := m.confs[name]; ok {
		return
	}

	m.confs[name] = conf
}

func (m *JsonConfigManager) findConf(name string, sn int) interface{} {
	if conf, ok := m.confs[name]; ok {
		return conf
	}

	return nil
}
