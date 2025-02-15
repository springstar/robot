package server

import (
	"github.com/springstar/robot/config"
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/pb"
	"github.com/springstar/robot/core"
)

type SkillTargetType int32

const (
	STT_ENEMY = 1
	STT_FRIEND = 2
)
type Skill struct {
	sn int32
	level int32
	pos int32
	nextRelease int64
}

type SkillQuestData struct {
	sn int
}

func newSkillQuestData() *SkillQuestData {
	return &SkillQuestData{}
}

func (d *SkillQuestData) resume(executor *RobotQuestExecutor) {

}

func (d *SkillQuestData) getQuestSn() int {
	return d.sn
}

func (d *SkillQuestData) onStatusUpdate(executor *RobotQuestExecutor, sn int, status QuestStatus) {
	if status == QSTATE_COMPLETED{
		executor.commitQuest(sn)
	}
}


func newSkill(sn int32, lv int32, pos int32, nextRelease int64) *Skill {
	return &Skill{
		sn: sn,
		level: lv, 
		pos: pos,
		nextRelease: nextRelease,
	}
}

func (r *Robot) handleSkillUpdate(packet *core.Packet) {
	msg := msg.ParseSCSkillUpdate(int32(msg.MSG_SCSkillUpdate), packet.Data)
	changes := msg.GetChanges()
	for _, c := range changes {
		skill := newSkill(c.GetSkillId(), c.GetLevel(), c.GetSlot(), c.GetSkillCd())
		r.addSkill(skill.sn, skill)
	}

}

func (r *Robot) handleHpChange(packet *core.Packet) {
	// resp := msg.ParseSCFightHpChg(int32(msg.MSG_SCFightHpChg), packet.Data)
	// targets := resp.GetDhpChgTar()
	// for _, t := range targets {
	// 	core.Info("hp chg ", t.Id)
	// }
}

func (r *Robot) initSkills(skillA []*pb.DSkill, skillB []*pb.DSkill) {
	for _, sa := range skillA {
		skill := newSkill(sa.GetSkillSn(), sa.GetSkillLevel(), sa.GetPosition(), sa.GetNextRealse())
		r.addSkill(skill.sn, skill)
	}

	for _, sb := range skillB {
		skill := newSkill(sb.GetSkillSn(), sb.GetSkillLevel(), sb.GetPosition(), sb.GetNextRealse())
		r.addSkill(skill.sn, skill)
	}
}

func (r *Robot) addSkill(sn int32, s *Skill) {
	r.skills[sn] = s
}

func (r *Robot) fight(enemyId int64) {
	enemy := r.findObj(enemyId)
	if enemy == nil {
		return
	}

	sn := r.pickSkill()
	if sn == 0 {
		// core.Info("no skill pick, enemy id ", enemyId, r.lastRelease)
		// for _, s := range r.skills {
		// 	core.Info("skill nextRelease ", s.sn, s.nextRelease)
		// }
		return
	}

	confSkill := config.FindConfSkill(int(sn))

	tarId := enemyId
	tarPos := enemy.getPos()
	dirX := tarPos.X - r.pos.X
	dirY := tarPos.Y - r.pos.Y
	dir := &pb.DVector2{}
	dir.X = dirX
	dir.Y = dirY
	spos := &pb.DVector2{}
	spos.X = r.pos.X
	spos.Y = r.pos.Y

	distance := r.pos.DistanceTo(tarPos)

	if confSkill.TargetType == STT_FRIEND {
		tarId = r.humanId
		tarPos = r.pos
	}

	if distance > float32(confSkill.Range) + 1.5 {
		return
	}

	// core.Info("fighting robot pos ", r.pos.X, r.pos.Y, distance)
	// core.Info("fighting target pos ", tarPos.X, tarPos.Y, distance)

	tpos := &pb.DVector2{}
	tpos.X = tarPos.X
	tpos.Y = tarPos.Y
	
	msg := msg.SerializeCSFightAtk(uint32(msg.MSG_CSFightAtk), r.humanId, sn, tarId, tpos, 0, false, dir, spos, 1)
	r.sendPacket(msg)
	// core.Info("send CSFightAtk to attack ", sn, tarId, r.lastRelease)
	r.updateLastRelease(int(sn))
	r.updateSkillCooling(int(sn))

}

func (r *Robot) dumpSkills() {
	core.Info("dump skills ", len(r.skills))
	for sn, skill := range r.skills {
		core.Info("skill ", sn, skill.level, skill.nextRelease)
	}
}

func (r *Robot) findSkill(sn int) *Skill {
	if sk, ok := r.skills[int32(sn)]; ok {
		return sk
	}

	return nil
}

func (r *Robot) resetSkillCooling(sn int) {
	sk := r.findSkill(sn)
	if sk != nil {
		sk.nextRelease = 0
	}
}

func (r *Robot) updateLastRelease(sn int) {
	confSkill := config.FindConfSkill(sn)
	r.lastRelease = core.GetCurrentTime()
	r.commonCooolTime = confSkill.PublicCoolTime
	if confSkill.GroupType == 1 {
		r.commonCooolTime = 3000
	}
	// core.Info("update last release ", r.lastRelease)

}

func (r *Robot) resetLastRelease() {

}

func (r *Robot) updateSkillCooling(sn int) {
	confSkill := config.FindConfSkill(sn)
	coolInfo, _ := core.Str2IntSlice(confSkill.CoolTime)
	coolTime := 0
	for i := 0; i < len(coolInfo); i++ {
		coolTime += coolInfo[i]
	}

	sk := r.findSkill(sn) 
	if sk != nil {
		sk.nextRelease = core.GetCurrentTime() + int64(coolTime)
		// core.Info("update next release ", sn, sk.nextRelease)
	}
}

func (r *Robot) pickSkill() int32 {
	// r.dumpSkills()
	now := core.GetCurrentTime()
	for sn, skill := range r.skills {
		confSkill := config.FindConfSkill(int(sn))
		if confSkill == nil {
			continue
		}

		if !confSkill.Active {
			continue
		}

		if confSkill.TargetType == STT_FRIEND {
			continue
		}

		if skill.level == 0 {
			continue
		}

		if r.lastRelease > 0 && r.lastRelease + int64(r.commonCooolTime) > now {
			continue
		}

		if skill.nextRelease > 0 && skill.nextRelease > now {
			continue
		}

		return sn
	}

	return 0
}

func (r *Robot) upgradeSkill() {
	for sn, _ := range r.skills {
		request := msg.SerializeCSSkillLevelup(msg.MSG_CSSkillLevelup, sn)
		r.sendPacket(request)
	}
}

func (r *Robot) handleDeath(packet *core.Packet) {
	resp := msg.ParseSCStageObjectDead(int32(msg.MSG_SCStageObjectDead), packet.Data)
	objId := resp.GetObjId()
	obj := r.findObj(objId)
	if obj == nil {
		return
	}

	if obj.getType() == WOT_MONSTER {
		monsterObj := obj.(*MonsterObj)
		monsterObj.curHp = 0
	}
}

func (r *Robot) handleFightResult(packet *core.Packet) {
	// resp := msg.ParseSCFightAtkResult(int32(msg.MSG_SCFightAtkResult), packet.Data)
	// code := resp.GetResultCode()
	// sn := resp.GetSkillSn()
	// if code == -1 {
	// 	r.resetSkillCooling(int(sn))
	// }
}