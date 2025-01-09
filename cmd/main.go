package main

import (
	"ToDoGo"
	"ToDoGo/pkg/handler"
	"ToDoGo/pkg/repository"
	"ToDoGo/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	if initConfig() != nil {
		log.Fatalf("config initializing failed: %s", initConfig().Error())
	}
	log.Default().Println("Config init  done")
	if godotenv.Load(".env") != nil {
		log.Fatalf("env initializing failed: %s", godotenv.Load(".env").Error())
	}
	db, err := repository.NewPostgresDB(repository.ConnConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("database.database"),
		SSLMode:  viper.GetString("database.sslmode"),
		//Host: "localhost", Port: "5436", Username: "postgres", Password: "qwerty", Database: "postgres", SSLMode: "disable",
	})

	if err != nil {
		log.Fatalf("database initializing failed: %s", err.Error())
	}
	log.Default().Println("DB init done")
	repos := repository.NewRepository(db)
	log.Default().Println("Repository init done")
	services := service.NewService(repos)
	log.Default().Println("Service init done")
	handlers := handler.NewHandler(services)
	log.Default().Println("Handler init done")
	srv := new(ToDoGo.Server)
	err = srv.Run(viper.GetString("port"), handlers.InitRoutes())
	if err != nil {
		log.Fatalf("error while starting server: %e", err.Error())
	} else {
		log.Println("Start Server Success")
	}
	log.Default().Println("DB done")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
