package gettime

import "time"

func GetTimeStr() (string) {
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	return timeNow
}