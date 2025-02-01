package handler

import (
	"ToDoGo/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.User                      // получаем на вход структуру User и валидируем запрос
	if err := c.BindJSON(&input); err != nil { // записываем данные из JSON по сссылке input'а
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Authorization.CreateUser(input) // отправляем введенные данные на уровень сервиса
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // если ответ кривой возвращаем ошибку
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id}) // иначе возвращаем айдишник пользователя
}

type signInInput struct { //Кастом структура, так как структура User подразумевает обязательное имя
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil { // если тело запроса кривое, то возвращаем Bad Request Error
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password) // если все на кондициях то отправляем логи на уровень сервиса
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // обрабатываем ошибку ответа
		fmt.Println("213423423423")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token}) // возвращаем токен

}
