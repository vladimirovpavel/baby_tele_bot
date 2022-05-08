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

var token string

func telegramBot() {
	var UA userActivityInterface

	//bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slogger.Error("error connecting to telegram service!")
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	slogger.Info("bot created and start polling")
	for update := range updates {
		if update.Message != nil {
			slogger.Debugw("--receive message",
				"message", update.Message.Text,
				"user", update.Message.From.ID,
				"username", fmt.Sprintf(
					"%s %s", update.Message.From.FirstName, update.Message.From.LastName))
			var msgToUser tgbotapi.MessageConfig

			if reflect.TypeOf(update.Message.Text).Kind() != reflect.String || update.Message.Text == "" {
				slogger.Debugw("message not text", "user", update.Message.From.ID)
				msgToUser = tgbotapi.NewMessage(update.Message.Chat.ID, "Use words in request")
				bot.Send(msgToUser)
			} else {
				parentId := int64(update.Message.From.ID)

				if update.Message.IsCommand() {
					slogger.Debugw("message is command",
						"command", update.Message.Command(),
						"user", update.Message.From.ID)
					switch update.Message.Command() {

					case "start":
						var text string
						_, err := RegisterNewParent(parentId, update.Message.From.LastName)
						if err != nil {
							slogger.Errorw("error register parent",
								"error", err.Error(),
								"user", update.Message.From.ID,
							)
							text = "Ошибка регистрации вас в системе. Свяжитель с администратором"
							break
						} else {
							text = fmt.Sprintf("hello, %s!\nВы можете:\n\tсоздавать и просматривать события с помощью /create_event\nуправлять детьми через /babyes_data\n\tполучить статистику через /get_state", update.Message.From.FirstName)
							slogger.Debugw("parent already registred",
								"user", update.Message.From.ID,
							)
						}
						msgToUser = tgbotapi.NewMessage(update.Message.Chat.ID, text)

					case "create_event":
						{
							var text string
							// получим текущего ребенка родителя, и все события за сегодня
							// для него
							currentBaby, err := GetCurrentBaby(parentId)
							if err != nil {
								slogger.Errorw("error get current baby",
									"error", err.Error(),
									"user", parentId)
								msgToUser = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении ребенка")
								break
							} else if currentBaby == nil {
								slogger.Errorw("try create event without current baby",
									"error", "no current baby set",
									"user", parentId)
								msgToUser = tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала выберите ребенка в /babyes_data")
								break
							} else {
								// TODO: переработал вывод - надо проверить что вышло
								events := GetEventsByBabyDate(currentBaby.Id(), time.Now())
								text = "События за сегодня:\n"
								for _, e := range events {
									text = fmt.Sprintf("%s\t-%s %s\n",
										text, e, e.GetSpecialValueString())
								}
								text = fmt.Sprintf("%s\nВыберите действие:", text)
								UA = NewEventActivity()
							}
							msgToUser = tgbotapi.NewMessage(update.Message.Chat.ID, text)
							msgToUser.ReplyMarkup = UA.getKeyboard()
						}
					case "get_state":
						{
							UA = NewStateActivity()
							msgToUser = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите, какая статистика вам нужна")
							msgToUser.ReplyMarkup = UA.getKeyboard()
						}
					case "babyes_data":
						{
							var text string
							UA = NewBabyActivity()
							babyes, err := GetBabyesByParent(parentId)
							if err != nil {
								slogger.Errorw("error get babyes by parent",
									"error", err.Error(),
									"user", parentId)
								break
							}
							currentBaby, err := GetCurrentBaby(parentId)
							if err != nil {
								slogger.Errorw("error get current baby",
									"error", err.Error(),
									"user", parentId)
								break
							}
							if len(babyes) == 0 {
								text = "Вы еще не зарегистрировали ребенка. Сделайте это, нажав на соответствующую кнопку"

							} else if currentBaby == nil || currentBaby.Id() == 0 {
								text = "Вам нужно выбрать ребенка"
							} else {
								text = "Ваши дети:"
								for counter, baby := range babyes {
									text = fmt.Sprintf("\n%s  %d\t%s", text, counter+1, baby)
									if baby.Id() == currentBaby.Id() {
										text = fmt.Sprintf("%s (выбран активным)", text)
									}
								}

							}

							msgToUser = tgbotapi.NewMessage(update.Message.Chat.ID, text)
							msgToUser.ReplyMarkup = UA.getKeyboard()

						}

					}
					if _, err := bot.Send(msgToUser); err != nil {
						slogger.Errorw("error sending message to user",
							"error", err.Error(),
							"user", parentId)
						continue
					}

					// если получаем ПРОСТО ТЕКСТ
				} else {
					slogger.Debugw("message its just text",
						"user", parentId)
					var msg tgbotapi.MessageConfig
					var text string
					if UA == nil {
						// текст может быть только передачей параметров для выбранного действия.
						// так что если действие не выбрано - ничего не делаем
						text = "Пожалуйста, выберите команду:\n/create_event\n/get_state\n/babyes_data"
						slogger.Debugw("not setted action - no use command",
							"user", parentId)
					} else {
						// если же действие выбрано, к тексту прикрепляем id пользователя
						// и отправляем в action.doAction
						if UA.Action() == "" {
							text = "Пожалуйста, выберите действие"
						} else {
							args := []string{strconv.Itoa(int(parentId))}
							args = append(args, strings.Split(update.Message.Text, " ")...)

							text, err = UA.doActivity(args)
							if err != nil {
								slogger.Errorw("do activity returns a error",
									"error", err.Error(),
									"user", parentId)
								text = err.Error()
							}
							UA = nil

						}

					}
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
					if _, err := bot.Send(msg); err != nil {
						slogger.Errorw("error sending message to user",
							"error", err.Error(),
							"user", parentId)
						continue
					}
				}
			}
		} else if update.CallbackQuery != nil {
			slogger.Debugw("message is callback",
				"action", update.CallbackQuery.Data,
				"user", update.CallbackQuery.From.ID)
			// TODO: PARSE FROM THIS

			var msg tgbotapi.MessageConfig
			callback := tgbotapi.NewCallback(
				update.CallbackQuery.ID,
				update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				fmt.Println(err)
				continue
			}
			var text string
			if UA != nil {
				UA.setAction(update.CallbackQuery.Data)

				text = fmt.Sprintf("Для действия %s.%s, введите:\n%s",
					UA, UA.Action(), UA.Description(UA.Action()))
			} else {
				slogger.Debugw("callback without checkec command",
					"user", update.CallbackQuery.From.ID)
				text = "Пожалуйста, сначала выберите команду\n"
			}
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			if _, err := bot.Send(msg); err != nil {
				slogger.Errorw("error sending message to user",
					"error", err.Error(),
					"user", update.CallbackQuery.From.ID)
				continue

			}

		}

	}
}
