package main

import (
	"database/sql"
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

	CheckExisting(table string, id int64) bool

	//SelectByBabyDate(queryParam []string) eventI
}

type event struct {
	babyId int64
	start  time.Time
}

func newEvent(b_id int64) *event {
	e := event{babyId: b_id}
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
	return fmt.Sprintf("Event of baby %d, started ad %s",
		e.babyId, e.Start().Format(("2006-01-02")))
}

func (e event) CheckExisting(table string, id int64) bool {
	queryString := fmt.Sprintf("select id from %s where id = %d", table, id)
	row, err := DBReadRow(queryString)
	if err != nil {
		return false
	}
	var i int
	if err := row.Scan(&i); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false
		}
	}
	return true
}

//-------------------------------SLEEP interfaces ------------------------

//----------------------SLEEP interface-----------------------
type sleepI interface {
	eventBaseWorker
	eventI
	SetId(int64)
	SetEndTime(time.Time)
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
		"WHERE id = %d", s.End().Format("2006-01-02 15:04"), s.Id())
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
	query_string := fmt.Sprintf("select baby_id, start, sleep_end "+
		"from sleep where id = %d", id)
	row, err := DBReadRow(query_string)
	if err != nil {
		return err
	}
	var babyId int64
	var start time.Time
	//var sleepEnd time.Time
	var sleepEnd sql.NullTime
	if err := row.Scan(&babyId, &start, &sleepEnd); err != nil {
		return err
	}
	s.SetId(id)
	s.SetBabyId(babyId)
	s.SetStart(start)
	s.SetEndTime(sleepEnd.Time)

	return nil
}

func (s *sleep) writeStructToBase() error {
	// at first, check existing current event in base.
	// if exists - update, not insert
	existing := s.CheckExisting("sleep", s.id)
	var queryString string
	if existing {
		if s.endTime.IsZero() {
			queryString = fmt.Sprintf("update sleep set (baby_id, start, sleep_end) "+
				"= (%d, '%s') where id = %d",
				s.babyId,
				s.start.Format("2006-01-02 15:04"),
				s.id)
		} else {
			queryString = fmt.Sprintf("update sleep set (baby_id, start, sleep_end) "+
				"= (%d, '%s', '%s') where id = %d",
				s.babyId,
				s.start.Format("2006-01-02 15:04"),
				s.endTime.Format("2006-01-02 15:04"),
				s.id)
		}
	} else {
		if s.endTime.IsZero() {
			queryString = fmt.Sprintf("insert into sleep (baby_id, start) "+
				"values (%d, '%s') RETURNING id",
				s.babyId, s.start.Format("2006-01-02 15:04"))
		} else {
			queryString = fmt.Sprintf("insert into sleep (baby_id, start, sleep_end) "+
				"values (%d, '%s', '%s') RETURNING id",
				s.babyId, s.start.Format("2006-01-02 15:04"), s.endTime.Format("2006-01-02 15:04"))
		}

	}

	fmt.Println(queryString)
	pIdRow, err := DBInsertAndGet(queryString)
	if err != nil {
		if existing && err.Error() == "sql: no rows in result set" {
		} else {
			return err
		}
	}

	if !existing {
		var sleepId int64

		err = pIdRow.Scan(&sleepId)
		if err != nil {
			return err
		}
		s.id = sleepId
	}
	return nil
}

func (s *sleep) SetEndTime(t time.Time) {
	s.endTime = t
	s.calcDuration()
}

func (s *sleep) SetId(id int64) {
	s.id = id
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

//------------------------------PAMPERS INTERFACES------------------------------

type pampersState int64

const (
	wet pampersState = iota
	dirty
	combined
)

//----------------------PAMPERS interface-----------------------
type pampersI interface {
	eventBaseWorker
	eventI
	SetState(pampersState) error
	SetId(int64)

	Id() int64
	State() pampersState
}

type pampers struct {
	event
	id    int64
	state pampersState
}

func newPampers(initEvent event) *pampers {
	p := &pampers{event: initEvent}
	return p
}

func (p *pampers) writeStructToBase() error {
	existing := p.CheckExisting("pampers", p.id)
	var queryString string
	if existing {
		queryString = fmt.Sprintf("update pampers set (baby_id, start, state) "+
			"values(%d, '%s', %d) where id = %d RETURNING id",
			p.babyId, p.start.Format("2006-01-02 15:04"), p.state, p.id)
	} else {

		queryString = fmt.Sprintf("insert into pampers(baby_id, start, state) "+
			"values(%d, '%s', %d) RETURNING id",
			p.babyId, p.start.Format("2006-01-02 15:04"), p.state)
	}
	fmt.Println(queryString)
	pIdRow, err := DBInsertAndGet(queryString)
	if err != nil {
		return err
	}

	var pampersId int64
	if err = pIdRow.Scan(&pampersId); err != nil {
		return err
	}
	p.id = pampersId
	return nil
}

//read sleep by id
func (p *pampers) readStructFromBase(id int64) error {

	query_string := fmt.Sprintf("select baby_id, start, state "+
		"from pampers where id = %d", id)

	row, err := DBReadRow(query_string)
	if err != nil {
		return err
	}
	var babyId int64
	var start time.Time
	var state pampersState
	if err := row.Scan(&babyId, &start, &state); err != nil {
		return err
	}
	p.SetBabyId(babyId)
	p.SetStart(start)
	p.SetState(state)
	p.SetId(id)

	return nil

}

func (p *pampers) SetState(ps pampersState) error {
	if ps > 2 {
		return fmt.Errorf("not valid pampers state")
	}
	p.state = ps
	return nil
}

func (p *pampers) SetId(id int64) {
	p.id = id
}

func (p pampers) State() pampersState {
	return p.state
}

func (p pampers) Id() int64 {
	return p.id
}

//---------------------------------------EAT INTERFACE----------------------------------
type eatI interface {
	eventBaseWorker
	eventI
	SetDescription(string)
	SetId(int64)

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
	existing := e.CheckExisting("eat", e.id)
	var queryString string
	if existing {
		queryString = fmt.Sprintf("update eat set (baby_id, start, description) "+
			"= (%d, '%s', '%s') where id = %d RETURNING id",
			e.babyId, e.start.Format("2006-01-02 15:04"), e.description, e.id)
	} else {

		queryString = fmt.Sprintf("insert into eat(baby_id, start, description) "+
			"VALUES (%d, '%s', '%s') RETURNING id;",
			e.event.babyId, e.start.Format("2006-01-02 15:04"), e.description)
	}
	row, err := DBInsertAndGet(queryString)
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
	e.SetId(id)

	return nil
}

func (e *eat) SetId(id int64) {
	e.id = id
}
