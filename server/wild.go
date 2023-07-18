package server

import (
	"github.com/springstar/robot/core"

)



type WildExecutor struct {
	*Executor
}

func newWildExecutor(r *Robot) *WildExecutor {
	return &WildExecutor{
		Executor: newExecutor(r),
	}
}

func (q *WildExecutor) exec(params []string, delta int) {
	mapSn, err := core.Str2Int(params[0])
	if err != nil {
		core.Error("wild map error %s", params)
		return
	}

	if int(q.mapSn) != mapSn {
		q.switchStage(1, mapSn, mapSn)
		q.setPause()	
	}

	if q.getState() == EXEC_PAUSE {
		// core.Info("monster quest attach ctx function ", q.mapSn)
		q.attachCtxFun(asyncWild, q)
	} else {

	}


}

func (q *WildExecutor) resume(params []string, delta int) {
	q.Executor.resume()
}

func asyncWild(e iExecutor) {
	q := e.(*WildExecutor)
	q.execWildFight()
}

func (q *WildExecutor) execWildFight() {

}



