package main

import (
	"git.ddd.rip/ptrcnull/watchface/system"
	"time"
)

var SimpleClockTime = Text(108, 162, 5, 36, LightGray)
var SimpleClockBattery = Text(90, 240, 3, 24, LightGray)

type SimpleClock struct{}

func (s SimpleClock) Init(face *Face) {
	s.draw(face, time.Now())
}

func (s SimpleClock) Render(face *Face) {
	t := time.Now()
	if t.Second() == 0 {
		s.draw(face, t)
	}
}

func (SimpleClock) draw(face *Face, t time.Time) {
	face.Text(SimpleClockTime, t.Format("15:04"))
	face.Text(SimpleClockBattery, system.GetBattery())
}
