package main

/*
TODO:
1. реализовать Write struct to base для всех событий
2. реализовать update для sleep события
3. реализовать выборку по дате и ребенку для всех событий
	3.1. тут можно сделать максимальную выборку select * from ... в event, дергать ее
	3.2. а потом уже парсить результат в каждом родителе в зависимости от типа
4. а что у нас вообще с получением из базы? по каким ключам мы будем получать?
*/

func main() {
	createTestData()
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
