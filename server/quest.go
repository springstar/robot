package server

type Quest struct {
	sn int32
	step int32
	total int32
	status int32
	typ int32
}

func newQuest() *Quest {
	return &Quest{

	}
}

type RobotQuestSet struct {
	quests map[int]*Quest
}

func newQuestSet() *RobotQuestSet{
	return &RobotQuestSet{
		quests: make(map[int]*Quest),
	}
}

func (qs *RobotQuestSet) addQuest(q *Quest) {
	if _, ok := qs.quests[int(q.sn)]; !ok {
		qs.quests[int(q.sn)] = q
	}
}

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
	 