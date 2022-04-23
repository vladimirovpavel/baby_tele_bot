package main

import (
	"database/sql"
	"fmt"
	"time"
)

//переделал конструктор событий
// методы GetTypedEventsIdsByBabyDate GetEventsByBabyDate GetNotEndedSleepForBaby
// removed testdata, crates event_test
var dbInfo = "host=127.0.0.1 port=5432 user=postgres password=!QAZxsw2 dbname=test_db"

// TODO: стоит подумать о том, чтобы сделать одну большую таблицу EVENTS
// из которой будут идти ссылки на таблицы sleeps, eats, pampers.
// тогда в любом случае будем делать один запрос конкретного eventа, и уже оптом
// добирать из базы "дополнительные" данные для него.
// сейчас же мы делаем два запроса из одной из той же таблицы - сначала заполняя
// структуру event, уже потом заполняя дополнительные данные
// TODO: нужно ли, может, удалить??? или наоборот - внедрить всюду?
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

func DBDeleteData(queryString string) error {
	db, err := sql.Open("pgx", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	result := db.QueryRow(queryString)
	return result.Err()
}

func DBInsertAndGet(queryString string) (*sql.Row, error) {
	fmt.Printf("Inserting to DB wiht query: \n%s\n", queryString)
	db, err := sql.Open("pgx", dbInfo)

	if err != nil {
		return nil, err
	}
	defer db.Close()
	result := db.QueryRow(queryString)
	return result, result.Err()
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

func RegisterNewParent(parentId int64, parentName string) (*parent, error) {
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

func GetTypedEventsIdsByBabyDate(babyId int64, t time.Time, tableName string) ([]int64, error) {
	var results []int64
	queryString := fmt.Sprintf("select id from %s where baby_id = %d and start > '%s'",
		tableName, babyId, t.Format("2006-01-02"))

	rows, err := DBReadRows(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventId int64
	for rows.Next() {
		if err := rows.Scan(&eventId); err != nil {
			break
		}
		results = append(results, eventId)
	}

	return results, nil
}

func GetEventsByBabyDate(babyId int64, t time.Time) []eventI {
	var resultTables []eventI
	eventsTables := []string{"sleep", "eat", "pampers"}
	for _, table := range eventsTables {
		idsList, err := GetTypedEventsIdsByBabyDate(babyId, t, table)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, eventId := range idsList {

			nEvent := newEvent(babyId)

			// TODO: вот прямо испытываю ощущение что ту можно сделать гораздо лучше
			// TODO: пожалуй, можно вообще не слайс объектов возвращать, а тупо слайс строк
			// тогда этот код совсем упростится
			switch table {
			case "sleep":
				{
					nSleep := newSleep(*nEvent)
					if err := nSleep.readStructFromBase(eventId); err != nil {
						fmt.Println(err)
						continue
					}
					resultTables = append(resultTables, nSleep)
				}
			case "eat":
				{
					nEat := newEat(*nEvent)
					if err := nEat.readStructFromBase(eventId); err != nil {
						fmt.Println(err)
						continue
					}
					resultTables = append(resultTables, nEat)

				}
			case "pampers":
				{
					nPampers := newPampers(*nEvent)
					if err := nPampers.readStructFromBase(eventId); err != nil {
						fmt.Println(err)
						continue
					}
					resultTables = append(resultTables, nPampers)
				}

				// хотим для id каждого типа создать объект соотв типа и добавить
				// в слайс eventI

			}
		}

	}
	return resultTables
}

func GetNotEndedSleepForBaby(babyId int64) (sleepI, error) {
	queryString := fmt.Sprintf("select id from sleep where (sleep_end is null) AND (baby_id = %d);",
		babyId)
	row, err := DBReadRow(queryString)
	if err != nil {
		return nil, err
	}

	var notEndedSleepId int64
	err = row.Scan(&notEndedSleepId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	var s = newSleep(*newEvent(babyId))
	// TODO: !!!!!!!!!!!!!!!! проблема в том, что при чтении свежесозданного события
	// TODO: не читается null конец события
	if err := s.readStructFromBase(notEndedSleepId); err != nil {
		return nil, err
	}

	return s, nil
}

func RemoveBabyFromBase(babyId int64) error {
	queryString := fmt.Sprintf("delete from baby where baby_id = %d", babyId)
	_, err := DBReadRow(queryString)
	if err != nil {
		return err
	}
	return nil
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

func GetCurrentBaby(parentId int64) (babyI, error) {
	queryString := fmt.Sprintf("select current_baby from parent where parent_id = %d", parentId)
	row, err := DBReadRow(queryString)
	if err != nil {
		return nil, err
	}

	var babyId int64

	if err := row.Scan(&babyId); err != nil || babyId == 0 {
		return nil, err
	}
	b := newBaby()
	if err := b.readStructFromBase(babyId); err != nil {
		fmt.Printf("Error reading baby for parent %d:\n%s", parentId, err)
	}

	return b, nil
}

func GetParentsIds() ([]int64, error) {
	queryString := "select parent_id from parent"
	rows, err := DBReadRows(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	var id int64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			break
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func CheckParentRegistred(parentId int64) bool {
	ids, err := GetParentsIds()
	if err != nil {
		return false
	}

	founded := false
	for _, i := range ids {
		if int64(parentId) == i {
			founded = true
		}
	}
	return founded

}
