package server

import (
	_ "bytes"
	
	_ "net"
	"log"
	"time"
	_ "sync"
	_ "context"

	"github.com/springstar/robot/core"
	"github.com/springstar/robot/msg"
	_ "github.com/gobwas/ws"
)

type iExecutor interface {
	exec(params []string, delta int) ExecState
	checkIfExec() bool
}

type Robot struct {
	core.IDispatcher
	*Character
	conn core.NetConnection
	packetQ chan []*core.Packet
	buffer *core.PacketBuffer
	mgr *RobotManager
	account *Account
	fsm *RobotFsm

	moduleMgr *ModuleManager
	ticker *time.Ticker
	quit chan struct{}
	executors map[string]iExecutor
	pc int
	isQuit bool

}

func newRobot(account *Account, robotMgr *RobotManager, fsm *RobotFsm) *Robot {
	r := &Robot{
		IDispatcher: core.NewMsgDispatcher(),
		Character: newCharacter(),
		mgr : robotMgr,
		account : account,
		fsm : fsm,
		packetQ : make(chan []*core.Packet),
		buffer : core.NewBuffer(),
		moduleMgr : newModuleManager(),
		quit: make(chan struct{}),
		pc: -1,
		executors: make(map[string]iExecutor),
		isQuit: false,
	}

	if (r.account != nil) {
		robotMgr.add(r.account.id, r)		
	}

	return r
}

func (r *Robot) loadModules() {
	r.executors["move"] = newMovement(r)
	r.executors["quest"] = newQuestExecutor(r)
	r.executors["match"] = newMatchExecutor(r)
}

func (r *Robot) registerMsgHandler() {
	r.Register(msg.MSG_SCLoginResult, r)
	r.Register(msg.MSG_SCQueryCharactersResult, r)
	r.Register(msg.MSG_SCCharacterCreateResult, r)
	r.Register(msg.MSG_SCCharacterLoginResult, r)
	r.Register(msg.MSG_SCInitData, r)
	r.Register(msg.MSG_SCStageEnterResult, r)
	r.Register(msg.MSG_SCStageSwitch, r)
	r.Register(msg.MSG_SCHumanKick, r)
	r.Register(msg.MSG_SCStageMove, r)
	r.Register(msg.MSG_SCAccountLoginQueue, r)
	r.Register(msg.MSG_SCSoulAwaken, r)
	r.Register(msg.MSG_SCStageObjectAppear, r)
	r.Register(msg.MSG_SCStageObjectDisappear, r)

}

func (r *Robot) startup() {
	r.loadModules()
	r.registerMsgHandler()
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
	case "createchar":
		r.createChar()
	case "sendCharacterLogin":
		r.sendCharacterLogin()
	case "waitForInit":
		r.waitForInit()
	case "enterStage":
		r.enterStage()
	case "ready":
		r.ready()
	case "onswitch":
		r.onSwitchStage()		
	default:
		core.Warn(action)	
	}
}

func (r *Robot) connect() {
	r.conn = core.NewWsConnection()
	err := r.conn.Connect(r.mgr.url)
	if err != nil {
		core.Error(err)
		r.fsm.trigger("connecting", "cfail", r)
	}

	r.fsm.trigger("connecting", "cok", r)
	
}

func (r *Robot) on_connection_established() {
	// log.Printf("connection established")
	r.sendLoginRequest()

	go r.readPackets()
	go r.mainLoop()
}

func (r *Robot)readPackets() {
	Loop:
	for {
		select {
		case <- r.quit:
			r.isQuit = true
			break Loop
		default:
			bytes, err := r.conn.Read()
			if (err != nil) {
				log.Fatal(err)
				break Loop
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
}

func (r *Robot) done() {
	r.quit <- struct{}{}
}

func (r *Robot) dispatch(packets []*core.Packet) {
	for _, packet := range packets {
		r.Dispatch(packet)
	}
}

func (r *Robot) sendPulse() {
	packet := msg.SerializeCSPing(msg.MSG_CSPing)
	r.sendPacket(packet)
}

func (r *Robot)mainLoop() {
	r.ticker = time.NewTicker(ROBOT_PULSE)
	for {
		select {
			case packets := <- r.packetQ:
				r.dispatch(packets)	
			case <- r.ticker.C:	
				if !r.checkQuit() {
					r.update()
				}	
		}
	}
}

func (r *Robot) checkQuit() bool {
	if r.isQuit {
		close(r.packetQ)
		r.conn.Close()
		r.ticker.Stop()		
		return true
	}

	return false
}

func (r *Robot) ready() {
	if r.profession == 0 {
		return
	}	
	
	r.pc = core.GenRandomInt(serv.icount())
}
 
func (r *Robot)update() {
	r.sendPulse()
	r.vm()

	r.ticker.Reset(ROBOT_PULSE)

}

func (r *Robot) findExecutor(cmd string) iExecutor {
	if executor, ok := r.executors[cmd]; !ok {
		return nil
	} else {
		return executor
	}
}

func (r *Robot) vm() {
	if r.pc != -1 {
		instruction := serv.fetch(r.pc)
		executor := r.findExecutor(instruction.cmd)
		if executor == nil {
			core.Error("no executor ", instruction.cmd)
			// log.Fatal("no executor ", instruction.cmd)
		} else {
			if !executor.checkIfExec() {
				return
			}
			
			state := executor.exec(instruction.params, 30)
			if state == EXEC_COMPLETED {
				r.pc, instruction = serv.next(r.pc)
			}
		}
	}	
}

func (r *Robot) HandleMessage(packet *core.Packet) {
	switch packet.Type {
		case msg.MSG_SCLoginResult:
			r.handleLoginResult(packet)
		case msg.MSG_SCAccountLoginQueue:
			r.handleLoginQueue(packet)	
		case msg.MSG_SCQueryCharactersResult:
			r.handleQueryCharacters(packet)
		case msg.MSG_SCCharacterCreateResult:
			r.handleCreateResult(packet)
		case msg.MSG_SCCharacterLoginResult:
			r.handleCharacterLogin(packet)
		case msg.MSG_SCInitData:
			r.handleInitData(packet)	
		case msg.MSG_SCStageEnterResult:
			r.handleEnterStage(packet)	
		case msg.MSG_SCStageSwitch:
			r.handleSwitchStage(packet)
		case msg.MSG_SCStageObjectAppear:
			r.handleObjAppear(packet)				
		case msg.MSG_SCStageObjectDisappear:
			r.handleObjDisappear(packet)
		case msg.MSG_SCHumanKick:
			r.handleKick(packet)
		case msg.MSG_SCStageMove:
			r.handleStageMove(packet)
		case msg.MSG_SCSoulAwaken:
			r.handleSoulAwaken(packet)				
		default:
			core.Warn("recv packet type ", packet.Type)	
	}

	queueMsgStat(STAT_RECV_PACKETS, int32(packet.Type), int32(packet.Length))

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

func (mgr *RobotManager) stopRobots() {
	for _, r := range mgr.robots {
		r.done()
	}

	for k := range mgr.robots {
		delete(mgr.robots, k)
	}


}




