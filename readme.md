# REST Api на Golang с использованием внедрения зависимостей 🚀

## Принцип работы
Приложение реализует принцип внедрения зависимостей (Луковичная архитектура или Clean Architecture), который позволяет расширять функционал, не затрагивая основные элементы кода.

Данные, получаемые от пользователя, проходят через 4 уровня:
* Веб-сервер 🌐
* Обработчик 🔄
* Сервис ⚙️
* Репозиторий 🗄️

## Стек технологий
### Backend
* ЯП - Golang 🐹
* Gin Framework 🏗️
* Docker 🐳
* Postgres 🗃️

### Frontend
* ReactJs ⚛️
* Material3 🎨

## Развертывание
1. Установить Golang
2. Выполнить `go mod tidy`
3. Установить Docker
4. Выполнить следующую команду для создания контейнера Postgres:
   ```
   docker run --name=todo-db -e POSTGRES_PASSWORD=qwerty -p 5436:5432 -d postgres
   ```
    * `name` - имя контейнера
    * `POSTGRES_PASSWORD` - пароль от базы данных Postgres 🔑
    * `-p` - порты для Postgres
5. Установить migrate
6. Выполнить миграцию:
   ```
   migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up
   ```
    * `qwerty` - пароль для базы данных
    * `5436` - порт Postgres
7. Скопировать файл .env.example, переименовать его в .env и заполнить
8. Запуск сервера `go run cmd/main.go`

## Документация (Полная документация /requests) 📖
### Регистрация
```
POST http://localhost:8090/auth/sign-up
Content-Type: application/json

{
  "name": "iliya22",
  "username": "bkmz1153434311",
  "password": "qwerty"
}
```

### Авторизация (Возвращает Bearer token) 🔑
```
POST http://localhost:8090/auth/sign-in
Content-Type: application/json

{
  "username": "bkmz1153434311",
  "password": "qwerty"
}
```

### Создание списка 📋
```
POST http://localhost:8090/api/lists/
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDA2Mzc3MTcsImlhdCI6MTc0MDU5NDUxNywidXNlcl9pZCI6MTB9.d4NgeHXl9zT_bXG9Ad-NvvC49MI892SiMdbF5lw4G4I

{
  "title": "затащит физтех",
  "description": "1 марта"
}
```