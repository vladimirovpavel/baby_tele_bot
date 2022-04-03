package main

import (
	"fmt"
	"time"
)

//----------------------EVENT interface-----------------------
type eventI interface {
	SetBabyId(int64)
	SetStart(time.Time)

	BabyId() int64
	Start() time.Time

	String() string

	//SelectByBabyDate(queryParam []string) eventI
}

type event struct {
	babyId int64
	start  time.Time
}

// RETURNS POINTER TO EVENT!!!!!!
func newEvent(b_id int64, t time.Time) *event {
	e := event{babyId: b_id, start: t}
	return &e
}

func (e event) BabyId() int64 {
	return e.babyId
}
func (e event) Start() time.Time {
	return e.start
}
func (e *event) SetBabyId(id int64) {
	e.babyId = id
}
func (e *event) SetStart(t time.Time) {
	e.start = t
}

func (e event) String() string {
	return fmt.Sprintf("Event of baby %d, started ad %s", e.babyId, e.Start().Format(("2006-01-02")))
}
