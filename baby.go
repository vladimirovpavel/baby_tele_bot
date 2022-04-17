package main

import (
	"fmt"
	"time"
)

type babyI interface {
	eventBaseWorker
	getState() string
	SetId(int64)
	SetName(string)
	SetParent(int64)
	SetBirth(time.Time)

	Id() int64
	ParentId() int64
	Name() string
	Birth() time.Time
}

type baby struct {
	id       int64
	parentId int64
	name     string
	birth    time.Time
}

func newBaby() *baby {
	b := &baby{
		id:       0,
		parentId: 0,
		name:     "",
		birth:    time.Time{},
	}
	return b
}

func (b *baby) writeStructToBase() error {
	query_string := fmt.Sprintf("insert into baby(parent_id, name, birth) "+
		"values (%d, '%s', '%s') RETURNING baby_id",
		b.ParentId(), b.Name(), b.Birth().Format("2006-01-02"))
	bIdRow, err := DBInsertAndGet(query_string)
	if err != nil {
		return err
	}

	var babyId int64
	if err = bIdRow.Scan(&babyId); err != nil {
		return err
	}
	b.id = babyId
	return nil
}

//read baby by babyid
func (b *baby) readStructFromBase(id int64) error {
	queryString := fmt.Sprintf("select parent_id, name, birth from baby where baby_id=%d", id)

	row, err := DBReadRow(queryString)
	if err != nil {
		return err
	}
	var parentId int64
	var name string
	var birth time.Time

	if err := row.Scan(&parentId, &name, &birth); err != nil {
		return err
	}

	b.id = id
	b.parentId = parentId
	b.name = name
	b.birth = birth
	return nil
}

func (b baby) getState() string {
	stateString := ""
	return stateString
}

func (b baby) Id() int64 {
	return b.id
}

func (b baby) ParentId() int64 {
	return b.parentId
}

func (b baby) Name() string {
	return b.name
}
func (b baby) Birth() time.Time {
	return b.birth
}

func (b *baby) SetParent(parentId int64) {
	b.parentId = parentId
}

func (b *baby) SetName(name string) {
	b.name = name
}

func (b *baby) SetBirth(birth time.Time) {
	b.birth = birth
}

func (b *baby) SetId(id int64) {
	b.id = id
}

func (b baby) String() string {
	return fmt.Sprintf("%s %s", b.Name(), b.Birth())
}
