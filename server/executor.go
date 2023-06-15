package server

import (
	"github.com/springstar/robot/core"
)

type iExecutor interface {
	exec(params []string, delta int) ExecState
	checkIfExec(params []string) bool
	handleBreak()
}

type Executor struct {
	ExecState
}

func (e *Executor) exec(params []string, delta int) ExecState {
	return EXEC_COMPLETED
}

func (e *Executor) handleBreak(params []string) {
	
}

func (e *Executor) checkIfExec() bool {
	return false
}

type RepeatedExecutor struct {
	*Executor
}

type AsyncExecutor struct {
	*Executor
}

func (r *Robot) vm() {
	if r.pc != -1 {
		instruction := serv.fetch(r.pc)
		executor := r.findExecutor(instruction.cmd)
		if executor == nil {
			core.Error("no executor ", instruction.cmd)
		} else {
			if !executor.checkIfExec(instruction.params) {
				executor.handleBreak()
				return
			}
			
			state := executor.exec(instruction.params, 30)
			if state == EXEC_COMPLETED {
				r.pc, instruction = serv.next(r.pc)
			}
		}
	}	
}