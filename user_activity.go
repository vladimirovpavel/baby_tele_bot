package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userActivityInterface interface {
	doActivity(args []string) (string, error)
	setAction(action string)
	Action() string
	Description(action string) string

	getKeyboard() tgbotapi.InlineKeyboardMarkup
}

type userActivity struct {
	possible_actions []string
	actionType       string
	action           string
	description      map[string]string
}

func (ua userActivity) getKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, action := range ua.possible_actions {
		buttonData := strings.Split(action, " ")[0]
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(action, buttonData),
		)
		rows = append(rows, row)
	}
	// TODO: send success to user
	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func (ua *userActivity) setAction(action string) {

	ua.action = action
}

func (ua userActivity) Action() string {
	return ua.action
}

func (ua userActivity) Description(action string) string {
	return ua.description[action]
}

type eventActivity struct {
	userActivity
}

func (a eventActivity) String() string {
	return a.actionType
}

func NewEventActivity() *eventActivity {
	a := new(eventActivity)
	a.actionType = "event"
	a.possible_actions = []string{"sleep event", "pampers event", "eat event"}
	a.description = map[string]string{
		"sleep":   "<sleep_time>\nevent time in format hh:mm or YYYY-MM-DD_hh:mm",
		"eat":     "<eat_time> description\n<event_time> in format HH:MM, description - любой текст",
		"pampers": "<pampers_time> dirty\\wet\\combine\n<event_time> in format HH:MM и тип памперса",
	}
	return a
}

type babyActivity struct {
	userActivity
}

func (a babyActivity) String() string {
	return a.actionType
}

func NewBabyActivity() *babyActivity {
	a := new(babyActivity)
	a.actionType = "baby"
	a.possible_actions = []string{"add baby", "current baby", "remove baby"}
	a.description = map[string]string{
		"add":     "<name> <date_of_birth>\ndate_of_birth in format YYYY-MM-DD",
		"current": "<number_of_baby>",
		"remove":  "<number_of_baby>",
	}
	return a
}

type stateActivity struct {
	userActivity
}

func (a stateActivity) String() string {
	return a.actionType
}

func NewStateActivity() *stateActivity {
	a := new(stateActivity)
	a.actionType = "state"
	a.possible_actions = []string{"today state", "week state", "month state"}
	a.description = map[string]string{
		"today state": "",
		"week state":  "",
		"month state": "",
	}
	return a
}

func (ea eventActivity) doActivity(args []string) (string, error) {
	// args is <parent_id>;<time>;descr (для всех, кроме sleep)
	var result string
	if len(args) < 2 {
		return result, fmt.Errorf("error: args must be: %s", ea.Description(ea.Action()))
	}
	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v\n",
		ea.actionType, ea.action, args)

	parentId, err := strconv.Atoi(args[0])
	if err != nil {
		return result, fmt.Errorf("error on convert pid to string:\n%s", err)
	}

	if !CheckParentRegistred(int64(parentId)) {
		return result, fmt.Errorf("error, parent with id not existing")
	}

	baby, err := GetCurrentBaby(int64(parentId))
	if err != nil {
		return result, fmt.Errorf("error getting current baby: %s", err)
	}
	if baby == nil {
		return result, fmt.Errorf("current baby not found")
	}

	var timeString string
	// TODO: month len, after pm time
	switch len(args[1]) {
	case 5:
		timeString = fmt.Sprintf("%s_%s", time.Now().Format("2006-01-02"), args[1])
	case 15:
		timeString = args[1]
	default:
		return "", fmt.Errorf("time must be in format 2006-1-15:04 or 15:04")
	}
	eventTime, err := time.Parse("2006-1-02_15:04", timeString)
	if err != nil {
		return result, err
	}
	var dbEntity eventBaseWorker

	event := newEvent(baby.Id())
	event.SetStart(eventTime)

	switch ea.action {
	case "sleep":
		{
			sleep, err := GetNotEndedSleepForBaby(baby.Id())
			if err != nil {
				slogger.Errorw("error get not ended sleep",
					"action", ea.action,
					"error", err.Error(),
					"user", parentId)
				return "", err
			}
			if sleep != nil {
				sleep.SetEndTime(event.Start())
				dbEntity = sleep
			} else {
				dbEntity = newSleep(*event)
			}
		}
	case "eat":
		{
			eat := newEat(*event)
			if len(args) == 3 {
				eat.SetDescription(args[2])
			} else {
				eat.SetDescription("")
			}
			dbEntity = eat
		}
	case "pampers":
		{
			if len(args) < 3 {
				return "", fmt.Errorf("пожалуйста, укажите состояние памперса")
			}
			var pSt pampersState
			switch args[2] {
			case "dirty":
				{
					pSt = dirty
				}
			case "wet":
				{
					pSt = wet
				}
			case "combined":
				{
					pSt = combined
				}
			default:
				{
					return "", fmt.Errorf("wrong type of pampers. must be dirty\\wet\\combine")
				}
			}
			pampers := newPampers(*event)
			pampers.SetState(pSt)
			dbEntity = pampers
		}
	default:
		{
			slogger.Errorw("not setted aciton, but doAction",
				"user", parentId)
		}
	}
	if err := dbEntity.writeStructToBase(); err != nil {
		slogger.Errorw("error write new action to base",
			"error", err.Error(),
			"user", parentId)
		return "", err
	}
	resultString := fmt.Sprintf("%s in %s successfully added to base",
		ea.actionType, event.Start())
	slogger.Debugw("action writed to base",
		"action", ea.actionType,
		"time", event.Start(),
		"user", parentId)
	return resultString, err
}

func (ba babyActivity) doActivity(args []string) (string, error) {
	//first arg is parentId
	var result string
	if len(args) < 2 {
		return result, fmt.Errorf("error: args must be: %s", ba.Description(ba.Action()))
	}
	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v\n",
		ba.actionType, ba.action, args)

	parentId, err := strconv.Atoi(args[0])
	if err != nil {
		return result, fmt.Errorf("error on convert pid to string:\n%s", err)
	}
	if !CheckParentRegistred(int64(parentId)) {
		return result, fmt.Errorf("error, parent with id not existing")
	}

	switch ba.action {
	case "add":
		{
			//args is ParentId;<name>;<date_of_birth>(in format YYYY-MM-DD);,
			if len(args) < 3 {
				return result, fmt.Errorf("error, baby format must be <name> <birth_date>")
			}
			baby := newBaby()
			t, err := time.Parse("2006-01-02", args[2])
			if err != nil {
				fmt.Printf("Error on reading date")
				return result, err
			}

			baby.SetBirth(t)
			baby.SetParent(int64(parentId))
			baby.SetName(args[1])
			if err := baby.writeStructToBase(); err != nil {
				return result, fmt.Errorf("ошибка на записи ребенка в базу данных:\n%s", err)
			}
			result = fmt.Sprintf("Ребенок %s успешно записан в базу данных", baby)
			fmt.Println(result)
		}
	case "current":
		{
			//arg is ParentId, babyNumber
			if len(args) < 2 {
				return result, fmt.Errorf("ошибка, не указан номер ребенка")
			}
			babyNumber, err := strconv.Atoi(args[1])
			if err != nil {
				return result, err
			}
			babyes, err := GetBabyesByParent(int64(parentId))
			if err != nil {
				return result, err
			}
			if babyNumber > len(babyes) || babyNumber == 0 {
				return result, fmt.Errorf("ошибка, указан неправильный номер")
			}
			neededBaby := babyes[babyNumber-1]

			p := newParent()
			if err := p.readStructFromBase(int64(parentId)); err != nil {
				return result, err
			}
			if err := p.SetCurrentBaby(neededBaby.Id()); err != nil {
				return result, err
			}
			if err := p.writeStructToBase(); err != nil {
				return result, err
			}
			result = fmt.Sprintf("Ребенок %s установлен", neededBaby)
			fmt.Println(result)

		}
	case "remove":
		{
			if len(args) < 2 {
				return result, fmt.Errorf("ошибка, не указан номер ребенка")
			}
			//args is parentId;baby number
			babyes, err := GetBabyesByParent(int64(parentId))
			if err != nil {
				return result, err
			}
			babyNumber, err := strconv.Atoi(args[1])
			if err != nil {
				return result, err
			}
			if len(babyes) < babyNumber || babyNumber == 0 {
				return result, fmt.Errorf("ошибка, выбран не существующий ребенок")
			}

			currentBaby, err := GetCurrentBaby(int64(parentId))
			if err != nil {
				return result, err
			}

			p := newParent()
			if err := p.readStructFromBase(int64(parentId)); err != nil {
				return result, err
			}

			if currentBaby.Id() == p.CurrentBaby() {
				//if deleted current baby - set to 0
				p.SetCurrentBaby(0)
				if err := p.writeStructToBase(); err != nil {
					return result, err
				}
			}

			if err := RemoveBabyFromBase(babyes[babyNumber-1].Id()); err != nil {
				return result, err
			}

			result = fmt.Sprintf("Ребенок %s удален из базы данных", babyes[babyNumber-1])
			fmt.Println(result)

		}
	}
	return result, nil
}

func (sa stateActivity) doActivity(args []string) (string, error) {
	var result string
	if len(args) < 1 {
		return result, fmt.Errorf("error: not have parentid")
	}
	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v\n",
		sa.actionType, sa.action, args)

	parentId, err := strconv.Atoi(args[0])
	if err != nil {
		return result, fmt.Errorf("error on convert pid to string:\n%s", err)
	}

	if !CheckParentRegistred(int64(parentId)) {
		return result, fmt.Errorf("error, parent with id not existing")
	}

	baby, err := GetCurrentBaby(int64(parentId))
	if err != nil {
		return result, fmt.Errorf("error getting current baby: %s", err)
	}
	if baby == nil {
		return result, fmt.Errorf("current baby not found")
	}

	switch sa.action {
	case "today":
		{
		}
	case "week":
		{
		}
	case "month":
		{
		}
	}
	result = fmt.Sprintf("Result for %s state", sa.action)

	return result, nil
}
