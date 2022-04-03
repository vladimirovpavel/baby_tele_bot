package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var dbInfo = "host=127.0.0.1 port=5432 user=postgres password=!QAZxsw2 dbname=test_db"

type DBEventQuery struct {
	table  string
	babyId int64
	start  time.Time
}

func (eventQuery DBEventQuery) String() string {
	return fmt.Sprintf("Event of type %s baby_id is %d, started %s",
		eventQuery.table, eventQuery.babyId, eventQuery.start.Format("2006-01-02"))
}

type eventBaseWorker interface {
	writeStructToBase() error
	readStructFromBase(interface{}) error
}

func DBInsertAndGet(queryString string) (*sql.Row, error) {
	fmt.Printf("Inserting to DB wiht query: \n%s\n", queryString)
	db, err := sql.Open("pgx", dbInfo)

	if err != nil {
		return nil, err
	}
	defer db.Close()
	result := db.QueryRow(queryString)
	return result, nil
	// TODO: case data inserged, but id not received

}

func DBReadRow(queryString string) (*sql.Row, error) {
	fmt.Printf("Reading Row from DB with query:\n%s\n", queryString)
	db, err := sql.Open("pgx", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result := db.QueryRow(queryString)
	return result, nil
}

//Read ROWS from db. DO NOT FORGET TO CALL rows.close()!!!
func DBReadRows(queryString string) (*sql.Rows, error) {
	fmt.Printf("Reading RowS from DB with query:\n%s\n", queryString)
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

//receive eventI from table of event type. Then, receive
func SelectEventByBabyDate(queryData DBEventQuery) ([]eventI, error) {
	var results []eventI

	queryString := fmt.Sprintf("select id, baby_id, start from %s where baby_id=%d and start > '%s';",
		queryData.table, queryData.babyId, queryData.start.Format("2006-01-02"))
	rows, err := DBReadRows(queryString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventId int64
		var babyId int64
		var start string
		if err := rows.Scan(&eventId, &babyId, &start); err != nil {
			fmt.Println(err)
			break
		}
		splittedDate := strings.Split(start, "T")
		startTime, err := time.Parse("2006-01-02", splittedDate[0])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		currentEvent := newEvent(babyId, startTime)

		results = append(results, currentEvent)
	}

	return results, err

}
