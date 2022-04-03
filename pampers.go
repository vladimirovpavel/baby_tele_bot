package main

import "fmt"

type pampersState int64

const (
	wet pampersState = iota
	dirty
	combined
)

//----------------------PAMPERS interface-----------------------
type pampersI interface {
	eventI
	SetState(pampersState)

	Id() int64
	State() pampersState
}

type pampers struct {
	event
	id    int64
	state pampersState
}

func newPampers(initEvent event) pampers {
	p := pampers{event: initEvent}
	return p
}

func (p *pampers) writeStructToBase() error {
	query_string := fmt.Sprintf("insert into pampers(baby_id, start, state) "+
		"values(%d, '%s', %d) RETURNING pampers_id", p.BabyId(), p.Start().Format("2006-01-02"), p.State())

	pIdRow, err := DBInsertAndGet(query_string)
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

func (p *pampers) SetState(ps pampersState) {
	p.state = ps
}

func (p pampers) State() pampersState {
	return p.state
}

func (p pampers) Id() int64 {
	return p.id
}
