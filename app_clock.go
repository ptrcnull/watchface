package main

import (
	"git.ddd.rip/ptrcnull/watchface/system"
	"time"
)

var ClockTime = Text(62, 162, 8, 36, LightGray)
var ClockBattery = Text(90, 240, 3, 24, LightGray)

type Clock struct{}

func (s Clock) Init(face *Face) {
	s.draw(face, time.Now())
}

func (s Clock) Render(face *Face) {
	t := time.Now()
	if t.Second() == 0 {
		s.draw(face, t)
	}
}

func (Clock) draw(face *Face, t time.Time) {
	face.Text(ClockTime, t.Format("15:04"))
	face.Text(ClockBattery, system.GetBattery())
}
