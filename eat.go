package main

import (
	"fmt"
	"time"
)

//----------------------EAT interface-----------------------
type eatI interface {
	eventBaseWorker
	eventI
	SetDescription(string)

	Id() int64
	Description() string
}
type eat struct {
	event
	id          int64
	description string
}

func newEat(initEvent event) *eat {
	e := &eat{event: initEvent}
	return e
}

func (e *eat) SetDescription(d string) {
	e.description = d
}

func (e eat) Description() string {
	return e.description
}

func (e eat) Id() int64 {
	return e.id
}

func (e *eat) writeStructToBase() error {
	query_string := fmt.Sprintf("insert into eat(baby_id, start, description) "+
		"VALUES (%d, '%s', '%s') RETURNING id;",
		e.event.BabyId(), e.event.Start().Format("2006-01-02"), e.Description())
	row, err := DBInsertAndGet(query_string)
	if err != nil {
		return err
	}
	var eatId int64
	if err = row.Scan(&eatId); err != nil {
		return err
	}
	e.id = eatId
	return nil
}

//fill current eat by id
func (e *eat) readStructFromBase(id int64) error {

	queryString := fmt.Sprintf("select baby_id, start, description "+
		"from eat where id = %d", id)
	row, err := DBInsertAndGet(queryString)
	if err != nil {
		return err
	}

	var babyId int64
	var start time.Time
	var description string

	if err := row.Scan(&babyId, &start, &description); err != nil {
		return err
	}

	e.SetBabyId(babyId)
	e.SetStart(start)
	e.SetDescription(description)

	return nil
}
