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
	*Robot
	exstates map[int]ExecState
}

func newExecutor(r *Robot) *Executor {
	return &Executor{
		Robot: r,
		exstates: make(map[int]ExecState),
	}
}

func (e *Executor) exec(params []string, delta int) ExecState {
	return EXEC_COMPLETED
}

func (e *Executor) handleBreak(params []string) {
	
}

func (e *Executor) checkIfExec() bool {
	return false
}

func (e *Executor) setOngoing() {
	e.exstates[e.pc] = EXEC_ONGOING
}

func (e *Executor) setCompleted() {
	e.exstates[e.pc] = EXEC_COMPLETED
}

func (e *Executor) getState() ExecState {
	return e.exstates[e.pc]
}

type AsyncExecutor struct {
	*Executor
}

func newAsyncExecutor(r *Robot) *AsyncExecutor {
	return &AsyncExecutor{
		Executor: newExecutor(r),
	}
}

func (ae *AsyncExecutor) checkIfExec(params []string) bool {
	if ae.getState() == EXEC_ONGOING {
		return false
	}

	return true
}

func (r *Robot) vm() {
	if r.pc != -1 {
		instruction := r.fetch(r.pc)
		executor := r.findExecutor(instruction.cmd)
		if executor == nil {
			core.Error("no executor ", instruction.cmd)
		} else {
			if !executor.checkIfExec(instruction.params) {
				executor.handleBreak()
				r.pc, _ = r.next(r.pc)
				return
			}
			
			state := executor.exec(instruction.params, 30)
			if state == EXEC_COMPLETED {
				r.pc, instruction = r.next(r.pc)
			}
		}
	}	
}