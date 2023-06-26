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
	*AsyncContext
	exstates map[int]ExecState
}

func newExecutor(r *Robot) *Executor {
	return &Executor{
		Robot: r,
		AsyncContext: newAsyncContext(),
		exstates: make(map[int]ExecState),
	}
}

func (e *Executor) exec(params []string, delta int) {
	
}

func (e *Executor) resumt(params []string, delta int) {
	e.AsyncContext.resume()
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

func (e *Executor) setPause() {
	e.exstates[e.pc] = EXEC_PAUSE
}

func (e *Executor) setResume() {
	e.exstates[e.pc] = EXEC_RESUME
}

func (e *Executor) setCompleted() {
	e.exstates[e.pc] = EXEC_COMPLETED
}

func (e *Executor) getState() ExecState {
	return e.exstates[e.pc]
}

type AsyncContext struct {
	ctxFunc func(iExecutor)
	e iExecutor
}

func newAsyncContext() *AsyncContext {
	return &AsyncContext{

	}
}

func (ac *AsyncContext) attachCtxFun(f func(iExecutor), e iExecutor) {
	ac.ctxFunc = f
	ac.e = e
}

func (ac *AsyncContext) resume() {
	if ac.ctxFunc != nil {
		ac.ctxFunc(ac.e)
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