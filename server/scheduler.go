package server

type Scheduler struct {
	
}

type WaitQuestData struct {
	questSn int
	waitStart int64
}

func newWaitQuestData() *WaitQuestData {
	return &WaitQuestData{

	}
}

func (d *WaitQuestData) resume(executor *RobotQuestExecutor) {

}

func (d *WaitQuestData) getQuestSn() int {
	return d.questSn
}

func (d *WaitQuestData) onStatusUpdate(executor *RobotQuestExecutor, sn int, status QuestStatus) {
	if status == QSTATE_COMPLETED{
		executor.commitQuest(sn)
	}
}


