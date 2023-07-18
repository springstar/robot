package server



type WildExecutor struct {
	*Executor
}

func newWildExecutor(r *Robot) *WildExecutor {
	return &WildExecutor{
		Executor: newExecutor(r),
	}
}

func (q *WildExecutor) exec(params []string, delta int) {
}

func (q *WildExecutor) resume(params []string, delta int) {
	// core.Info("RobotQuestExecutor resume exec")
	q.Executor.resume()
}


