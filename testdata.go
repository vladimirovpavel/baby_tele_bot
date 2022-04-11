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

//func createParents()

func CreatingData() {
	var parents []parentI
	var babyes []babyI
	var eats []eatI
	var pamperses []pampersI
	var sleeps []sleepI

	var p *parent
	var b *baby
	var e *event
	var ea *eat
	var sl *sleep
	var pam *pampers

	// TODO: Create many events func!!!
	p = newParent()
	p.SetName("Elisa")
	if err := p.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	parents = append(parents, p)

	p = newParent()
	p.SetName("John")
	if err := p.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	parents = append(parents, p)

	b = newBaby()
	b.SetBirth(getDataForTest())
	b.SetName("Mary")

	b.SetParent(parents[0].Id())
	if err := b.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	babyes = append(babyes, b)

	b = newBaby()
	b.SetBirth(getDataForTest())
	b.SetName("Bobby")
	b.SetParent(parents[0].Id())
	if err := b.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	babyes = append(babyes, b)

	b = newBaby()
	b.SetBirth(getDataForTest())
	b.SetName("Roger")
	b.SetParent(parents[1].Id())
	if err := b.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	babyes = append(babyes, b)

	e = newEvent(babyes[0].Id(), getDataForTest())
	pam = newPampers(*e)
	pam.SetState(wet)
	if err := pam.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	pamperses = append(pamperses, pam)

	e = newEvent(babyes[1].Id(), getDataForTest())
	pam = newPampers(*e)
	pam.SetState(dirty)
	if err := pam.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	pamperses = append(pamperses, pam)

	e = newEvent(babyes[2].Id(), getDataForTest())
	pam = newPampers(*e)
	pam.SetState(combined)
	if err := pam.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	pamperses = append(pamperses, pam)

	e = newEvent(babyes[2].Id(), getDataForTest())
	ea = newEat(*e)
	ea.SetDescription("50g porrige")
	if err := ea.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	eats = append(eats, ea)

	e = newEvent(babyes[0].Id(), getDataForTest())
	ea = newEat(*e)
	ea.SetDescription("100g spahetty")
	if err := ea.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	eats = append(eats, ea)

	e = newEvent(babyes[1].Id(), getDataForTest())
	sl = newSleep(*e)
	sl.setEndTime(getDataForTest())
	if err := sl.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sl)

	e = newEvent(babyes[1].Id(), getDataForTest())
	sl = newSleep(*e)
	sl.setEndTime(getDataForTest())
	if err := sl.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sl)

	e = newEvent(babyes[1].Id(), getDataForTest())
	sl = newSleep(*e)
	sl.setEndTime(getDataForTest())
	if err := sl.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sl)

	sleeps[0].setEndTime(getDataForTest())
	if err := sleeps[0].updateEndSleepTime(); err != nil {
		fmt.Println(err)
	}

	e = newEvent(babyes[1].Id(), getDataForTest())
	sl = newSleep(*e)
	sl.setEndTime(getDataForTest())
	if err := sl.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sl)

	e = newEvent(babyes[1].Id(), getDataForTest())
	sl = newSleep(*e)
	sl.setEndTime(getDataForTest())
	if err := sl.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sl)

	e = newEvent(babyes[1].Id(), getDataForTest())
	sl = newSleep(*e)
	sl.setEndTime(getDataForTest())
	if err := sl.writeStructToBase(); err != nil {
		fmt.Println(err)
	}
	sleeps = append(sleeps, sl)

	sleeps[3].setEndTime(getDataForTest())
	if err := sleeps[3].updateEndSleepTime(); err != nil {
		fmt.Println(err)
	}

	sleeps[1].setEndTime(getDataForTest())
	if err := sleeps[1].updateEndSleepTime(); err != nil {
		fmt.Println(err)
	}
	sleeps[2].setEndTime(getDataForTest())
	if err := sleeps[2].updateEndSleepTime(); err != nil {
		fmt.Println(err)
	}

	//require.True(t, true)
	//return parents, babyes, pamperses, eats, sleeps

}
