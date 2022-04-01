package main

import (
	"database/sql"
)

var dbInfo = "host=127.0.0.1 port=5432 user=postgres password=!QAZxsw2 dbname=test_db"

type eventBaseWorker interface {
	writeStructToBase() error
	readStructFromBase(interface{}) error
}

func DBInsertAndGet(query_string string) (*sql.Row, error) {
	db, err := sql.Open("pgx", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result := db.QueryRow(query_string)
	return result, nil
	// TODO: case data inserged, but id not received

}

func DBReadRow(queryString string) (*sql.Row, error) {
	db, err := sql.Open("pgx", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result := db.QueryRow(queryString)
	return result, nil
}

func DBReadRows(queryString string) (*sql.Rows, error) {
	db, err := sql.Open("pgx", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result, err := db.Query(queryString)
	if err != nil {
		return nil, err
	}
	return result, nil

}
