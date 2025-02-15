package server

import (
	"log"
	"fmt"
	"encoding/json"
)

type CommandType int32

const (
	COMMAND_TEST  CommandType = 1
	COMMAND_BENCH CommandType = 2
	COMMAND_REPORT CommandType = 3
)

type iCommand interface {
	getType() CommandType
	repeated() int
	exec()
}

type Command struct {
	iCommand
	typ CommandType
}

func (cmd Command) getType() CommandType {
	return cmd.typ
}

func (cmd Command) repeated() int {
	return 1
}

func (cmd Command) exec() {
	fmt.Println(cmd.getType())	
}


type TestCommand struct {
	Command
	desc string
}


type BenchCommand struct {
	Command
	Batch int `start:"batch"`
	Count int `start:"count"`
	Interval int32 `start:"interval"`
	Url string `start:"url"`
}

func (cmd BenchCommand) repeated() int {
	return cmd.Batch
}

func (cmd BenchCommand) exec() {
	for i := 0; i < cmd.Batch; i++ {
		cmd.runRobots()	
	}
}

func (cmd BenchCommand) runRobots() {
	accountMgr := serv.accountMgr
	robotMgr := serv.robotMgr

	count := 0
	for i := 0; i < cmd.Count; i++ {
		err, account := accountMgr.alloc()
		if (err != nil) {
			log.Fatal(err)
			break
		}

		robot := newRobot(account, robotMgr, newFsm())
		robot.startup()
		count += 1
	}

	queueStat(STAT_SEND_ROBOTS, int32(count))
}

type DebugCommand struct {
	Command
	Account string `account:"account name"`
}

func (cmd DebugCommand) exec() {
	name := cmd.Account
	account := serv.accountMgr.findAccountByName(name)
	if account == nil {
		_, account = serv.accountMgr.allocName(name)

	}

	robot := newRobot(account, serv.robotMgr, newFsm())
	robot.startup()
}


type ExitCommand struct {

}

func (cmd ExitCommand) exec() {

}

type StopCommand struct {
	Command
}

func (cmd StopCommand) exec() {
	serv.robotMgr.stopRobots()
}

type QuitCommand struct {
	Command
}

func (cmd QuitCommand) exec() {

}

type ReportCommand struct {
	Command
}

func (cmd ReportCommand) exec() {
	b, err := json.Marshal(serv.RunStat)
	if err != nil {
		serv.driver.rq <- err.Error()
		return
	}


	serv.driver.rq <- string(b)
}

