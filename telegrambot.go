package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: проблема работы с несколькими клиентами вообще. Нужно где-то учитывать статус пользвоателя

// EventActivity меняем - вывод информации по текущему дню можно сделать по умолчанию при
// выборе соответствующей команды. А inline клавиатуру сделаем позволяющей выбрать тип добавляемого события
// добавлена проверка уже идущего сна, которая в любом случае даст команду на завершение
// предыдущего, прежде чем можно будет ввести новый
// добавлено получение и вывод всех событий при выборе команды событий (с проверкой выборки
// актуального ребенка)

//Добавление нового события в список обрабатываемых:
//	1. вводим его описание в NewEventActivity
//	2. реализуем для него структуру в events, включая event
//	3. реализуем добавление в базу и таблицу в базе

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
	if len(args) < 2 {
		return "", fmt.Errorf("error: args must be: %s", ea.Description(ea.Action()))
	}
	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v\n",
		ea.actionType, ea.action, args)

	parentId, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	baby, err := GetCurrentBaby(int64(parentId))
	if err != nil {
		err := fmt.Errorf("error getting current baby: %s", err)
		return "", err
	}

	currentTimeString := args[1]

	// TODO: case когда данные полностью неправильные, но содержат _
	if !strings.ContainsAny(currentTimeString, "_") {
		// если указано только время - получим текущую дату
		tempTime := time.Now()
		tempTimeString := fmt.Sprintf("%d-%d-%d_%s",
			tempTime.Year(), tempTime.Month(), tempTime.Day(), currentTimeString)
		currentTimeString = tempTimeString
	}

	eventTime, err := time.Parse("2006-1-02_03:04", currentTimeString)
	if err != nil {
		return "", err
	}
	var dbEntity eventBaseWorker

	event := newEvent(baby.Id())
	event.SetStart(eventTime)

	switch ea.action {
	case "sleep":
		{
			sleep, err := GetNotEndedSleepForBaby(baby.Id())
			if err != nil {
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
			if len(args) > 2 {
				eat.SetDescription(args[2])
			}
			dbEntity = eat
		}
	case "pampers":
		{
			if len(args) < 3 {
				return "", fmt.Errorf("please, set type of pampers")
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
			case "combine":
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
	}
	if err := dbEntity.writeStructToBase(); err != nil {
		return "", err
	}
	resultString := fmt.Sprintf("%s in %s successfully added to base",
		ea.actionType, event.Start())
	return resultString, err
}

func (ba babyActivity) doActivity(args []string) (string, error) {
	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v\n",
		ba.actionType, ba.action, args)
	//first arg is parentId
	var result string
	parentId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error on convert pid to string:\n%s", err)
		return result, err
	}
	switch ba.action {
	case "add":
		{
			//args is ParentId;<name>;<date_of_birth>(in format YYYY-MM-DD);,
			baby := newBaby()
			if len(args) < 3 {
				return result, fmt.Errorf("error, baby format must be <name> <birth_date>")
			}
			t, err := time.Parse("2006-01-02", args[2])
			if err != nil {
				fmt.Printf("Error on reading date")
				return result, err
			}
			baby.SetBirth(t)

			baby.SetParent(int64(parentId))
			baby.SetName(args[1])
			if err := baby.writeStructToBase(); err != nil {
				fmt.Printf("Error on writing new baby to base:\n%s", err)
				return result, err
			}
			result = fmt.Sprintf("Baby %s successfully writed to base", baby)
			fmt.Println(result)
		}
	case "current":
		{
			if len(args) < 2 {
				return result, fmt.Errorf("error, not receive baby number")
			}
			//arg is ParentId, babyNumber
			babyNumber, err := strconv.Atoi(args[1])
			if err != nil {
				return result, err
			}
			if babyNumber <= 0 {
				return result, fmt.Errorf("Error: baby number must be greather then 0")
			}
			babyes, err := GetBabyesByParent(int64(parentId))
			if err != nil {
				return result, err
			}
			if len(babyes) < babyNumber {
				return result, fmt.Errorf("error, baby number more then babyes count")
			}
			neededBaby := babyes[babyNumber-1]

			p := newParent()
			if err := p.readStructFromBase(int64(parentId)); err != nil {
				return result, err
			}
			p.SetCurrentBaby(neededBaby.Id())
			if err := p.writeStructToBase(); err != nil {
				return result, err
			}
			result = fmt.Sprintf("Baby %s setted as current", neededBaby)
			fmt.Println(result)

		}
	case "remove":
		{
			if len(args) < 2 {
				return result, fmt.Errorf("error, not receive baby number")
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
			if babyNumber <= 0 {
				return result, fmt.Errorf("Error: baby number must be greather then 0")
			}
			if len(babyes) < babyNumber {
				return result, fmt.Errorf("error, checked not existing baby")
			}

			currentBaby, err := GetCurrentBaby(int64(parentId))
			if err != nil {
				return result, err
			} else if currentBaby == nil {
				return result, fmt.Errorf("At first, set current baby")
			}
			if currentBaby.Id() == babyes[babyNumber-1].Id() {
				p := newParent()
				if err := p.readStructFromBase(int64(parentId)); err != nil {
					return result, err
				}
				p.SetCurrentBaby(0)
				if err := p.writeStructToBase(); err != nil {
					return result, err
				}
			}

			if err := RemoveBabyFromBase(babyes[babyNumber-1].Id()); err != nil {
				return result, err
			}

			result = fmt.Sprintf("Baby %s remove from base", babyes[babyNumber-1])
			fmt.Println(result)

		}
	}
	return result, nil
}

func (sa stateActivity) doActivity(args []string) (string, error) {
	// validate args
	var result string

	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v\n",
		sa.actionType, sa.action, args)
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

var eventNumericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("7"),
		tgbotapi.NewKeyboardButton("8"),
		tgbotapi.NewKeyboardButton("9"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(":"),
		tgbotapi.NewKeyboardButton("0"),
		tgbotapi.NewKeyboardButton("."),
	),
)

func telegramBot() {
	var UA userActivityInterface

	//bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	bot, err := tgbotapi.NewBotAPI("5164256009:AAEDA4OGpZTt2CPedXYS7Yn9dj9y86TLH_k")
	if err != nil {
		panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			var msg tgbotapi.MessageConfig

			if reflect.TypeOf(update.Message.Text).Kind() != reflect.String || update.Message.Text == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use words in request")
				bot.Send(msg)
			} else {
				fmt.Println(update)
				fmt.Println(update.Message.Text)
				parentId := int64(update.Message.From.ID)
				// ЕСЛИ ПОЛУЧАЕМ КОМАНДУ
				if update.Message.IsCommand() {
					switch update.Message.Command() {

					case "start":

						_, err := RegisterNewParent(parentId, update.Message.From.LastName)
						if err != nil {
							fmt.Println(err)
							continue
						}
						message := fmt.Sprintf("hello, %s!\nВы можете:\n\tработать с событиями через /create_event\nуправлять детьми через /babyes_data\n\tполучить статистику через /get_state", update.Message.From.FirstName)
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, message)

					case "create_event":
						{
							var text string
							// получим текущего ребенка родителя, и все события за сегодня
							// для него
							currentBaby, err := GetCurrentBaby(parentId)
							if err != nil || currentBaby == nil {
								msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала выберите ребенка в /babyes_data")
								break
								// TODO: break из switch выходит?
							} else {
								events := GetEventsByBabyDate(currentBaby.Id(), time.Now())
								text = "События за сегодня:\n"
								for _, e := range events {
									text = fmt.Sprintf("%s\t- %s\n", text, e)
								}
								text = fmt.Sprintf("%s\nВыберите действие:", text)
								UA = NewEventActivity()
							}
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
							msg.ReplyMarkup = UA.getKeyboard()
						}
					case "get_state":
						{
							UA = NewStateActivity()
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите, какая статистика вам нужна")
							msg.ReplyMarkup = UA.getKeyboard()
						}
					case "babyes_data":
						{
							UA = NewBabyActivity()
							babyes, err := GetBabyesByParent(parentId)
							if err != nil {
								continue
							}
							var message string
							if len(babyes) == 0 {
								message = "Вы еще не зарегистрировали ребенка. Сделайте это, нажав на соответствующую кнопку"

							} else {
								// TODO: ВЫДЕЛЯТЬ АКТУАЛЬНОГО РЕБЕНА
								message = fmt.Sprintf("Ваши дети:\n")
								for counter, baby := range babyes {
									message = fmt.Sprintf("%s  %d\t%s", message, counter+1, baby)
								}

							}

							msg = tgbotapi.NewMessage(update.Message.Chat.ID, message)
							msg.ReplyMarkup = UA.getKeyboard()

						}

					}
					if _, err := bot.Send(msg); err != nil {
						fmt.Println(err)
						continue
					}

					// если получаем ПРОСТО ТЕКСТ
				} else {
					var msg tgbotapi.MessageConfig
					var text string
					if UA == nil {
						text = "Please, select the command:\n/create_event\n/get_state\n/babyes_data"
					} else {
						// если же и action был и данные пользователем переданы - выполняем
						// для все-таки, видимо, нужен слайс строк с ктнтролем целостности уже внутри
						// прямо сейчас проблема - при попытке выставить актуального ребенка
						// передаем больше аргументов, чем есть ( а есть - один, номер ребенка)
						args := []string{strconv.Itoa(int(parentId))}
						args = append(args, strings.Split(update.Message.Text, " ")...)

						text, err = UA.doActivity(args)
						if err != nil {
							fmt.Println(err)
							text = err.Error()
						}

						UA = nil

					}
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
					if _, err := bot.Send(msg); err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
		} else if update.CallbackQuery != nil {
			var msg tgbotapi.MessageConfig
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				fmt.Println(err)
				continue
			}
			var text string
			if UA != nil {
				UA.setAction(update.CallbackQuery.Data)
				//description := UA.Description(UA.Action())
				//if action without descrition and additional data needed

				// TODO : проверить что именно приходит клиенту
				text = fmt.Sprintf("For %s.%s, enter data:\n%s",
					UA, UA.Action(), UA.Description(UA.Action()))
			} else {
				text = "Please, at start check command\n"
			}
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err)
				continue

			}

		}

	}
}
