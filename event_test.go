package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// замена require на require
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
		fmt.Println(err.Error())
	}
	if err := DBDeleteData("delete from parent where parent_id=1 or parent_id=2 or parent_id=125"); err != nil {
		fmt.Println(err.Error())

	}
}

func TestWriteStructToBaseReadStructFromBase(t *testing.T) {
	cleanupDB()
	createTestData()
	defer cleanupDB()

	// PARENT
	testParent := newParent()
	require.Error(t, testParent.readStructFromBase(-123), "Tests read not existing parent id")
	require.NoError(t, testParent.readStructFromBase(1), "Test read existing parent id")
	require.Equal(t, testParent.name, "TestParent1", "test receive parent name")

	//checks succesfully write parent with new id to base
	testParent.SetName("NewNameTestParent1")
	require.NoError(t, testParent.writeStructToBase(), "Test write parent struct to base")
	testParent2 := newParent()
	require.NoError(t, testParent2.readStructFromBase(1), "Test read parent with corret id from base")
	require.Equal(t, testParent2.name, "NewNameTestParent1", "test succesfylly write and read new ID")

	// BABY
	//checks get babyes from parent
	b, err := GetBabyesByParent(testParent2.Id())
	require.NoError(t, err, "Tests get babyes by parent")
	require.NotEmpty(t, b, "Tests get babyes by parent")
	require.Equal(t, len(b), 2, "Test babyes count")

	//checks read concrete baby and write concrete baby
	currentBaby := b[0]
	require.Error(t, currentBaby.readStructFromBase(-123),
		"test baby read with wrong id")
	require.NoError(t, currentBaby.readStructFromBase(currentBaby.Id()),
		"tests baby read struct")

	currentBaby.SetBirth(getDataForTest())
	currentBaby.SetName("testsbabyname")
	require.Error(t, currentBaby.SetParent(-125))
	require.NoError(t, currentBaby.SetParent(testParent.Id()))

	require.NoError(t, currentBaby.writeStructToBase(), "test write new baby value to base")

	nb := newBaby()
	require.NoError(t, nb.readStructFromBase(currentBaby.Id()), "Checks re read writed baby")
	require.Equal(t, nb.Id(), currentBaby.Id(), "Check equals writed baby and readed baby")
	require.Equal(t, nb.Birth().Format("2006-01-02"), currentBaby.Birth().Format("2006-01-02"), "Check equals writed baby and readed baby")
	require.Equal(t, nb.ParentId(), currentBaby.ParentId(), "Check equals writed baby and readed baby")

	// EAT
	event := newEvent(currentBaby.Id())
	e := newEat(*event)
	e.SetStart(getDataForTest())
	e.SetDescription("eat from tests")
	require.NoError(t, e.writeStructToBase(), "checks write eat to base")

	ne := newEat(*newEvent(currentBaby.Id()))
	require.NoError(t, ne.readStructFromBase(e.Id()), "Checks read eat from base")

	fmt.Println(e.Start())
	fmt.Println(ne.Start())

	require.Equal(t, ne.Id(), e.Id(), "Tests readed eat == writed eat")
	require.Equal(t, ne.Start().Format("2006-01-02 15:04"),
		e.Start().Format("2006-01-02 15:04"),
		"Tests readed eat == writed eat")
	require.Equal(t, ne.Description(), e.Description(), "Tests read eat == write eat")

	// PAMPERS
	event = newEvent(currentBaby.Id())
	p := newPampers(*event)
	p.SetStart(getDataForTest())
	require.Error(t, p.SetState(5), "Set not valid state to pampers")
	require.NoError(t, p.SetState(wet), "Set valid state to pampers")

	require.NoError(t, p.writeStructToBase(), "test writes pampers to base")

	np := newPampers(*newEvent(currentBaby.Id()))
	require.NoError(t, np.readStructFromBase(p.Id()), "test reads writed pampers from base")

	require.Equal(t, p.BabyId(), np.BabyId(), "tests readed pampers == writed pampers")
	require.Equal(t, p.Start().Format("2006-01-02 15:04"),
		np.Start().Format("2006-01-02 15:04"),
		"tests readed pampers == writed pampers")
	require.Equal(t, p.State(), np.State(), "tests readed pampers == writed pampers")

	// SLEEP
	event = newEvent(currentBaby.Id())
	sl := newSleep(*event)

	sl.SetStart(getDataForTest())

	require.NoError(t, sl.writeStructToBase(), "tests write not ended sleep to base")

	ns := newSleep(*newEvent(currentBaby.Id()))
	require.NoError(t, ns.readStructFromBase(sl.Id()), "tests read writed not ended sleep")

	require.Equal(t, ns.BabyId(), sl.BabyId(), "test readed not ended sleep == writed")
	require.Equal(t, ns.Start().Format("2006-01-02 15:04"),
		sl.Start().Format("2006-01-02 15:04"),
		"test readed not ended sleep == writed")

	sl.SetStart(getDataForTest())
	require.NoError(t, sl.writeStructToBase(), "tests updated not ended sleep to base")
	require.NoError(t, ns.readStructFromBase(sl.Id()), "tests read updated not ended sleep")
	require.Equal(t, ns.Start().Format("2006-01-02 15:04"),
		sl.Start().Format("2006-01-02 15:04"),
		"test readed updated not ended sleep == writed")

	ns.SetEndTime(getDataForTest())
	require.NoError(t, ns.writeStructToBase(), "Write end of sleep")
	require.NoError(t, sl.readStructFromBase(ns.Id()), "test read writed ended")

	require.Equal(t, ns.End().Format("2006-01-02 15:04"), sl.End().Format("2006-01-02 15:04"),
		"tests end time write == end time readed")

	sl = newSleep(*newEvent(currentBaby.Id()))
	sl.SetStart(getDataForTest())
	sl.SetEndTime(getDataForTest())

	require.NoError(t, sl.writeStructToBase(), "test write sleep with end")
	require.NoError(t, ns.readStructFromBase(sl.Id()), "test read sleep writed with end")

	require.Equal(t,
		ns.Start().Format("2006-01-02 15:04"),
		sl.Start().Format("2006-01-02 15:04"),
		"tests readed with end sleep == writed with end sleep")

	require.Equal(t,
		ns.End().Format("2006-01-02 15:04"),
		sl.End().Format("2006-01-02 15:04"),
		"tests readed with end sleep == writed with end sleep")
}

func TestRegisterParent(t *testing.T) {
	var testId int64 = 125
	var testName string = "testPar"
	tp, err := RegisterNewParent(testId, testName)
	require.Nil(t, err, "Testing parent created")
	require.Equal(t, tp.Id(), testId, "Test id is ok")
	require.Equal(t, tp.Name(), testName)
	require.Equal(t, tp.CurrentBaby(), int64(0), "test current baby is nil")
}

func TestCurrentBaby(t *testing.T) {
	createTestData()

	defer cleanupDB()

	p := newParent()
	p.readStructFromBase(1)
	require.Equal(t, p.Id(), int64(1))

	require.NotEqual(t, p.CurrentBaby(), int64(0))
	require.Error(t, p.SetCurrentBaby(-125), "tests not set not exist baby")
}

func TestGetEventsByBabyDate(t *testing.T) {
	createTestData()
	defer cleanupDB()

	p1 := newParent()
	require.Nil(t, p1.readStructFromBase(1), "Tests read parent")
	require.NotEqual(t, p1.CurrentBaby(), int64(0), "Test parent have baby")

}

func TestBabyActivity(t *testing.T) {
	createTestData()
	defer cleanupDB()

	UA := NewBabyActivity()

	UA.setAction("add")
	// test basic activity
	res, err := UA.doActivity([]string{""})
	require.Error(t, err, "Tests babyactivity with less args")
	require.Empty(t, res, "Tests babyactivity with less args")

	res, err = UA.doActivity([]string{"abc"})
	require.Error(t, err, "Tests babyactivity with not numeric parent id")
	require.Empty(t, res, "Tests babyactivity with not numeric parent id")

	res, err = UA.doActivity([]string{"-135"})
	require.Error(t, err, "Tests babyactivity with not numeric parent id")
	require.Empty(t, res, "Tests babyactivity with not numeric parent id")

	// test add action
	UA.setAction("add")

	res, err = UA.doActivity([]string{"1", "testname"})
	require.Error(t, err, "Tests add baby with less args")
	require.Empty(t, res, "Tests add baby with less args")

	res, err = UA.doActivity([]string{"1", "testname", "abc"})
	require.Error(t, err, "Tests add baby with not valid time")
	require.Empty(t, res, "Tests add baby with not valid time")

	res, err = UA.doActivity([]string{"1", "testname", "abc"})
	require.Error(t, err, "Tests add baby with not valid time")
	require.Empty(t, res, "Tests add baby with not valid time")

	res, err = UA.doActivity([]string{"1", "testname", "1975-03-12"})
	require.Nil(t, err, "Tests baby added")
	require.NotEmpty(t, res, "Tests baby succesfully added")

	//test set action
	UA.setAction("current")

	res, err = UA.doActivity([]string{"1", "abc"})
	require.Error(t, err, "Test set current baby with not digit number")
	require.Empty(t, res, "Test set current baby with not digit number")

	res, err = UA.doActivity([]string{"1", "5392"})
	require.Error(t, err, "Test set current baby with not existing number")
	require.Empty(t, res, "Test set current baby with not existing number")

	res, err = UA.doActivity([]string{"1", "0"})
	require.Error(t, err, "Test set current baby to zero")
	require.Empty(t, res, "Test set current baby to zero")

	res, err = UA.doActivity([]string{"1", "1"})
	require.NoError(t, err, "Test set current baby")
	require.NotEmpty(t, res, "Test set current baby")

	currentB, err := GetCurrentBaby(1)
	require.NoError(t, err, "test current baby is setted")
	allBs, _ := GetBabyesByParent(1)

	founded := false
	for _, b := range allBs {
		if currentB.Id() == b.Id() {
			founded = true
			break
		}
	}

	require.True(t, founded, "tests id of setted current baby == current_baby of parent")

	// test remove
	UA.setAction("remove")
	res, err = UA.doActivity([]string{"1", "abc"})
	require.Error(t, err, "Test remove current baby with not digit number")
	require.Empty(t, res, "Test remove current baby with not digit number")

	res, err = UA.doActivity([]string{"1", "5392"})
	require.Error(t, err, "Test  remove baby with not existing number")
	require.Empty(t, res, "Test  remove baby with not existing number")

	res, err = UA.doActivity([]string{"1", "1"})
	require.NoError(t, err, "Test remove current baby")
	require.NotEmpty(t, res, "Test remove current baby")

	p := newParent()
	err = p.readStructFromBase(1)
	require.Equal(t, int64(0), p.CurrentBaby(), "Test when remove current baby - it setted to 0")
	require.NoError(t, err, "Test when remove current baby - it setted to 0")

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
	require.Error(t, err, "Tests EventActivity with less args")
	require.Empty(t, res, "Tests EventActivity with less args")

	res, err = UA.doActivity([]string{"abc"})
	require.Error(t, err, "Tests EventActivity with not numeric parent id")
	require.Empty(t, res, "Tests EventActivity with not numeric parent id")

	res, err = UA.doActivity([]string{"-135"})
	require.Error(t, err, "Tests EventActivity with no existing parent id")
	require.Empty(t, res, "Tests EventActivity with no existing parent id")

	currentB, err := GetCurrentBaby(1)
	require.NoError(t, err, "tests current baby exists")
	require.NotNil(t, currentB, "tests current baby exists")

	p := newParent()
	p.readStructFromBase(1)
	actualCurrentBaby := p.CurrentBaby()

	p.SetCurrentBaby(0)
	p.writeStructToBase()

	res, err = UA.doActivity([]string{"1", "03.05"})
	require.Error(t, err, "try to write with current_baby = 0 for parent")
	require.Empty(t, res, "try to write with current_baby = 0 for parent")

	p.SetCurrentBaby(actualCurrentBaby)
	p.writeStructToBase()

	res, err = UA.doActivity([]string{"1", "1"})
	require.Error(t, err, "tests event with not valid time")
	require.Empty(t, res, "tests event with not valid time")

	res, err = UA.doActivity([]string{"1", "03.05"})
	require.Error(t, err, "tests event with not valid time format with only time")
	require.Empty(t, res, "tests event with not valid time format with only time")

	res, err = UA.doActivity([]string{"1", "2006+01*02_03.05"})
	require.Error(t, err, "tests event with not valid time format with time and date")
	require.Empty(t, res, "tests event with not valid time format with time and date")

	res, err = UA.doActivity([]string{"1", "2006+01*02_03.05"})
	require.Error(t, err, "tests event with not valid time format with time and date")
	require.Empty(t, res, "tests event with not valid time format with time and date")

	// ACTION SLEEP
	UA.setAction("sleep")
	res, err = UA.doActivity([]string{"1", "15:04"})
	require.NoError(t, err, "valid only time")
	require.NotEmpty(t, res, "valid only time")

	b := newBaby()
	b.readStructFromBase(p.CurrentBaby())

	//new created sleep must be not ended
	notEndedSleep, _ := GetNotEndedSleepForBaby(b.Id())
	sleepsIds, err := GetTypedEventsIdsByBabyDate(p.CurrentBaby(), time.Now(), "sleep")
	require.NoError(t, err, fmt.Sprintf("get all sleeps by date error: %s", err))

	require.Equal(t, notEndedSleep.Id(), sleepsIds[0],
		"baby have not ended sleep. And it single sleep today")

	res, err = UA.doActivity([]string{"1", time.Now().Format("2006-1-02_15:04")})
	require.NoError(t, err, "set end of sleep")
	require.NotEmpty(t, res, "set end of sleep")

	notEndedSleep.readStructFromBase(notEndedSleep.Id())
	require.NotEmpty(t, notEndedSleep.End(), "end of sleep setted")
	require.Equal(t,
		time.Now().Format("2006-01-02 15:04"),
		notEndedSleep.End().Format("2006-01-02 15:04"),
		"sleep ended tommorow")

	// ACTION PAMPERS
	UA.setAction("pampers")
	p = newParent()
	p.readStructFromBase(1)

	res, err = UA.doActivity([]string{"1", "15:04"})
	require.Error(t, err, "without parmers state")
	require.Empty(t, res, "without parmers state")

	res, err = UA.doActivity([]string{"1", "15:04", "abc"})
	require.Error(t, err, "wrong parmers state")
	require.Empty(t, res, "wrong parmers state")

	res, err = UA.doActivity([]string{"1", "15:04", "wet"})
	require.NoError(t, err, "wet parmers state")
	require.NotEmpty(t, res, "wet parmers state")

	res, err = UA.doActivity([]string{"1", "15:04", "dirty"})
	require.NoError(t, err, "dirty parmers state")
	require.NotEmpty(t, res, "dirty parmers state")

	res, err = UA.doActivity([]string{"1", "15:04", "combined"})
	require.NoError(t, err, "combined parmers state")
	require.NotEmpty(t, res, "combined parmers state")

	pampersIds, err := GetTypedEventsIdsByBabyDate(p.CurrentBaby(), time.Now(), "pampers")

	require.NoError(t, err, fmt.Sprintf("tests get pampers by date: %s", err))
	require.Equal(t, 3, len(pampersIds), "tests count of added pampers  == 3")

	pamp := newPampers(*newEvent(p.CurrentBaby()))
	pamp.readStructFromBase(pampersIds[0])
	require.Equal(t, wet, pamp.state, "tests wet state of pamp succesfully writed")

	pamp.readStructFromBase(pampersIds[1])
	require.Equal(t, dirty, pamp.state, "tests dirty state of pamp succesfully writed")

	pamp.readStructFromBase(pampersIds[2])
	require.Equal(t, combined, pamp.state, "tests combined state of pamp succesfully writed")

	UA.setAction("eat")

	res, err = UA.doActivity([]string{"1", "15:04"})
	require.NoError(t, err, "tests add eat without description")
	require.NotEmpty(t, res, "tests adde at without description")
	res, err = UA.doActivity([]string{"1", "15:04", ""})
	require.NoError(t, err, "tests add eat with empty description")
	require.NotEmpty(t, res, "tests add eat with empty description")

	res, err = UA.doActivity([]string{"1", "15:04", "this is test eat"})
	require.NoError(t, err, "tests add eat with description")
	require.NotEmpty(t, res, "tests add eat with description")

}

func TestStateActivity(t *testing.T) {
	createTestData()
	defer cleanupDB()

	UA := NewStateActivity()

	UA.setAction("today")
	// test basic activity

	res, err := UA.doActivity([]string{""})
	require.Error(t, err, "Tests StateActivity with less args")
	require.Empty(t, res, "Tests StateActivity with less args")

	res, err = UA.doActivity([]string{"abc"})
	require.Error(t, err, "Tests StateActivity with not numeric parent id")
	require.Empty(t, res, "Tests StateActivity with not numeric parent id")

	res, err = UA.doActivity([]string{"-135"})
	require.Error(t, err, "Tests StateActivity with not numeric parent id")
	require.Empty(t, res, "Tests StateActivity with not numeric parent id")

	res, err = UA.doActivity([]string{"0"})
	require.Error(t, err, "Tests StateActivity with not numeric parent id")
	require.Empty(t, res, "Tests StateActivity with not numeric parent id")

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
