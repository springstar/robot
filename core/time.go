package core

import (
	"time"
)

func Now() int64 {
	now := time.Now()
	return now.Unix()
}

func GetCurrentTime() int64 {
	now := time.Now()
	msec := now.UnixMilli()
	return msec
}