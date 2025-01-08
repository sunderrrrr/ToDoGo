package main

import (
	"ToDoGo"
	"ToDoGo/pkg/handler"
	"ToDoGo/pkg/repository"
	"ToDoGo/pkg/service"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if initConfig() != nil {
		log.Fatalf("config initializing failed: %s", initConfig().Error())
	}
	db, err := repository.NewPostgresDB(repository.ConnConfig{})
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(ToDoGo.Server)
	err = srv.Run(viper.GetString("port"), handlers.InitRoutes())
	if err != nil {
		log.Fatalf("error while starting server: %e", err.Error())
	} else {
		log.Println("Start Server Success")
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
