package models

type User struct { //Структура пользователя
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required" db:"password_hash"`
}

type UserReset struct {
	Username string `json:"username" binding:"required"`
	Token    string `json:"token" binding:"required"`
	OldPass  string `json:"old_password" binding:"required"`
	NewPass  string `json:"new_password" binding:"required"`
}

type ResetRequest struct {
	Login string `json:"login" binding:"required"`
}
