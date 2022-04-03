package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getDataForTest() time.Time {

	min := time.Date(1975, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(1975, 3, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min

	return time.Unix(sec, 0)

}

func createTestData() ([]parent, []baby, []pampers, []eat, []sleep) {
	var parents []parent
	var babyes []baby
	var eats []eat
	var pamperses []pampers
	var sleeps []sleep

	// TODO: Create many events func!!!
	p := newParent()
	p.phone = "892943253"
	p.name = "Elisa"
	if err := p.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	parents = append(parents, p)

	p = newParent()
	p.phone = "23453251"
	p.name = "John"
	if err := p.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	parents = append(parents, p)

	b := newBaby()
	b.birth = getDataForTest()
	b.name = "Mary"
	b.parentId = parents[0].Id()
	if err := b.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	babyes = append(babyes, b)

	b = newBaby()
	b.birth = getDataForTest()
	b.name = "Bobby"
	b.parentId = parents[0].Id()
	if err := b.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	babyes = append(babyes, b)

	b = newBaby()
	b.birth = getDataForTest()
	b.name = "Roger"
	b.parentId = parents[0].Id()
	if err := b.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	babyes = append(babyes, b)

	e := newEvent(babyes[0].Id(), getDataForTest())
	pamp := newPampers(*e)
	pamp.SetState(wet)
	if err := pamp.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	pamperses = append(pamperses, pamp)

	e = newEvent(babyes[1].Id(), getDataForTest())
	pamp = newPampers(*e)
	pamp.SetState(dirty)
	if err := pamp.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	pamperses = append(pamperses, pamp)

	e = newEvent(babyes[2].Id(), getDataForTest())
	pamp = newPampers(*e)
	pamp.SetState(combined)
	if err := pamp.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	pamperses = append(pamperses, pamp)

	e = newEvent(babyes[2].Id(), getDataForTest())
	eat := newEat(*e)
	eat.SetDescription("50g porrige")
	if err := eat.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	eats = append(eats, eat)

	e = newEvent(babyes[0].Id(), getDataForTest())
	eat = newEat(*e)
	eat.SetDescription("100g spahetty")
	if err := eat.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	eats = append(eats, eat)

	e = newEvent(babyes[1].Id(), getDataForTest())
	sleep := newSleep(*e)
	sleep.setEndTime(getDataForTest())
	if err := sleep.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sleep)

	e = newEvent(babyes[1].Id(), getDataForTest())
	sleep = newSleep(*e)
	sleep.setEndTime(getDataForTest())
	if err := sleep.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sleep)
	e = newEvent(babyes[1].Id(), getDataForTest())
	sleep = newSleep(*e)
	sleep.setEndTime(getDataForTest())
	if err := sleep.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sleep)

	sleeps[0].setEndTime(getDataForTest())
	if err := sleeps[0].updateEndSleepTime(); err != nil {
		fmt.Println(err)
	}

	return parents, babyes, pamperses, eats, sleeps

}
