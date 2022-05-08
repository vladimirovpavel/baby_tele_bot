package main

import (
	"database/sql"
	"fmt"
	"time"
)

var dbInfo string

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
	slogger.Debugf("Inserting to DB wiht query: \n%s\n", queryString)
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
	slogger.Debugf("Reading Row from DB with query:\n%s\n", queryString)
	db, err := sql.Open("pgx", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result := db.QueryRow(queryString)
	return result, result.Err()
}

//Read ROWS from db. DO NOT FORGET TO CALL rows.close()!!!
func DBReadRows(queryString string) (*sql.Rows, error) {
	slogger.Debugf("Reading RowS from DB with query:\n%s\n", queryString)
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

// функция проверит существование родителя в базе.
// Если он не существует - создаст его
// если сущесвует - получит данные из базы и вернет экземпляр
func RegisterNewParent(parentId int64, parentName string) (*parent, error) {
	queryString := fmt.Sprintf("select parent_id from parent where parent_id = %d", parentId)
	row, err := DBReadRow(queryString)
	if err != nil {
		return nil, fmt.Errorf("error reading parent id %d from base:\n%s", parentId, err)
	}
	var pId int64
	var p *parent

	err = row.Scan(&pId)

	if err != nil && err.Error() != "sql: no rows in result set" {
		// ошибка чтения данных из базы, не "НЕТ РЕЛЕВАНТНЫХ СТРОЧЕК"
		return nil, fmt.Errorf("error reading parent id %d from base:\n%s", parentId, err)
	} else if err == nil { // ошибки нет -> запрос успешен ->  родитель уже в базе
		slogger.Debugf("Parent with id %d and name %s already in base\n", parentId, parentName)
		p = newParent()
		err := p.readStructFromBase(parentId)
		if err != nil {
			return nil, fmt.Errorf("error reading parent struct for parent id %d from base:\n%s", parentId, err)
		}
	} else {
		//if error == "sql: no rows in result set" - записи в базе нет, родитель
		//не зарегистрирован - можно регистрировать
		p = newParent()
		p.SetName(parentName)
		p.SetId(parentId)
		if err := p.writeStructToBase(); err != nil {
			return nil, fmt.Errorf("error writing parent struct for id %d  to base:\n%s", parentId, err)
		}
		slogger.Debugf("Parent with id %d and name %s added to db\n", parentId, parentName)

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

func GetCurrentBaby(parentId int64) (babyI, error) {
	queryString := fmt.Sprintf("select current_baby from parent where parent_id = %d", parentId)
	row, err := DBReadRow(queryString)
	if err != nil {
		return nil, fmt.Errorf("error get current baby for parent %d\n%s", parentId, err)
	}

	var babyId int64

	if err := row.Scan(&babyId); err != nil || babyId == 0 {
		return nil, err
	}
	b := newBaby()
	if err := b.readStructFromBase(babyId); err != nil {
		return nil, fmt.Errorf("error reading baby for parent %d:\n%s", parentId, err)
	}

	return b, nil
}

func GetParentsIds() ([]int64, error) {
	queryString := "select parent_id from parent"
	rows, err := DBReadRows(queryString)
	if err != nil {
		return nil, fmt.Errorf("error on get all parents ids:\n%s", err.Error())
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

func RemoveBabyFromBase(babyId int64) error {
	queryString := fmt.Sprintf("delete from baby where baby_id = %d", babyId)
	_, err := DBReadRow(queryString)
	if err != nil {
		return fmt.Errorf("error on remove baby with id %d from base:\n%s", babyId, err.Error())
	}
	return nil
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

// returns events IDs
func GetTypedEventsIdsByBabyDate(babyId int64, t time.Time, tableName string) ([]int64, error) {
	var results []int64
	queryString := fmt.Sprintf("select id from %s where baby_id = %d and start > '%s' "+
		"and start < '%s'",
		tableName, babyId, t.Format("2006-01-02"), (t.AddDate(0, 0, 1)).Format("2006-01-02"))

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
			slogger.Debugw("On GetEventsByBabyDate has error",
				"error", err.Error())
			continue
		}
		for _, eventId := range idsList {

			nEvent := newEvent(babyId)
			switch table {
			case "sleep":
				{
					nSleep := newSleep(*nEvent)
					if err := nSleep.readStructFromBase(eventId); err != nil {
						slogger.Debugw("On GetEventsByBabyDate has error",
							"error", err.Error())
						continue
					}
					resultTables = append(resultTables, nSleep)
				}
			case "eat":
				{
					nEat := newEat(*nEvent)
					if err := nEat.readStructFromBase(eventId); err != nil {
						slogger.Debugw("On GetEventsByBabyDate has error",
							"error", err.Error())
						continue
					}
					resultTables = append(resultTables, nEat)

				}
			case "pampers":
				{
					nPampers := newPampers(*nEvent)
					if err := nPampers.readStructFromBase(eventId); err != nil {
						slogger.Debugw("On GetEventsByBabyDate has error",
							"error", err.Error())
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
		return nil, fmt.Errorf("error get not ended sleep id for baby %d:\n%s", babyId, err.Error())
	}

	var notEndedSleepId int64
	err = row.Scan(&notEndedSleepId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		} else {
			return nil, fmt.Errorf("error on receive not ended sleep for baby %d:\n%s", babyId, err.Error())
		}
	}
	var s = newSleep(*newEvent(babyId))
	if err := s.readStructFromBase(notEndedSleepId); err != nil {
		return nil, fmt.Errorf("error on read struct from base for sleep id %d\n%s", notEndedSleepId, err.Error())
	}

	return s, nil
}
