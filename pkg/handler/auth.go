package handler

import (
	"ToDoGo/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
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

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil { // если тело запроса кривое, то возвращаем Bad Request Error
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password) // если все на кондициях то отправляем логи на уровень сервиса
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // обрабатываем ошибку ответа
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token}) // возвращаем токен

}
