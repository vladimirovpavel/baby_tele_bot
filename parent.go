package main

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

//меняем ID на телеграм ID. Удаляем phone, добавляем current_baby

//TODO: current baby
type parentI interface {
	eventBaseWorker

	SetId(id int64)
	SetName(name string)
	SetCurrentBaby(babyId int64) //TODO: implement this

	Id() int64
	Name() string
	CurrentBaby() int64
}

type parent struct {
	id          int64
	name        string
	currentBaby int64
}

func newParent() *parent {
	p := &parent{}
	return p
}

func (p *parent) writeStructToBase() error {
	// checks existing parent in base. If exists - not INSERT, UPDATE
	queryString := fmt.Sprintf("select parent_id from parent where parent_id=%d", p.id)

	row, err := DBReadRow(queryString)
	if err != nil {
		return err
	}

	var id string
	if err := row.Scan(&id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			queryString = fmt.Sprintf("insert into parent(parent_id, name, current_baby) "+
				"values('%d', '%s', '%d') RETURNING parent_id", p.Id(), p.Name(), p.currentBaby)
		} else {
			return err
		}
	} else {
		queryString = fmt.Sprintf("update parent set (name, current_baby) "+
			"= ('%s', %d) where parent_id = %d", p.Name(), p.CurrentBaby(), p.Id())
	}

	//pIdRow, err := DBInsertAndGet(query)
	_, err = DBInsertAndGet(queryString)
	if err != nil {
		return err
	}

	/* var parentId int64
	if err = pIdRow.Scan(&parentId); err != nil {
		return err
	}
	p.id = parentId */
	return nil
}

//read parent by id
func (p *parent) readStructFromBase(id int64) error {
	queryString := fmt.Sprintf("select name, current_baby from parent where parent_id=%d", id)

	row, err := DBReadRow(queryString)
	if err != nil {
		return err
	}

	var name string
	var currentBaby int64
	if err := row.Scan(&name, &currentBaby); err != nil {
		return err
	}

	p.SetId(id)
	p.SetName(name)
	p.SetCurrentBaby(currentBaby)

	return nil
}

func (p parent) Id() int64 {
	return p.id
}

func (p parent) Name() string {
	return p.name
}

func (p parent) CurrentBaby() int64 {
	return p.currentBaby
}

func (p *parent) SetId(id int64) {
	p.id = id
}
func (p *parent) SetName(name string) {
	p.name = name
}

func (p *parent) SetCurrentBaby(id int64) {
	p.currentBaby = id
}
