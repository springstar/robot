package server

import (
	"github.com/springstar/robot/core"
)

type ExecState int32

const (
	EXEC_NO_START ExecState = iota
	EXEC_ONGOING // 进行中，不可重入
	EXEC_REPEATED // 可重入
	EXEC_PAUSE
	EXEC_RESUME
	EXEC_COMPLETED
)

type iExecutor interface {
	exec(params []string, delta int)
	resume(params []string, delta int)
	getStatus(pc int) ExecState
	onEvent(k EventKey)
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

func (e *Executor) exec(params []string, delta int) {
	
}

func (e *Executor) resumt(params []string, delta int) {

}

func (e *Executor) getStatus(pc int) ExecState {
	return e.exstates[pc]
}

func (e *Executor) onEvent(k EventKey) {

}


func (e *Executor) setRepeated() {
	e.exstates[e.pc] = EXEC_REPEATED
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

type AsyncContext struct {

}

func newAsyncContext() *AsyncContext {
	return &AsyncContext{

	}
}


func (r *Robot) vm() {
	if r.pc != -1 {
		instruction := r.fetch(r.pc)
		executor := r.findExecutor(instruction.cmd)
		if executor == nil {
			core.Error("no executor ", instruction.cmd)
			return
		}

		status := executor.getStatus(r.pc)

		switch status {
		case EXEC_NO_START, EXEC_REPEATED:
			executor.exec(instruction.params, 30)
		case EXEC_RESUME:
			executor.resume(instruction.params, 30)	
		default:
			break	
		}
		
		status = executor.getStatus(r.pc)
		if status == EXEC_COMPLETED {
			r.pc, _ = r.next(r.pc)
			executor.exec(instruction.params, 30)
		}
	}	

	
}