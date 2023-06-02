package server

type RobotQuestHandler struct {

}

func newQuestHandler() *RobotQuestHandler {
	return &RobotQuestHandler{}
}

func (q *RobotQuestHandler) exec(params []string, delta int) ExecState {
	return EXEC_COMPLETED
}
	 