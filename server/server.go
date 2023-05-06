package server

import (
	"gopkg.in/yaml.v3"
	"log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

var (
	serv *Server
)

type Server struct {
	cfg ServerConfig
	exit chan struct {}
	engine *gin.Engine
	driver *RobotDriver
	accountMgr *AccountManager
	robotMgr *RobotManager

}


func NewServer() *Server {
	serv = &Server{
		engine : gin.Default(),
		driver : NewDriver(),
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
	startAccountId := getStartAccountId(serv.cfg.ServerId)
	serv.accountMgr = newAccountManager(startAccountId, serv.cfg.MaxNum)

	serv.robotMgr = newRobotManager(serv.cfg.Url)

}

func (serv *Server) Run() {
	go serv.driver.Start()
	serv.engine.Run()
	
}

func (serv *Server) PostCommand(cmd iCommand) {
	serv.driver.PostCommand(cmd)
}