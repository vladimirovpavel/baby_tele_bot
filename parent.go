package main

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type parentI interface {
	// TODO: what about set funcs&

	Id()
	Phone()
	Name()
}

type parent struct {
	id    int64
	name  string
	phone string
	// #TODO: store phone in DB as 10 digits for Russia
}

func newParent() parent {
	p := parent{}
	return p
}

func (p *parent) writeStructToBase() error {
	query := fmt.Sprintf("insert into parent(name, phone) "+
		"values('%s', '%s') RETURNING parent_id", p.Name(), p.Phone())

	pIdRow, err := DBInsertAndGet(query)
	if err != nil {
		return err
	}

	var parentId int64
	if err = pIdRow.Scan(&parentId); err != nil {
		return err
	}
	p.id = parentId
	return nil
}

//read parent by phone
func (p *parent) readStructFromBase(query interface{}) error {
	var parentId int64
	var name string

	parentPhone := query.(string)
	query_string := fmt.Sprintf("select parent_id, name from parent where phone=%s", parentPhone)

	row, err := DBReadRow(query_string)
	if err != nil {
		return err
	}
	if err := row.Scan(&parentId, &name); err != nil {
		return err
	}

	p.id = parentId
	p.phone = parentPhone
	p.name = name

	return nil
}

func (p parent) Id() int64 {
	return p.id
}
func (p parent) Phone() string {
	return p.phone
}

func (p parent) Name() string {
	return p.name
}
