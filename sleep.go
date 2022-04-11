package main

import (
	"fmt"
	"time"
)

//----------------------SLEEP interface-----------------------
type sleepI interface {
	eventBaseWorker
	eventI
	setEndTime(time.Time)
	updateEndSleepTime() error
	calcDuration()

	Duration() time.Duration
	End() time.Time
	Id() int64
}

type sleep struct {
	event
	endTime  time.Time
	duration time.Duration
	id       int64
}

func newSleep(initEvent event) *sleep {
	s := &sleep{event: initEvent}
	return s
}

func (s *sleep) updateEndSleepTime() error {
	query_string := fmt.Sprintf("update sleep set sleep_end = '%s' "+
		"WHERE id = %d", s.End().Format("2006-01-02 03:04"), s.Id())
	_, err := DBInsertAndGet(query_string)
	if err != nil {
		return err
	}
	return nil

}

//read sleep by id
func (s *sleep) readStructFromBase(id int64) error {

	/* query_string := fmt.Sprintf("select (baby_id, start, sleep_end) "+
	"from sleep where sleep_start > '%s' and sleep_start < ('%s' + '1 day'::interval",
	date, date) */
	query_string := fmt.Sprintf("select (baby_id, start, sleep_end) "+
		"from sleep where id = %d", id)

	row, err := DBReadRow(query_string)
	if err != nil {
		return err
	}
	var babyId int64
	var start time.Time
	var sleepEnd time.Time
	if err := row.Scan(&babyId, &start, &sleepEnd); err != nil {
		return err
	}
	s.SetBabyId(babyId)
	s.SetStart(start)
	s.setEndTime(sleepEnd)

	return nil
}

func (s *sleep) writeStructToBase() error {
	query_string := fmt.Sprintf("insert into sleep (baby_id, start) "+
		"values (%d, '%s') RETURNING id", s.BabyId(), s.Start().Format("2006-01-02 03:04"))
	fmt.Println(query_string)
	pIdRow, err := DBInsertAndGet(query_string)
	if err != nil {
		return err
	}

	var sleepId int64

	err = pIdRow.Scan(&sleepId)
	if err != nil {
		return err
	}
	s.id = sleepId
	return nil
}

func (s *sleep) setEndTime(t time.Time) {
	s.endTime = t
	s.calcDuration()
}

func (s *sleep) calcDuration() {
	s.duration = s.End().Sub(s.Start())
}

func (s sleep) Duration() time.Duration {
	return s.duration
}
func (s sleep) End() time.Time {
	return s.endTime
}
func (s sleep) Id() int64 {
	return s.id
}
