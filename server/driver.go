package server

import (
	_ "fmt"
	"runtime"
	

	"github.com/Jeffail/tunny"
	
)

type RobotDriver struct {
	cq chan iCommand
	rq chan string
	pool *tunny.Pool
}

func NewDriver() *RobotDriver {
	return &RobotDriver{
		cq : make(chan iCommand),
		rq : make(chan string),
	}
}

func (driver *RobotDriver) Start() {
	numCPUs := runtime.NumCPU()
	driver.pool = tunny.NewFunc(numCPUs, run)
	defer driver.pool.Close()
	driver.process()
}

func run(i interface{}) interface{} {
	cmd := i.(iCommand)
	cmd.exec()
	return cmd
}

func (driver *RobotDriver) PostCommand(cmd iCommand) {
	driver.cq <- cmd
}

func (driver *RobotDriver) process() {
	for {
		select {
		case cmd := <- driver.cq:
			driver.exec(cmd)
		}	
	}	
}

func (driver *RobotDriver) exec(cmd iCommand) {
	driver.pool.Process(cmd)
}

func (driver *RobotDriver) genRobots() {

}

