package main

import (
	message "awesomeProject"
	"awesomeProject/pkg/cache"
	"awesomeProject/pkg/handler"
	"awesomeProject/pkg/repository"
	"awesomeProject/pkg/service"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reflect"
	"sync"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	//Инициализируем конфигурационный файл config.yml,
	//чтобы получить возможность получать из него данные для подключения к бд
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()
	//Создаем объект базы данных
	repos := repository.NewRepository(db)
	//Создаем объект сервиса
	services := service.NewService(repos)
	//Создаем объект хендлера
	handlers := handler.NewHandler(services)
	//Запускаем http сервер
	srv := new(message.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print("MessageApp Started")
	//Чтобы избежать потери данных кэша при подении сервиса проверяем их наличие в кэше
	//В случае отсутствия каких-либо данных восстанавливаем их из бд
	list, err := services.MessagesS.GetAll()
	c := cache.C
	for _, v := range list {
		_, ok := c.Read(v.MessageId)
		if ok {
			continue
		}
		c.Update(v.MessageId, v)
	}
	//Подключаемся к каналу nats-streaming
	sc, _ := stan.Connect("mess", "sub")
	//Подписываемся на канал nats-streaming
	sc.Subscribe("message", func(m *stan.Msg) {
		//Декодируем полученный из канала массив байт в объект Message
		var response message.Message
		err = json.Unmarshal(m.Data, &response)
		if err != nil {
			fmt.Printf("error occured while subing: %s", err.Error())
		}
		//Проверяем подходит ли полученный объект под структуру Message
		//Если мы получили объект в котором ни одно поле не совпадает с нужной нам структурой выводим сообщение об ошибке
		if reflect.DeepEqual(response, message.Message{}) {
			fmt.Println("invalid message")
		} else {
			//В обратном случае сохраняем объект в бд
			id, err := services.MessagesS.Create(response)
			if err != nil {
				fmt.Printf("message error: %s", err.Error())
			}
			fmt.Printf("order whith id : %d was created", id)
			fmt.Println()
		}
		//Durable Subscription дает возможность не пропускать сообщения при падении сервиса
		//и продолжить чтение с последнего прочитанного сообщения
	}, stan.DurableName("my-durable"))
	Block()
	sc.Close()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
func Block() {
	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}
