package server

import (
	"gopkg.in/yaml.v3"
	"log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	_ "github.com/springstar/robot/core"
)

var (
	serv *Server
)


type Server struct {
	*RunStat
	*InstructionList
	cfg ServerConfig
	exit chan struct {}
	engine *gin.Engine
	driver *RobotDriver
	accountMgr *AccountManager
	robotMgr *RobotManager
	confMgr *JsonConfigManager
	nameMgr *NameManager
	sceneMgr *SceneManager

}


func NewServer() *Server {
	serv = &Server{
		engine : gin.Default(),
		driver : NewDriver(),
		RunStat: newRunStat(),
		exit : make(chan struct{}),

	}

	return serv
}

func (serv *Server) Init() {
	content, err := ioutil.ReadFile("robot.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}
	
	err = serv.parseConfig(content)
	if (err != nil) {
		log.Fatal(err)
		return
	}

	serv.InstructionList = loadInstructions("server/orders.txt")

	serv.initManager()

	for _, api := range apis {
		if (api.method == "GET") {
			serv.engine.GET(api.url, api.handler)			
		} else if (api.method == "POST") {
			serv.engine.POST(api.url, api.handler)
		}
	}
}

func (serv *Server) parseConfig(content []byte) error {
	if err := yaml.Unmarshal([]byte(content), &serv.cfg); err != nil {
		return err
	}
	
	return nil
}

func (serv *Server) initManager() {
	serv.confMgr  = newJsonConfigManager()
	serv.confMgr.init("config")

	startAccountId := getStartAccountId(serv.cfg.ServerId)
	serv.accountMgr = newAccountManager(startAccountId, serv.cfg.MaxNum)

	serv.robotMgr = newRobotManager(serv.cfg.Url)	

	serv.nameMgr = newNameManager()
	serv.nameMgr.loadNameFiles()


}

func (serv *Server) Run() {
	go serv.driver.Start()
	go serv.RunStat.Start()
	serv.engine.Run()
	
}

func (serv *Server) PostCommand(cmd iCommand) {
	serv.driver.PostCommand(cmd)
}

func (serv *Server) getNameManager() *NameManager {
	return serv.nameMgr
}