package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func getDataForTest() time.Time {

	min := time.Date(1975, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(1975, 3, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min

	return time.Unix(sec, 0)

}

func createTestData() {
	var parents []parentI
	p1 := newParent()
	p1.SetName("TestParent1")
	p1.SetId(1)

	p2 := newParent()
	p2.SetName("TestParent2")
	p2.SetId(2)

	parents = append(parents, p1, p2)

	var babyes []babyI
	b1 := newBaby()
	b1.SetBirth(getDataForTest())
	b1.SetName("testBaby1")
	b1.SetParent(p1.Id())

	b2 := newBaby()
	b2.SetBirth(getDataForTest())
	b2.SetName("testBaby1")
	b2.SetParent(p1.Id())

	b3 := newBaby()
	b3.SetBirth(getDataForTest())
	b3.SetName("testBaby1")
	b3.SetParent(p2.Id())

	p1.SetCurrentBaby(b1.Id())
	p2.SetCurrentBaby(b2.Id())
	babyes = append(babyes, b1, b2, b3)

	event1 := newEvent(b1.Id())
	event2 := newEvent(b1.Id())
	event3 := newEvent(b2.Id())
	event4 := newEvent(b2.Id())
	event5 := newEvent(b2.Id())
	event6 := newEvent(b3.Id())
	event7 := newEvent(b3.Id())
	event8 := newEvent(b3.Id())
	event9 := newEvent(b3.Id())

	sl1 := newSleep(*event1)
	sl1.SetStart(getDataForTest())
	sl2 := newSleep(*event4)
	sl2.SetStart(getDataForTest())
	sl3 := newSleep(*event6)
	sl3.SetStart(getDataForTest())

	var sleeps []sleepI
	sleeps = append(sleeps, sl1, sl2, sl3)

	ea1 := newEat(*event2)
	ea1.SetDescription("test_eat1")
	ea2 := newEat(*event3)
	ea2.SetDescription("test_eat2")
	ea3 := newEat(*event7)
	ea3.SetDescription("test_eat1")

	var eats []eatI
	eats = append(eats, ea1, ea2, ea3)

	pamp1 := newPampers(*event5)
	pamp1.setState(dirty)
	pamp2 := newPampers(*event8)
	pamp2.setState(dirty)
	pamp3 := newPampers(*event9)
	pamp3.setState(dirty)

	var pampers []pampersI
	pampers = append(pampers, pamp1, pamp2, pamp3)

}

func TestEventCreate(t *testing.T) {
	require.Equal(t, 1, 1, "ok")
	//require.True(t, true, "ok")
}

func TestGetTypedEventsIdsByBabyDate(t *testing.T) {

}
func TestGetEventsByBabyDate(t *testing.T) {

}
func TestGetNotEndedSleepForBaby(t *testing.T) {

}
