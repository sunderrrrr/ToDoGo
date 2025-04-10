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

func (h *Handler) passwordReset(c *gin.Context) {
	var input models.UserReset
	err := c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	response, err := h.services.User.ResetPassword(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": response})

}

func (h *Handler) deleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "user deleted"})
}
