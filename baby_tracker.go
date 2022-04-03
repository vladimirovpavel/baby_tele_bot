package main

import (
	"fmt"
	"time"
)

/*
TODO:
1. OK реализовать Write struct to base для всех событий
2. OK реализовать update для sleep события
3. реализовать выборку по дате и ребенку для всех событий
	3.1. тут можно сделать максимальную выборку select * from ... в event, дергать ее
	3.1.1. потестить написанную selectbybabydate
	3.2. а потом уже парсить результат в каждом родителе в зависимости от типа
5. OK duration не нужен в sleep
6. OK потестить update для sleep
7. OK event init должен полуать не baby, а его id.
8. для всех структур реализовать метод stirng для выдачи клиенту
9. для событий изменил SLEEP_ID на просто id. проверить что ничего не сломалось
10. для запроса события понадобится второй метод - можно в eventI - readEventSpecial, который будет считывать все специфичное для конкретного эвента

*/

func main() {
	start := time.Date(1975, 2, 17, 0, 0, 0, 0, time.UTC)
	query := DBEventQuery{table: "eat", babyId: 1, start: start}
	res, _ := SelectEventByBabyDate(query)
	fmt.Printf("%s", res[0])
	//a, b, c, d, e := createTestData()
	//fmt.Println(a, b, c, d, e)
	/* p := newParent()
	b := newBaby()
	b.SetParent(p.Id())
	b.SetId(1)

	e := newEvent(b.Id(), time.Now())
	s := newSleep(e)
	if err := s.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
	s.setEndTime(time.Now())
	if err := s.updateEndSleepTime(); err != nil {
		fmt.Println(err)
	} */
	//e[0].SelectByBabyDate(map[string]string{"table": "sleep_id", "babyId": "1", "start": "1975-02-18"})

	/* time.Sleep(1 * time.Second)
	e := newEvent(b, time.Now())

	time.Sleep(2 * time.Second)
	s := newSleep(e)

	time.Sleep(1 * time.Second)
	s.setEndTime(time.Now())

	fmt.Printf("Sleep starts at: %s\n\tends at: %s\n\tduration is: %s",
		s.Start(), s.End(), s.Duration())

	pe := newEvent(b, time.Now())
	p := newPampers(pe)
	fmt.Printf("\n\nPampers time is: %s", p.Start()) */
}
