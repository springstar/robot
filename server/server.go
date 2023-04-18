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
	engine *gin.Engine
	driver *RobotDriver
	accountMgr *AccountManager
}


func NewServer() *Server {
	serv = &Server{
		engine : gin.Default(),
		driver : NewDriver(),
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

	startAccountId := getStartAccountId(serv.cfg.ServerId)
	serv.accountMgr = newAccountManager(startAccountId, serv.cfg.MaxNum)

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

func (serv *Server) Run() {
	go serv.driver.Start()
	serv.engine.Run()
	
}

func (serv *Server) PostCommand(cmd iCommand) {
	serv.driver.PostCommand(cmd)
}