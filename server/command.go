package server

import (
	"fmt"
)

type CommandType int32

const (
	COMMAND_TEST  CommandType = 1
	COMMAND_BENCH CommandType = 2
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
	batch int `start:"batch"`
	count int `start:"count"`
	interval int32 `start:"interval"`
	url string `start:"url"`
}

func (cmd BenchCommand) repeated() int {
	return cmd.batch
}

func (cmd BenchCommand) exec() {
	for i := 0; i < cmd.batch; i++ {
		cmd.runRobots()	
	}
}

func (cmd BenchCommand) runRobots() {
	accountMgr := serv.accountMgr
	robotMgr := serv.robotMgr

	for i := 0; i < cmd.count; i++ {
		err, account := accountMgr.alloc()
		if (err != nil) {
			break
		}
		
		robot := newRobot(account, robotMgr, newDefaultStateMachine(defaultTransitions()))
		robot.startup()
	}
}