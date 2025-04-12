package handler

import (
	"ToDoGo/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getUserInfo(c *gin.Context) {
	username, err := getUsername(c)
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if username == "" {
		newErrorResponse(c, http.StatusBadRequest, "username is empty")
	}
	c.JSON(http.StatusOK, gin.H{"id": id,
		"name": username})
}

func (h *Handler) passwordResetRequest(c *gin.Context) { //Запрос сброса пароля
	var input models.ResetRequest
	err := c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	err = h.services.Authorization.ResetPasswordRequest(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}
func (h *Handler) passwordResetConfirm(c *gin.Context) {
	var input models.UserReset
	err := c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	err = h.services.Authorization.ResetPassword()

}
func (h *Handler) deleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "user deleted"})
}
