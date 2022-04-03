package main

import (
	"fmt"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func telegramBot() {
	//bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	bot, err := tgbotapi.NewBotAPI("5164256009:AAEDA4OGpZTt2CPedXYS7Yn9dj9y86TLH_k")
	if err != nil {
		panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		fmt.Println(update)
		fmt.Println(update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		bot.Send(msg)
	}
	/* bot, err := tgbotapi.NewBotAPI(os.GETENV("TOKEN"))
	if err != nil {
		panic(err)
	} */

}
