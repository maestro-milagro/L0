package main

import (
	message "awesomeProject"
	"awesomeProject/pkg/cache"
	"awesomeProject/pkg/handler"
	"awesomeProject/pkg/repository"
	"awesomeProject/pkg/service"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

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

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(message.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print("MessageApp Started")
	list, err := services.MessagesS.GetAll()
	c := cache.C
	for _, v := range list {
		_, ok := c.Read(v.MessageId)
		if ok {
			continue
		}
		c.Update(v.MessageId, v)
	}
	sc, _ := stan.Connect("mess", "sub")

	sc.Subscribe("message", func(m *stan.Msg) {
		var response message.Message
		err = json.Unmarshal(m.Data, &response)
		if err != nil {
			fmt.Printf("error occured while subing: %s", err.Error())
		}
		if reflect.DeepEqual(response, message.Message{}) {
			fmt.Println("invalid message")
		} else {
			id, err := services.MessagesS.Create(response)
			if err != nil {
				fmt.Printf("message error: %s", err.Error())
			}
			fmt.Printf("order whith id : %d was created", id)
			fmt.Println()
		}
	})
	if err != nil {
		return
	}
	Block()
	sc.Close()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("MessageApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
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
