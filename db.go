package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var dbInfo = "host=127.0.0.1 port=5432 user=postgres password=!QAZxsw2 dbname=test_db"

// TODO: стоит подумать о том, чтобы сделать одну большую таблицу EVENTS
// из которой будут идти ссылки на таблицы sleeps, eats, pampers.
// тогда в любом случае будем делать один запрос конкретного eventа, и уже оптом
// добирать из базы "дополнительные" данные для него.
// сейчас же мы делаем два запроса из одной из той же таблицы - сначала заполняя
// структуру event, уже потом заполняя дополнительные данные
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
	readStructFromBase(id int64) error
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

// Регистрация нового родителя. проверяем - есть ли уже в базе запись о родителе,
// если есть - вернем уже существующую структуру
// если нет - создади запись и вернем структуру
func RegisterNewParent(parentId int64, parentName string) (*parent, error) {
	//проверяем - не существует ли родитель в настоящий момент
	queryString := fmt.Sprintf("select parent_id from parent where parent_id = %d", parentId)
	row, err := DBReadRow(queryString)
	if err != nil {
		return nil, err
	}
	var pId int64
	var p *parent

	err = row.Scan(&pId)

	if err != nil && err.Error() != "sql: no rows in result set" {
		// ошибка чтения данных из базы, не "НЕТ РЕЛЕВАНТНЫХ СТРОЧЕК"
		return nil, err
	} else if err == nil { // ошибки нет -> запрос успешен ->  родитель уже в базе
		fmt.Printf("Parent with id %d and name %s already in base\n", parentId, parentName)
		p = newParent()
		err := p.readStructFromBase(parentId)
		if err != nil {
			return nil, err
		}
	} else {
		//if error == "sql: no rows in result set" - записи в базе нет, родитель
		//не зарегистрирован - можно регистрировать
		p = newParent()
		p.SetName(parentName)
		p.SetId(parentId)
		if err := p.writeStructToBase(); err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Printf("Parent with id %d and name %s added to db\n", parentId, parentName)

	}
	return p, nil

}

func GetBabyesByParent(parentId int64) ([]babyI, error) {
	queryString := fmt.Sprintf("select baby_id from baby where parent_id = %d", parentId)

	rows, err := DBReadRows(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var babyes []babyI
	var babyId int64
	err = nil
	for rows.Next() {
		if err := rows.Scan(&babyId); err != nil {
			break
		}
		baby := newBaby()
		if err := baby.readStructFromBase(babyId); err != nil {
			fmt.Println(err)
			continue
		}
		babyes = append(babyes, baby)
	}
	return babyes, err
}

func GetBabyesEventsByDate(id int64, t time.Time) ([]eventI, error) {
	/* query_string := fmt.Sprintf("select (baby_id, start, sleep_end) "+
	"from sleep where sleep_start > '%s' and sleep_start < ('%s' + '1 day'::interval",
	date, date) */
	return nil, nil
}
