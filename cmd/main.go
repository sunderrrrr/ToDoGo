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
	// Инициализация конфигурации приложения
	if initConfig() != nil {
		log.Fatalf("config initializing failed: %s", initConfig().Error())
	}

	// Загрузка переменных окружения из файла .env
	if godotenv.Load(".env") != nil {
		log.Fatalf("env initializing failed: %s", godotenv.Load(".env").Error())
	}

	// Создание соединения с базой данных PostgreSQL с использованием конфигурации
	db, err := repository.NewPostgresDB(repository.ConnConfig{
		Host:     viper.GetString("database.host"),     // Хост базы данных
		Port:     viper.GetString("database.port"),     // Порт базы данных
		Username: viper.GetString("database.username"), // Имя пользователя базы данных
		Password: os.Getenv("DB_PASSWORD"),             // Пароль базы данных из переменных окружения
		Database: viper.GetString("database.database"), // Имя базы данных
		SSLMode:  viper.GetString("database.sslmode"),  // Режим SSL для базы данных
	})

	// Проверка на наличие ошибок при инициализации базы данных
	if err != nil {
		log.Fatalf("database initializing failed: %s", err.Error())
	}

	// Логирование успешной инициализации базы данных
	log.Default().Println("DB init done")

	// Создание нового репозитория на основе соединения с базой данных
	repos := repository.NewRepository(db)
	log.Default().Println("Repository init done")

	// Создание нового сервиса на основе репозитория
	services := service.NewService(repos)
	log.Default().Println("Service init done")

	// Создание нового обработчика на основе сервиса
	handlers := handler.NewHandler(services)
	log.Default().Println("Handler init done")

	// Создание нового экземпляра сервера
	srv := new(ToDoGo.Server)

	// Запуск сервера и инициализация маршрутов
	err = srv.Run(viper.GetString("srvport"), handlers.InitRoutes())

	// Проверка на ошибки при запуске сервера
	if err != nil {
		log.Fatalf("error while starting server: %e", err.Error())
	} else {
		log.Println("Start Server Success") // Успешный запуск сервера
	}

	log.Default().Println("DB done") // Завершение работы с базой данных
}

// Функция для инициализации конфигурации приложения
func initConfig() error {
	viper.AddConfigPath("configs") // Установка пути к конфигурационным файлам
	viper.SetConfigName("config")  // Установка имени конфигурационного файла
	return viper.ReadInConfig()    // Чтение конфигурации из файла
}
