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

//Добавление нового события в список обрабатываемых:
//	1. вводим его описание в NewEventActivity
//	2. реализуем для него структуру в events, включая event
//	3. реализуем добавление в базу и таблицу в базе

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
	slogger.Info("bot created and start polling")
	for update := range updates {
		slogger.Infof("receive message %s from user %s",
			update.Message.Text,
			update.Message.From.ID)
		if update.Message != nil {
			var msg tgbotapi.MessageConfig

			if reflect.TypeOf(update.Message.Text).Kind() != reflect.String || update.Message.Text == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use words in request")
				bot.Send(msg)
			} else {
				fmt.Println(update)
				fmt.Println(update.Message.Text)
				parentId := int64(update.Message.From.ID)

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
								fmt.Println(err)
								continue
							}
							currentBaby, err := GetCurrentBaby(parentId)
							if err != nil {
								fmt.Println(err)
								continue
							}
							var message string
							if len(babyes) == 0 {
								message = "Вы еще не зарегистрировали ребенка. Сделайте это, нажав на соответствующую кнопку"

							} else {
								message = fmt.Sprintf("Ваши дети:\n")
								for counter, baby := range babyes {
									if baby.Id() == currentBaby.Id() {
										message = fmt.Sprintf("%s  %d\t**__%s__**", message, counter+1, baby)

									}
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
