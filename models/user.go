package models

type User struct { //Структура пользователя
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required" db:"password_hash"`
}

type UserReset struct {
	OldPass string `json:"old_password" binding:"required"`
	NewPass string `json:"new_password" binding:"required"`
}
