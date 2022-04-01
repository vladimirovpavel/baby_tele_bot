package main

import (
	"database/sql"
	"fmt"
	"time"
)

//----------------------EAT interface-----------------------
type eatI interface {
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

func newEat(initEvent event) eat {
	e := eat{event: initEvent}
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
		"VALUES (%d, '%s', '%s') RETURNING eat_id;",
		e.event.Baby().Id(), e.event.Start().Format("2006-01-02"), e.Description())
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
func (e eat) readStructFromBase(query interface{}) (interface{}, error) {
	q := query.([]string)

	babyId := q[0]
	eatTime := q[1]
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	res, err := db.Query("select baby_id, start, description "+
		"from eat where start > %s, baby_id == %s",
		eatTime, babyId)
	if err != nil {
		return nil, err
	}
	var babyid int64
	var start time.Time
	var description string
	if err := res.Scan(&babyid, &start, &description); err != nil {
		return nil, err
	}

	return nil, nil
}
