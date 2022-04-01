package main

import (
	"time"
)

//----------------------EVENT interface-----------------------
type eventI interface {
	Baby() babyI
	Start() time.Time
}

type event struct {
	baby      baby
	startTime time.Time
}

func newEvent(b baby, t time.Time) event {
	e := event{baby: b, startTime: t}
	return e
}

func (e event) Baby() baby {
	return e.baby
}
func (e event) Start() time.Time {
	return e.startTime
}
