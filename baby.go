package main

import (
	"fmt"
	"time"
)

type babyI interface {
	getState() string

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

func newBaby() baby {
	b := baby{
		id:       0,
		parentId: 0,
		name:     "",
		birth:    time.Now(),
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

//read baby by parent
func (b *baby) readStructFromBase(query interface{}) error {
	parentId := query.(int64)

	var babyid int64
	var name string
	var birth time.Time

	query_string := fmt.Sprintf("select baby_id, name, birth from baby where parent_id=%d", parentId)

	row, err := DBReadRow(query_string)
	if err != nil {
		return err
	}

	if err := row.Scan(&babyid, name, birth); err != nil {
		return err
	}

	b.parentId = parentId
	b.id = babyid
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

func (b baby) String() string {
	return fmt.Sprintf("%s %s", b.Name(), b.Birth())
}
