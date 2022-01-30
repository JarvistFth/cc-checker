package utils

import "time"

type TimeCounter struct {
	startTime int64
	endTime int64
}

func (c *TimeCounter) Begin() {
	c.startTime = time.Now().UnixNano()
}

func (c *TimeCounter) End() (seconds, ms float64){
	c.endTime = time.Now().UnixNano()

	seconds = float64((c.endTime - c.startTime) / 1e9)
	ms = float64((c.endTime - c.startTime) / 1e6)
	return seconds,ms
}
