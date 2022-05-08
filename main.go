package main

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// repalces assert to require
// добавил метод GetSpecialValue к каждому из конкретных событий
// исправил в DBReadRow - возврат ошибки из result, если она происходит
// исправил GetTypedEventByBabyDate - поставил верхнюю границу для запроса по дате
// исправил вывод активного ребенка

//var logger = zap.NewDevelopmentConfig()

// сделать докер контейнер для baby-tracker с передачей токена бота через env
// сменить токен бота
// задеплоить контейнет в базой данных на сервер
// настроить ufw на сервере для правильного проброса порта
// вставить логирование

/*
TODO:
3. реализовать выборку по дате и ребенку для всех событий
3.1. тут можно сделать максимальную выборку select * from ... в event, дергать ее
3.1.1. потестить написанную selectbybabydate
3.2. а потом уже парсить результат в каждом родителе в зависимости от типа
8. для всех структур реализовать метод stirng для выдачи клиенту
10. для запроса события понадобится второй метод - можно в eventI - readEventSpecial, который будет считывать все специфичное для конкретного эвента


*/
var slogger *zap.SugaredLogger

func init() {
	var conf zap.Config
	runType := os.Getenv("RUN_TYPE")
	switch runType {
	case "develop":
		fmt.Println("develop logging")
		conf = zap.NewDevelopmentConfig()
	case "production":
		fmt.Println("product logging")
		conf = zap.NewProductionConfig()
	default:
		fmt.Println("nog found env, basic logging is development")
		conf = zap.NewDevelopmentConfig()
	}
	conf.OutputPaths = append(conf.OutputPaths, "/var/tmp/telebot.log")
	logger, err := conf.Build()
	if err != nil {
		return
	}
	slogger = logger.Sugar()

	token = os.Getenv("TOKEN")
	if token == "" {
		token = "5164256009:AAG7l3jCAu-WmexXKFh1TObeE43G3IEsvy4"
	}
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if user == "" || pass == "" || db == "" || host == "" || port == "" {
		fmt.Println("Using standart dbinfo")
		dbInfo = "host=127.0.0.1 port=5432 user=postgres password=!QAZxsw2 dbname=test_db"
	} else {
		fmt.Println("Using dbinfo from envs")
		dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			host, port, user, pass, db)
	}

}

func main() {
	//CreatingData()
	slogger.Infof("create telegrambot")
	time.Sleep(5 * time.Second)
	telegramBot()
	// создал RegisterNewParent с возвратом родителя (созданного, или существующего)
	// создал GetBabyesByParent c возвратом всех детей родителя
	// изменил new, writetobase и readfrombase для родителя и ребенка.
	//		теперь новый - принимает ID.
	//		read - считывает все поля из базы по существующему id
	//		write - записывает все поля в базу
	// описывать функции для действий пользователя в telegrambot.go

}
