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

	for i := 0; i < cmd.Count; i++ {
		err, account := accountMgr.alloc()
		if (err != nil) {
			break
		}

		robot := newRobot(account, robotMgr, newFsm())
		robot.startup()
	}
}

type ExitCommand struct {

}

func (cmd ExitCommand) exec() {

}