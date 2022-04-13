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
// добавлена doActivity для Baby
// pampers, eat, sleep структуры и интефрейсы перенесены в events.go

// далее нужны do activiti для event, do activity для state

type userActivityInterface interface {
	doActivity(args ...string) (string, error)
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
	a.possible_actions = []string{"add event", "today events"}
	a.description = map[string]string{
		"add event":    "<event_time>\nevent time in format HH:MM",
		"today events": "",
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
	a.possible_actions = []string{"add baby", "current baby", "view baby", "remove baby"}
	a.description = map[string]string{
		"add baby":     "<name> <date_of_birth>\ndate_of_birth in format YYYY-MM-DD",
		"current baby": "<number_of_baby>",
		"remove baby":  "<number_of_baby",
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

func (ea eventActivity) doActivity(args ...string) (string, error) {
	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
		ea.actionType, ea.action, args)
	switch ea.action {
	case "add":
		{
		}
	case "today":
		{
		}

	}

	return "", nil
}
func (ba babyActivity) doActivity(args ...string) (string, error) {
	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
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
			splittedDate := strings.Split(args[1], "-")
			year, yerr := strconv.Atoi(splittedDate[0])
			month, merr := strconv.Atoi(splittedDate[1])
			day, derr := strconv.Atoi(splittedDate[2])

			if yerr != nil || merr != nil || derr != nil {
				fmt.Printf("Error on reading date")
				return result, fmt.Errorf("%s, %s, %s", yerr, merr, derr)
			}
			baby.SetBirth(time.Date(year, time.Month(month+1), day, 0, 0, 0, 0, time.UTC))

			baby.SetParent(int64(parentId))
			baby.SetName(args[0])
			if err := baby.writeStructToBase(); err != nil {
				fmt.Printf("Error on writing new baby to base:\n%s", err)
				return result, err
			}
			result = fmt.Sprintf("Baby %s successfully writed to base", baby)
			fmt.Println(result)
		}
	case "current":
		{
			//arg is ParentId
			currentBaby, err := GetParentCurrentBaby(int64(parentId))
			if err != nil {
				return result, err
			}
			result = fmt.Sprintf("Current baby is %s", currentBaby)
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
			if len(babyes) < babyNumber {
				return result, fmt.Errorf("error, checked not existing baby")
			}
			if err := removeBabyFromBase(babyes[babyNumber-1].Id()); err != nil {
				return result, err
			}
			result = fmt.Sprintf("Baby %s remove from base", babyes[babyNumber-1])
			fmt.Println(result)
			// TODO: set current to zero if it removed

		}
	}
	return result, nil
}

func (sa stateActivity) doActivity(args ...string) (string, error) {
	// validate args
	var result string

	fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
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
				// ЕСЛИ ПОЛУЧАЕМ КОМАНДУ
				if update.Message.IsCommand() {
					parentId := int64(update.Message.From.ID)
					switch update.Message.Command() {

					case "start":

						_, err := RegisterNewParent(parentId, update.Message.From.LastName)
						if err != nil {
							fmt.Println(err)
							continue
						}
						message := fmt.Sprintf("hello, %s!\nТы можешь:\n\tработать с событиями через /create_event\nуправлять детьми через /babyes_data\n\tполучить статистику через /get_state", update.Message.From.FirstName)
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, message)
						//msg.ReplyMarkup = mainMenuKeyboard

					case "create_event":
						{
							// TODO: можно сразу выводить действия за сегодня
							UA = NewEventActivity()
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие")
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
					switch UA.Action() {
					case "": //если никакой action не был выбран пользователем ранее
						{
							text = "Please, select the command:\n/create_event\n/get_state\n/babyes_data"
						}
					default: // если же и action был и данные пользователем переданы - выполняем
						text, err = UA.doActivity(update.Message.Text)
						if err != nil {
							fmt.Println(err)
						}

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
			UA.setAction(update.CallbackQuery.Data)
			description := UA.Description(UA.Action())
			//if action without descrition and additional data needed
			if description == "" {
				UA.doActivity("")
				continue
			} else {

				// TODO : проверить что именно приходит клиенту
				text := fmt.Sprintf("For %s.%s, enter data:\n%s",
					UA, UA.Action(), UA.Description(UA.Action()))
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
				if _, err := bot.Send(msg); err != nil {
					fmt.Println(err)
					continue
				}
			}

		}

	}
}
