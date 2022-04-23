package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getDataForTest() time.Time {

	min := time.Date(1975, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(1975, 2, 0, 0, 0, 0, 0, time.UTC).Unix()
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
	for _, p := range parents {
		p.writeStructToBase()
	}

	var babyes []babyI
	b1 := newBaby()
	b1.SetBirth(getDataForTest())
	b1.SetName("testBaby1")
	b1.SetParent(p1.Id())
	p1.CurrentBaby()
	b2 := newBaby()
	b2.SetBirth(getDataForTest())
	b2.SetName("testBaby2")
	b2.SetParent(p1.Id())

	b3 := newBaby()
	b3.SetBirth(getDataForTest())
	b3.SetName("testBaby3")
	b3.SetParent(p2.Id())

	babyes = append(babyes, b1, b2, b3)
	for _, b := range babyes {
		b.writeStructToBase()
	}
	err := p1.SetCurrentBaby(b1.Id())
	fmt.Println(err)
	err = p2.SetCurrentBaby(b3.Id())
	fmt.Println(err)

	for _, p := range parents {
		p.writeStructToBase()
	}
	var ebw []eventBaseWorker
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
	sl1.SetEndTime(getDataForTest())
	sl2 := newSleep(*event4)
	sl2.SetStart(getDataForTest())
	sl2.SetEndTime(getDataForTest())
	sl3 := newSleep(*event6)
	sl3.SetStart(getDataForTest())
	sl3.SetEndTime(getDataForTest())

	//var sleeps []sleepI0
	//sleeps = append(sleeps, sl1, sl2, sl3)
	ebw = append(ebw, sl1, sl2, sl3)

	ea1 := newEat(*event2)
	ea1.SetDescription("test_eat1")
	ea1.SetStart(getDataForTest())
	ea2 := newEat(*event3)
	ea2.SetDescription("test_eat2")
	ea2.SetStart(getDataForTest())
	ea3 := newEat(*event7)
	ea3.SetDescription("test_eat1")
	ea3.SetStart(getDataForTest())
	//var eats []eatI
	//eats = append(eats, ea1, ea2, ea3)
	ebw = append(ebw, ea1, ea2, ea3)

	pamp1 := newPampers(*event5)
	pamp1.SetState(dirty)
	pamp1.SetStart(getDataForTest())
	pamp2 := newPampers(*event8)
	pamp2.SetState(dirty)
	pamp2.SetStart(getDataForTest())
	pamp3 := newPampers(*event9)
	pamp3.SetState(dirty)
	pamp3.SetStart(getDataForTest())
	p1.CurrentBaby()
	//var pampers []pampersI
	//TODO: if getcurrentbaby == err - set to 0
	//pampers = append(pampers, pamp1, pamp2, pamp3)
	ebw = append(ebw, pamp1, pamp2, pamp3)

	for _, e := range ebw {
		e.writeStructToBase()
	}

}

func cleanupDB() {
	eventDeleteQuery := "delete from %s where %s < '1977-01-01' or %s > '%s'"
	eventTables := []string{"eat", "pampers", "sleep"}
	for _, table := range eventTables {
		err := DBDeleteData(fmt.Sprintf(
			eventDeleteQuery, table, "start", "start", time.Now().Format("2006-01-02")))
		if err != nil {
			fmt.Println(err)
		}
	}
	DBDeleteData(fmt.Sprintf("delete from sleep * where start > '%s'",
		time.Now().Format("2006-01-02")))
	if err := DBDeleteData(fmt.Sprintf(
		eventDeleteQuery,
		"baby",
		"birth",
		"birth",
		time.Now().Format("2006-01-02"))); err != nil {
	}
	if err := DBDeleteData("delete from parent where parent_id=1 or parent_id=2 or parent_id=125"); err != nil {

	}
}

func TestWriteStructToBaseReadStructFromBase(t *testing.T) {
	cleanupDB()
	createTestData()
	//defer cleanupDB()

	// PARENT
	testParent := newParent()
	assert.Error(t, testParent.readStructFromBase(-123), "Tests read not existing parent id")
	assert.NoError(t, testParent.readStructFromBase(1), "Test read existing parent id")
	require.Equal(t, testParent.name, "TestParent1", "test receive parent name")

	//checks succesfully write parent with new id to base
	testParent.SetName("NewNameTestParent1")
	assert.NoError(t, testParent.writeStructToBase(), "Test write parent struct to base")
	testParent2 := newParent()
	assert.NoError(t, testParent2.readStructFromBase(1), "Test read parent with corret id from base")
	require.Equal(t, testParent2.name, "NewNameTestParent1", "test succesfylly write and read new ID")

	// BABY
	//checks get babyes from parent
	b, err := GetBabyesByParent(testParent2.Id())
	assert.NoError(t, err, "Tests get babyes by parent")
	assert.NotEmpty(t, b, "Tests get babyes by parent")
	assert.Equal(t, len(b), 2, "Test babyes count")

	//checks read concrete baby and write concrete baby
	currentBaby := b[0]
	assert.Error(t, currentBaby.readStructFromBase(-123),
		"test baby read with wrong id")
	assert.NoError(t, currentBaby.readStructFromBase(currentBaby.Id()),
		"tests baby read struct")

	currentBaby.SetBirth(getDataForTest())
	currentBaby.SetName("testsbabyname")
	assert.Error(t, currentBaby.SetParent(-125))
	assert.NoError(t, currentBaby.SetParent(testParent.Id()))

	assert.NoError(t, currentBaby.writeStructToBase(), "test write new baby value to base")

	nb := newBaby()
	assert.NoError(t, nb.readStructFromBase(currentBaby.Id()), "Checks re read writed baby")
	assert.Equal(t, nb.Id(), currentBaby.Id(), "Check equals writed baby and readed baby")
	assert.Equal(t, nb.Birth().Format("2006-01-02"), currentBaby.Birth().Format("2006-01-02"), "Check equals writed baby and readed baby")
	assert.Equal(t, nb.ParentId(), currentBaby.ParentId(), "Check equals writed baby and readed baby")

	// EAT
	event := newEvent(currentBaby.Id())
	e := newEat(*event)
	e.SetStart(getDataForTest())
	e.SetDescription("eat from tests")
	assert.NoError(t, e.writeStructToBase(), "checks write eat to base")

	ne := newEat(*newEvent(currentBaby.Id()))
	assert.NoError(t, ne.readStructFromBase(e.Id()), "Checks read eat from base")

	fmt.Println(e.Start())
	fmt.Println(ne.Start())

	assert.Equal(t, ne.Id(), e.Id(), "Tests readed eat == writed eat")
	assert.Equal(t, ne.Start().Format("2006-01-02 15:04"),
		e.Start().Format("2006-01-02 15:04"),
		"Tests readed eat == writed eat")
	assert.Equal(t, ne.Description(), e.Description(), "Tests read eat == write eat")

	// PAMPERS
	event = newEvent(currentBaby.Id())
	p := newPampers(*event)
	p.SetStart(getDataForTest())
	assert.Error(t, p.SetState(5), "Set not valid state to pampers")
	assert.NoError(t, p.SetState(wet), "Set valid state to pampers")

	assert.NoError(t, p.writeStructToBase(), "test writes pampers to base")

	np := newPampers(*newEvent(currentBaby.Id()))
	assert.NoError(t, np.readStructFromBase(p.Id()), "test reads writed pampers from base")

	assert.Equal(t, p.BabyId(), np.BabyId(), "tests readed pampers == writed pampers")
	assert.Equal(t, p.Start().Format("2006-01-02 15:04"),
		np.Start().Format("2006-01-02 15:04"),
		"tests readed pampers == writed pampers")
	assert.Equal(t, p.State(), np.State(), "tests readed pampers == writed pampers")

	// SLEEP
	event = newEvent(currentBaby.Id())
	sl := newSleep(*event)

	sl.SetStart(getDataForTest())

	assert.NoError(t, sl.writeStructToBase(), "tests write not ended sleep to base")

	ns := newSleep(*newEvent(currentBaby.Id()))
	assert.NoError(t, ns.readStructFromBase(sl.Id()), "tests read writed not ended sleep")

	assert.Equal(t, ns.BabyId(), sl.BabyId(), "test readed not ended sleep == writed")
	assert.Equal(t, ns.Start().Format("2006-01-02 15:04"),
		sl.Start().Format("2006-01-02 15:04"),
		"test readed not ended sleep == writed")

	sl.SetStart(getDataForTest())
	assert.NoError(t, sl.writeStructToBase(), "tests updated not ended sleep to base")
	assert.NoError(t, ns.readStructFromBase(sl.Id()), "tests read updated not ended sleep")
	assert.Equal(t, ns.Start().Format("2006-01-02 15:04"),
		sl.Start().Format("2006-01-02 15:04"),
		"test readed updated not ended sleep == writed")

	ns.SetEndTime(getDataForTest())
	assert.NoError(t, ns.writeStructToBase(), "Write end of sleep")
	assert.NoError(t, sl.readStructFromBase(ns.Id()), "test read writed ended")

	assert.Equal(t, ns.End().Format("2006-01-02 15:04"), sl.End().Format("2006-01-02 15:04"),
		"tests end time write == end time readed")

	sl = newSleep(*newEvent(currentBaby.Id()))
	sl.SetStart(getDataForTest())
	sl.SetEndTime(getDataForTest())

	assert.NoError(t, sl.writeStructToBase(), "test write sleep with end")
	assert.NoError(t, ns.readStructFromBase(sl.Id()), "test read sleep writed with end")

	assert.Equal(t,
		ns.Start().Format("2006-01-02 15:04"),
		sl.Start().Format("2006-01-02 15:04"),
		"tests readed with end sleep == writed with end sleep")

	assert.Equal(t,
		ns.End().Format("2006-01-02 15:04"),
		sl.End().Format("2006-01-02 15:04"),
		"tests readed with end sleep == writed with end sleep")
}

func TestRegisterParent(t *testing.T) {
	var testId int64 = 125
	var testName string = "testPar"
	tp, err := RegisterNewParent(testId, testName)
	assert.Nil(t, err, "Testing parent created")
	assert.Equal(t, tp.Id(), testId, "Test id is ok")
	assert.Equal(t, tp.Name(), testName)
	assert.Equal(t, tp.CurrentBaby(), int64(0), "test current baby is nil")
}

func TestCurrentBaby(t *testing.T) {
	createTestData()

	defer cleanupDB()

	p := newParent()
	p.readStructFromBase(1)
	assert.Equal(t, p.Id(), int64(1))

	assert.NotEqual(t, p.CurrentBaby(), int64(0))
	assert.Error(t, p.SetCurrentBaby(-125), "tests not set not exist baby")
}

func TestGetEventsByBabyDate(t *testing.T) {
	createTestData()
	defer cleanupDB()

	p1 := newParent()
	assert.Nil(t, p1.readStructFromBase(1), "Tests read parent")
	assert.NotEqual(t, p1.CurrentBaby(), int64(0), "Test parent have baby")

}

func TestBabyActivity(t *testing.T) {
	createTestData()
	defer cleanupDB()

	UA := NewBabyActivity()

	UA.setAction("add")
	// test basic activity
	res, err := UA.doActivity([]string{""})
	assert.Error(t, err, "Tests babyactivity with less args")
	assert.Empty(t, res, "Tests babyactivity with less args")

	res, err = UA.doActivity([]string{"abc"})
	assert.Error(t, err, "Tests babyactivity with not numeric parent id")
	assert.Empty(t, res, "Tests babyactivity with not numeric parent id")

	res, err = UA.doActivity([]string{"-135"})
	assert.Error(t, err, "Tests babyactivity with not numeric parent id")
	assert.Empty(t, res, "Tests babyactivity with not numeric parent id")

	// test add action
	UA.setAction("add")

	res, err = UA.doActivity([]string{"1", "testname"})
	assert.Error(t, err, "Tests add baby with less args")
	assert.Empty(t, res, "Tests add baby with less args")

	res, err = UA.doActivity([]string{"1", "testname", "abc"})
	assert.Error(t, err, "Tests add baby with not valid time")
	assert.Empty(t, res, "Tests add baby with not valid time")

	res, err = UA.doActivity([]string{"1", "testname", "abc"})
	assert.Error(t, err, "Tests add baby with not valid time")
	assert.Empty(t, res, "Tests add baby with not valid time")

	res, err = UA.doActivity([]string{"1", "testname", "1975-03-12"})
	assert.Nil(t, err, "Tests baby added")
	assert.NotEmpty(t, res, "Tests baby succesfully added")

	//test set action
	UA.setAction("current")

	res, err = UA.doActivity([]string{"1", "abc"})
	assert.Error(t, err, "Test set current baby with not digit number")
	assert.Empty(t, res, "Test set current baby with not digit number")

	res, err = UA.doActivity([]string{"1", "5392"})
	assert.Error(t, err, "Test set current baby with not existing number")
	assert.Empty(t, res, "Test set current baby with not existing number")

	res, err = UA.doActivity([]string{"1", "1"})
	assert.NoError(t, err, "Test set current baby")
	assert.NotEmpty(t, res, "Test set current baby")

	currentB, err := GetCurrentBaby(1)
	assert.NoError(t, err, "test current baby is setted")
	allBs, _ := GetBabyesByParent(1)

	founded := false
	for _, b := range allBs {
		if currentB.Id() == b.Id() {
			founded = true
			break
		}
	}

	assert.True(t, founded, "tests id of setted current baby == current_baby of parent")

	// test remove
	UA.setAction("remove")
	res, err = UA.doActivity([]string{"1", "abc"})
	assert.Error(t, err, "Test remove current baby with not digit number")
	assert.Empty(t, res, "Test remove current baby with not digit number")

	res, err = UA.doActivity([]string{"1", "5392"})
	assert.Error(t, err, "Test  remove baby with not existing number")
	assert.Empty(t, res, "Test  remove baby with not existing number")

	res, err = UA.doActivity([]string{"1", "1"})
	assert.NoError(t, err, "Test remove current baby")
	assert.NotEmpty(t, res, "Test remove current baby")

	p := newParent()
	err = p.readStructFromBase(1)
	assert.Equal(t, int64(0), p.CurrentBaby(), "Test when remove current baby - it setted to 0")
	assert.NoError(t, err, "Test when remove current baby - it setted to 0")

	res, err = UA.doActivity([]string{" "})

}

func TestEventActivity(t *testing.T) {
	cleanupDB()
	createTestData()
	t.Cleanup(cleanupDB)
	defer cleanupDB()
	UA := NewEventActivity()

	// test basic activity

	res, err := UA.doActivity([]string{""})
	assert.Error(t, err, "Tests EventActivity with less args")
	assert.Empty(t, res, "Tests EventActivity with less args")

	res, err = UA.doActivity([]string{"abc"})
	assert.Error(t, err, "Tests EventActivity with not numeric parent id")
	assert.Empty(t, res, "Tests EventActivity with not numeric parent id")

	res, err = UA.doActivity([]string{"-135"})
	assert.Error(t, err, "Tests EventActivity with no existing parent id")
	assert.Empty(t, res, "Tests EventActivity with no existing parent id")

	currentB, err := GetCurrentBaby(1)
	assert.NoError(t, err, "tests current baby exists")
	assert.NotNil(t, currentB, "tests current baby exists")

	p := newParent()
	p.readStructFromBase(1)
	actualCurrentBaby := p.CurrentBaby()

	p.SetCurrentBaby(0)
	p.writeStructToBase()

	res, err = UA.doActivity([]string{"1", "03.05"})
	assert.Error(t, err, "try to write with current_baby = 0 for parent")
	assert.Empty(t, res, "try to write with current_baby = 0 for parent")

	p.SetCurrentBaby(actualCurrentBaby)
	p.writeStructToBase()

	res, err = UA.doActivity([]string{"1", "1"})
	assert.Error(t, err, "tests event with not valid time")
	assert.Empty(t, res, "tests event with not valid time")

	res, err = UA.doActivity([]string{"1", "03.05"})
	assert.Error(t, err, "tests event with not valid time format with only time")
	assert.Empty(t, res, "tests event with not valid time format with only time")

	res, err = UA.doActivity([]string{"1", "2006+01*02_03.05"})
	assert.Error(t, err, "tests event with not valid time format with time and date")
	assert.Empty(t, res, "tests event with not valid time format with time and date")

	res, err = UA.doActivity([]string{"1", "2006+01*02_03.05"})
	assert.Error(t, err, "tests event with not valid time format with time and date")
	assert.Empty(t, res, "tests event with not valid time format with time and date")

	// ACTION SLEEP
	UA.setAction("sleep")
	res, err = UA.doActivity([]string{"1", "15:04"})
	assert.NoError(t, err, "valid only time")
	assert.NotEmpty(t, res, "valid only time")

	b := newBaby()
	b.readStructFromBase(p.CurrentBaby())

	//new created sleep must be not ended
	notEndedSleep, _ := GetNotEndedSleepForBaby(b.Id())
	sleepsIds, err := GetTypedEventsIdsByBabyDate(p.CurrentBaby(), time.Now(), "sleep")
	assert.NoError(t, err, fmt.Sprintf("get all sleeps by date error: %s", err))

	assert.Equal(t, notEndedSleep.Id(), sleepsIds[0],
		"baby have not ended sleep. And it single sleep today")

	res, err = UA.doActivity([]string{"1", time.Now().Format("2006-1-02_15:04")})
	assert.NoError(t, err, "set end of sleep")
	assert.NotEmpty(t, res, "set end of sleep")

	notEndedSleep.readStructFromBase(notEndedSleep.Id())
	assert.NotEmpty(t, notEndedSleep.End(), "end of sleep setted")
	assert.Equal(t,
		time.Now().Format("2006-01-02 15:04"),
		notEndedSleep.End().Format("2006-01-02 15:04"),
		"sleep ended tommorow")

	// ACTION PAMPERS
	UA.setAction("pampers")
	p = newParent()
	p.readStructFromBase(1)

	res, err = UA.doActivity([]string{"1", "15:04"})
	assert.Error(t, err, "without parmers state")
	assert.Empty(t, res, "without parmers state")

	res, err = UA.doActivity([]string{"1", "15:04", "abc"})
	assert.Error(t, err, "wrong parmers state")
	assert.Empty(t, res, "wrong parmers state")

	res, err = UA.doActivity([]string{"1", "15:04", "wet"})
	assert.NoError(t, err, "wet parmers state")
	assert.NotEmpty(t, res, "wet parmers state")

	res, err = UA.doActivity([]string{"1", "15:04", "dirty"})
	assert.NoError(t, err, "dirty parmers state")
	assert.NotEmpty(t, res, "dirty parmers state")

	res, err = UA.doActivity([]string{"1", "15:04", "combined"})
	assert.NoError(t, err, "combined parmers state")
	assert.NotEmpty(t, res, "combined parmers state")

	pampersIds, err := GetTypedEventsIdsByBabyDate(p.CurrentBaby(), time.Now(), "pampers")

	assert.NoError(t, err, fmt.Sprintf("tests get pampers by date: %s", err))
	assert.Equal(t, 3, len(pampersIds), "tests count of added pampers  == 3")

	pamp := newPampers(*newEvent(p.CurrentBaby()))
	pamp.readStructFromBase(pampersIds[0])
	assert.Equal(t, wet, pamp.state, "tests wet state of pamp succesfully writed")

	pamp.readStructFromBase(pampersIds[1])
	assert.Equal(t, dirty, pamp.state, "tests dirty state of pamp succesfully writed")

	pamp.readStructFromBase(pampersIds[2])
	assert.Equal(t, combined, pamp.state, "tests combined state of pamp succesfully writed")

	UA.setAction("eat")

	res, err = UA.doActivity([]string{"1", "15:04"})
	assert.NoError(t, err, "tests add eat without description")
	assert.NotEmpty(t, res, "tests adde at without description")
	res, err = UA.doActivity([]string{"1", "15:04", ""})
	assert.NoError(t, err, "tests add eat with empty description")
	assert.NotEmpty(t, res, "tests add eat with empty description")

	res, err = UA.doActivity([]string{"1", "15:04", "this is test eat"})
	assert.NoError(t, err, "tests add eat with description")
	assert.NotEmpty(t, res, "tests add eat with description")

}

func TestStateActivity(t *testing.T) {
	createTestData()
	defer cleanupDB()

	UA := NewStateActivity()

	UA.setAction("today")
	// test basic activity

	res, err := UA.doActivity([]string{""})
	assert.Error(t, err, "Tests StateActivity with less args")
	assert.Empty(t, res, "Tests StateActivity with less args")

	res, err = UA.doActivity([]string{"abc"})
	assert.Error(t, err, "Tests StateActivity with not numeric parent id")
	assert.Empty(t, res, "Tests StateActivity with not numeric parent id")

	res, err = UA.doActivity([]string{"-135"})
	assert.Error(t, err, "Tests StateActivity with not numeric parent id")
	assert.Empty(t, res, "Tests StateActivity with not numeric parent id")

	UA.setAction("week")
	UA.setAction("month")
}

/*

тестировать сам telegrambotapi я не могу.
но могу тестировать то, что я из него дергаю.

1. выборки, которые делаются по пользвоательским командам с тестовыми данными
*/

func TestCreate(t *testing.T) {
	createTestData()
}

func TestCleanup(t *testing.T) {
	cleanupDB()
}
