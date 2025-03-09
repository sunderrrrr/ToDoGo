package models

import "errors"

type ToDo struct { //Структура списка
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" binding:"required" db:"title"`
	Description string `json:"description" db:"description"`
}

type UserList struct { //Стукрутра, связывающая пользователя и список через ID
	Id     int
	UserId int
	ListId int
}

type TodoItem struct { //Структура члена списка
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type ListsItem struct { //Структура связывающая список и член списка
	Id     int
	ListId int
	ItemId int
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateItemInput) Valid() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("missing fields error")
	} else {
		return nil
	}
}
