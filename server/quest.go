package server

type RobotQuestHandler struct {
	r *Robot
}

func newQuestHandler(r *Robot) *RobotQuestHandler {
	return &RobotQuestHandler{
		r: r,
	}
}

func (q *RobotQuestHandler) exec(params []string, delta int) ExecState {

	return EXEC_COMPLETED
}
	 