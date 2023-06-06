package server

import (
	_ "fmt"
	"runtime"
	"time"

	"github.com/Jeffail/tunny"
	
)

type RobotDriver struct {
	*RunStat
	cq chan iCommand
	rq chan string
	pool *tunny.Pool
	ticker *time.Ticker
}

func NewDriver() *RobotDriver {
	return &RobotDriver{
		RunStat: newRunStat(),
		cq : make(chan iCommand),
		rq : make(chan string),
		ticker : time.NewTicker(SERVER_PULSE),
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
		case <- driver.ticker.C:
			driver.pulse()
			driver.ticker.Reset(SERVER_PULSE)	

		}	
	}	
}

func (driver *RobotDriver) pulse() {
	if len(driver.ch) == 0 {
		return
	}

	for stat := range driver.ch {
		driver.statistic(stat)
	}
}

func (driver *RobotDriver) exec(cmd iCommand) {
	driver.pool.Process(cmd)
}

func (driver *RobotDriver) genRobots() {

}

