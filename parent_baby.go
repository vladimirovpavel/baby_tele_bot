package main

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

//меняем ID на телеграм ID. Удаляем phone, добавляем current_baby

type parentI interface {
	eventBaseWorker

	SetId(id int64)
	SetName(name string)
	SetCurrentBaby(babyId int64) error

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
		return fmt.Errorf("error on write parent %d to base:\n%s", p.id, err.Error())
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			queryString = fmt.Sprintf("insert into parent(parent_id, name, current_baby) "+
				"values('%d', '%s', '%d') RETURNING parent_id", p.Id(), p.Name(), p.currentBaby)
		} else {
			return fmt.Errorf("error on check existing parent %d in base  to base:\n%s",
				p.id, err.Error())
		}
	} else {
		queryString = fmt.Sprintf("update parent set (name, current_baby) "+
			"= ('%s', %d) where parent_id = %d", p.name, p.currentBaby, p.id)
	}

	_, err = DBInsertAndGet(queryString)
	if err != nil {
		return fmt.Errorf("error on insert on update or create parent %d:\n%s", p.Id(), err.Error())
	}

	return nil
}

//read parent by id
func (p *parent) readStructFromBase(id int64) error {
	queryString := fmt.Sprintf("select name, current_baby from parent where parent_id=%d", id)

	row, err := DBReadRow(queryString)
	if err != nil {
		return fmt.Errorf("error on read data for parent %d:\n%s", p.id, err.Error())
	}

	var name string
	var currentBaby int64
	if err := row.Scan(&name, &currentBaby); err != nil {
		return fmt.Errorf("error on receive data for parent %d:\n%s", p.id, err.Error())
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

func (p *parent) SetCurrentBaby(id int64) error {
	if id != 0 {
		babyes, err := GetBabyesByParent(p.id)
		if err != nil {
			return fmt.Errorf("error on get babyes by parent %d:\n%s", p.id, err.Error())
		}
		founded := false
		for _, b := range babyes {
			if b.Id() == id {
				founded = true
				break
			}
		}
		if !founded {
			return fmt.Errorf("error, not found baby with id %d", p.id)
		}
	}
	p.currentBaby = id

	return nil
}

type babyI interface {
	eventBaseWorker
	getState() string
	SetId(int64)
	SetName(string)
	SetParent(int64) error
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
	queryString := fmt.Sprintf("select baby_id from baby where baby_id=%d", b.id)

	row, err := DBReadRow(queryString)
	if err != nil {
		return err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			queryString = fmt.Sprintf("insert into baby(parent_id, name, birth) "+
				"values (%d, '%s', '%s') RETURNING baby_id",
				b.parentId, b.name, b.birth.Format("2006-01-02"))
		} else {
			return err
		}
	} else {
		queryString = fmt.Sprintf("update baby set (parent_id, name, birth) = "+
			" (%d, '%s', '%s') where baby_id = %d RETURNING baby_id",
			b.parentId, b.name, b.birth.Format("2006-01-02"), b.id)
	}
	bIdRow, err := DBInsertAndGet(queryString)
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

func (b *baby) SetParent(parentId int64) error {
	if !CheckParentRegistred(int64(parentId)) {
		return fmt.Errorf("error, parent with id not existing")
	}
	b.parentId = parentId
	return nil
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
	return fmt.Sprintf("%s %s", b.Name(), b.Birth().Format("2006-01-02"))
}
