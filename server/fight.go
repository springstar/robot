package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/pb"
	"github.com/springstar/robot/core"
)
type Skill struct {
	sn int32
	level int32
	pos int32
	nextRelease int64
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
	core.Info("recv skill change")
	msg := msg.ParseSCSkillUpdate(int32(msg.MSG_SCSkillUpdate), packet.Data)
	changes := msg.GetChanges()
	for _, c := range changes {
		skill := newSkill(c.GetSkillId(), c.GetLevel(), c.GetSlot(), c.GetSkillCd())
		r.addSkill(skill.sn, skill)
	}

}

func (r *Robot) handleHpChange(packet *core.Packet) {
	core.Info("recv hp change")
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