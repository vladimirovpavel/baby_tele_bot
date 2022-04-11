package main

import (
	"fmt"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// добавлен интерфейс userActivityInterface и конструктор клавиатур
// в итоге
// добавлены структуры активностей пользователя
// добавлен вывод детей для подменю детей

// начата реализация doAction. Проверить корректность работы - на каждую команду
// дожлна появляться inline клавиатура. На каждое нажатие кнопок на ней
// должен отрабаывать соответствующий doAction

// далее надо реализовывать doAction с реальной записью или запросом в\из db

/*
const (
	stableState tgBotState = iota
	waitingEvent
	waitingBabyAdd
	waitingBabySet
	waitingBabyRemove
	waitingState
)
*/
type userActivityInterface interface {
	doActivity(args string) error
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

		//buttonData := fmt.Sprintf("%s_%s", ua.actionType, strings.Split(action, " ")[0])
		buttonData := strings.Split(action, " ")[0]
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(action, buttonData),
		)
		rows = append(rows, row)
	}
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
		"view baby":    "",
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

func (ea eventActivity) doActivity(args string) error {
	switch ea.action {
	case "add":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				ea.actionType, ea.action, args)
		}
	case "today":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				ea.actionType, ea.action, args)
		}

	}

	return nil
}
func (ba babyActivity) doActivity(args string) error {
	switch ba.action {
	case "add":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				ba.actionType, ba.action, args)
		}
	case "current":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				ba.actionType, ba.action, args)
		}
	case "remove":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				ba.actionType, ba.action, args)
		}
	case "view":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				ba.actionType, ba.action, args)
		}
	}
	// validate args

	return nil
}
func (sa stateActivity) doActivity(args string) error {
	// validate args
	switch sa.action {
	case "today":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				sa.actionType, sa.action, args)
		}
	case "week":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				sa.actionType, sa.action, args)
		}
	case "month":
		{
			fmt.Printf("in \"doActivity\" method of action %s. Action is %s, args %#v",
				sa.actionType, sa.action, args)
		}
	}

	return nil
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
								message = fmt.Sprintf("Ваши дети:\n%s\n", babyes)
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
						err := UA.doActivity(update.Message.Text)
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
			// TODO: switch for callbacks
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

				text := fmt.Sprintf("For %s.%s, enter data:\n%s",
					UA, UA.Action(), UA.Description(UA.Action()))
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
				if _, err := bot.Send(msg); err != nil {
					fmt.Println(err)
					continue
				}
			}
			/*callback_strings := strings.Split(update.CallbackData(), "_")
			UA.setAction(callback_strings[1])
			switch callback_strings[0] {
			case "event":
				{
					// если опция add - ждем ввода данных от пользователя
					// если опция today - делаем запрос событий за сегодня и выводим
					text := fmt.Sprintf("Type of option is %s, option is %s",
						callback_strings[0], callback_strings[1])
					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
					currentBotState = waitingEvent
				}
			case "state":
				{
					// тут в зависимости от опции просто дергаем функции
					// той или иной выборки и подсчета
					text := fmt.Sprintf("Type of option is %s, option is %s",
						callback_strings[0], callback_strings[1])
					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
					currentBotState = waitingState

				}
			case "baby":
				{
					//опция add - ждем вывода данных от пользователя
					//опция set - ждем вывода данных от пользователя
					//оция view - запрашиваем данные и выводим
					//опция remove - ждем ввода данных от пользователя
					text := fmt.Sprintf("Type of option is %s, option is %s",
						callback_strings[0], callback_strings[1])
					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
					switch callback_strings[1] {
					case "add":
						{
							currentBotState = waitingBabyAdd
						}
					case "set":
						{
							currentBotState = waitingBabySet
						}
					case "view":
						{
						}
					case "remove":
						{
							currentBotState = waitingBabyRemove
						}
					}

				}*/
			// ввод данных - в 		event add
			//						babyes add
			//						babyes set
			//						babyes remove

		}

	}
}
