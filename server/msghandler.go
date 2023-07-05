package server

import (

	"github.com/springstar/robot/core"
	"github.com/springstar/robot/msg"
)

func (r *Robot) registerMsgHandler() {
	r.Register(msg.MSG_SCLoginResult, r)
	r.Register(msg.MSG_SCQueryCharactersResult, r)
	r.Register(msg.MSG_SCCharacterCreateResult, r)
	r.Register(msg.MSG_SCCharacterLoginResult, r)
	r.Register(msg.MSG_SCInitData, r)
	r.Register(msg.MSG_SCStageEnterResult, r)
	r.Register(msg.MSG_SCStageSwitch, r)
	r.Register(msg.MSG_SCHumanKick, r)
	r.Register(msg.MSG_SCStageMove, r)
	r.Register(msg.MSG_SCStageMoveStop, r)
	r.Register(msg.MSG_SCAccountLoginQueue, r)
	r.Register(msg.MSG_SCSoulAwaken, r)
	r.Register(msg.MSG_SCStageObjectAppear, r)
	r.Register(msg.MSG_SCStageObjectDisappear, r)
	r.Register(msg.MSG_SCMatchEnrollResponse, r)
	r.Register(msg.MSG_SCMatchResult, r)
	r.Register(msg.MSG_SCInformMsg, r)
	r.Register(msg.MSG_SCTeamMine, r)
	r.Register(msg.MSG_SCPlatTeamListResponse, r)
	r.Register(msg.MSG_SCFightHpChg, r)
	r.Register(msg.MSG_SCQuestInfo, r)
	r.Register(msg.MSG_SCGatherFirst, r)
	r.Register(msg.MSG_SCGatherSecond, r)
	r.Register(msg.MSG_SCSkillUpdate, r)
	r.Register(msg.MSG_SCStageObjectDead, r)
}

func (r *Robot) HandleMessage(packet *core.Packet) {
	switch packet.Type {
		case msg.MSG_SCLoginResult:
			r.handleLoginResult(packet)
		case msg.MSG_SCAccountLoginQueue:
			r.handleLoginQueue(packet)	
		case msg.MSG_SCQueryCharactersResult:
			r.handleQueryCharacters(packet)
		case msg.MSG_SCCharacterCreateResult:
			r.handleCreateResult(packet)
		case msg.MSG_SCCharacterLoginResult:
			r.handleCharacterLogin(packet)
		case msg.MSG_SCInitData:
			r.handleInitData(packet)	
		case msg.MSG_SCStageEnterResult:
			r.handleEnterStage(packet)	
		case msg.MSG_SCStageSwitch:
			r.handleSwitchStage(packet)
		case msg.MSG_SCStageObjectAppear:
			r.handleObjAppear(packet)				
		case msg.MSG_SCStageObjectDisappear:
			r.handleObjDisappear(packet)
		case msg.MSG_SCHumanKick:
			r.handleKick(packet)
		case msg.MSG_SCStageMove:
			r.handleStageMove(packet)
		case msg.MSG_SCStageMoveStop:
			r.handleMoveStop(packet)	
		case msg.MSG_SCStageSetPos:
			r.handleStageSetPos(packet)	
		case msg.MSG_SCStagePullTo:
			r.handleStagePull(packet)					
		case msg.MSG_SCSoulAwaken:
			r.handleSoulAwaken(packet)
		case msg.MSG_SCMatchEnrollResponse:
			r.handleArenaEnroll(packet)					
		case msg.MSG_SCMatchResult:
			r.handleArenaMatchResult(packet)	
		case msg.MSG_SCInformMsg:
			r.handleInform(packet)
		case msg.MSG_SCTeamMine:
			r.handleTeamDetail(packet)
		case msg.MSG_SCPlatTeamListResponse:
			r.handleTeamList(packet)		
		case msg.MSG_SCFightHpChg:
			r.handleHpChange(packet)
		case msg.MSG_SCQuestInfo:
			r.handleQuestInfo(packet)
		case msg.MSG_SCRemoveQuest:
			r.handleRemoveQuest(packet)		
		case msg.MSG_SCSkillUpdate:
			r.handleSkillUpdate(packet)
		case msg.MSG_SCStageObjectDead:
			r.handleDeath(packet)	
		case msg.MSG_SCGatherFirst:
			r.HandleGatherFirst(packet)		
		case msg.MSG_SCGatherSecond:
			r.HandleGatherSecond(packet)	

		default:
			core.Warn("recv packet type ", packet.Type)	
	}

	queueMsgStat(STAT_RECV_PACKETS, int32(packet.Type), int32(packet.Length))

}