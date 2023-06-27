package server

import (
	"time"
)

const (
	SERVER_PULSE  = 30 * time.Millisecond
	ROBOT_PULSE = 3 * 1000 * time.Millisecond
	ROBOT_PREFIX = "robot"
	ROBOT_SECTION = 10000
	ROLE_SEX_MALE = 0
	ROLE_SEX_FEMALE = 1
)


type CampType int32

const (
	CAMP_ESCORTNPC = 1
	CAMP_TRAP = 2
	CAMP_GEAR = 3
	CAMP_ENEMY = 4
	CAPM_MONSTER = 5
	CAMP_ESCORTENEMY = 6
	CAMP_PARTNER = 7
	CAMP_NEUTRAL = 8
	CAMP_HUMAN = 9
)