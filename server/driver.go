package server

import (
	"runtime"
	"time"

	"github.com/Jeffail/tunny"
	
)

type RobotDriver struct {
	cq chan iCommand
	pool *tunny.Pool
	ticker *time.Ticker
}

func NewDriver() *RobotDriver {
	return &RobotDriver{
		cq : make(chan iCommand),
		ticker : time.NewTicker(SERVER_PULSE),
	}
}

func (driver *RobotDriver) Start() {
	numCPUs := runtime.NumCPU()
	driver.pool = tunny.NewFunc(numCPUs, run)
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
		driver.ticker.Reset(SERVER_PULSE)
		select {
		case cmd := <- driver.cq:
			driver.exec(cmd)
		case <- time.Tick(SERVER_PULSE):
			driver.pulse()
		}
	
	}
	
}

func (driver *RobotDriver) pulse() {
	
}

func (driver *RobotDriver) exec(cmd iCommand) {
	driver.pool.Process(cmd)
}

func (driver *RobotDriver) genRobots() {

}

