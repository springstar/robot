package server

import (
	_ "bytes"
	
	_ "net"
	"fmt"
	"log"
	"github.com/springstar/robot/core"
	_ "github.com/springstar/robot/msg"
	_ "github.com/gobwas/ws"
)

type Robot struct {
	core.IDispatcher
	conn core.NetConnection
	mgr *RobotManager
	account *Account
	fsm *RobotFsm
	packetQ chan []*core.Packet
	buffer *core.PacketBuffer
	moduleMgr *ModuleManager

}

func newRobot(account *Account, robotMgr *RobotManager, fsm *RobotFsm) *Robot {
	r := &Robot{
		IDispatcher: core.NewMsgDispatcher(),
		mgr : robotMgr,
		account : account,
		fsm : fsm,
		packetQ : make(chan []*core.Packet),
		buffer : core.NewBuffer(),
		moduleMgr : newModuleManager(),

	}

	if (r.account != nil) {
		robotMgr.add(r.account.id, r)		
	}

	return r
}

func (r *Robot) loadModules() {

}

func (r *Robot) startup() {
	r.Register(112, r)
	r.Register(1004, r)

	r.fsm.trigger("entry", "connect", r)
}

func (r *Robot) doAction(action string) {
	switch action {
	case "connect":
		r.connect()
	case "on_connection_established":
		r.on_connection_established()
	case "login":
		r.sendLoginRequest()
	case "querychars":
		r.querychars()
	default:
		fmt.Println(action)	
	}
}

func (r *Robot) connect() {
	r.conn = core.NewWsConnection()
	err := r.conn.Connect(r.mgr.url)
	if err != nil {
		fmt.Print(err)
		r.fsm.trigger("connecting", "cfail", r)
	}

	r.fsm.trigger("connecting", "cok", r)
	
}

func (r *Robot) on_connection_established() {
	go r.startWork()
	r.mainLoop()
}

func (r *Robot) startWork() {
	r.sendLoginRequest()

	for {
		bytes, err := r.conn.Read()
		if (err != nil) {
			log.Fatal(err)
		}

		if len(bytes) <= 0 {
			continue
		}

		// add to msg buffer
		r.buffer.Write(bytes)

		// split packet from msg buffer and send to packetQ channel
		packets := r.buffer.Read()
		if (packets == nil) {
			continue
		}

		r.packetQ <- packets
		
	}
}

func (r *Robot) dispatch(packets []*core.Packet) {
	for _, packet := range packets {
		r.Dispatch(packet)
	}
}

func (r *Robot) mainLoop() {
	for {
		select {
		case packets := <- r.packetQ:
			r.dispatch(packets)	
		}
	}
}

func (r *Robot) HandleMessage(packet *core.Packet) {
	switch packet.Type {
	case 112:
		r.handleLoginResult(packet)
	case 1004:
		r.handleQueryCharacters(packet)
	}
}

func (r *Robot) sendPacket(packet []byte) {
	r.conn.Write(packet)
}

type RobotManager struct {
	robots map[int]*Robot
	url string
}

func newRobotManager(url string) *RobotManager {
	return &RobotManager{
		robots : make(map[int]*Robot),
		url : url,
	}
}

func (mgr *RobotManager)add(account int, robot *Robot) {
	mgr.robots[account] = robot
}




