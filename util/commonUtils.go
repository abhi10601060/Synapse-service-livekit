package util

import "time"

func GetCurrentTime() string {
	return time.Now().Format("01-02-2006 15:04:05")
}