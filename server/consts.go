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