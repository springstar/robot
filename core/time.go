package core

import (
	"time"
)

func GetCurrentTime() int64 {
	now := time.Now()
	msec := now.UnixMilli()
	return msec
}