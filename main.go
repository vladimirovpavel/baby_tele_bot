package main

import "go.uber.org/zap"

var logger = zap.NewExample()
var slogger = logger.Sugar()

/*
TODO:
3. реализовать выборку по дате и ребенку для всех событий
	3.1. тут можно сделать максимальную выборку select * from ... в event, дергать ее
	3.1.1. потестить написанную selectbybabydate
	3.2. а потом уже парсить результат в каждом родителе в зависимости от типа
8. для всех структур реализовать метод stirng для выдачи клиенту
10. для запроса события понадобится второй метод - можно в eventI - readEventSpecial, который будет считывать все специфичное для конкретного эвента

*/

func main() {

	//CreatingData()
	slogger.Infof("create telegrambot")
	telegramBot()
	// создал RegisterNewParent с возвратом родителя (созданного, или существующего)
	// создал GetBabyesByParent c возвратом всех детей родителя
	// изменил new, writetobase и readfrombase для родителя и ребенка.
	//		теперь новый - принимает ID.
	//		read - считывает все поля из базы по существующему id
	//		write - записывает все поля в базу
	// описывать функции для действий пользователя в telegrambot.go

}
