package handler

import (
	"ToDoGo/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	tempId := c.Param("id")
	listId, err := strconv.Atoi(tempId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}
	var input models.TodoItem
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid input: %s", err.Error()))
	}
	id, err := h.services.TodoItem.CreateItem(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	tempId := c.Param("id")
	listId, err := strconv.Atoi(tempId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}
	items, err := h.services.TodoItem.GetAllItemsOfList(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("item-id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}
	item, err := h.services.TodoItem.GetItemById(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": item})
}

func (h *Handler) updateItem(c *gin.Context) {
	UserId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ListId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid list id: %s", err.Error()))
	}
	ItemId, err := strconv.Atoi(c.Param("item-id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid item id: %s", err.Error()))
		return
	}
	var input models.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("list.go: invalid input body: %s", err.Error()))
		return
	}

	updateList := h.services.TodoItem.UpdateItem(UserId, ListId, ItemId, input)
	if updateList != nil {
		newErrorResponse(c, http.StatusInternalServerError, updateList.Error())
		c.JSON(http.StatusBadRequest, updateList.Error())
	} else {
		c.JSON(http.StatusOK, "success item update")
	}

}

func (h *Handler) deleteItem(c *gin.Context) {
	UserId, err := getUserId(c)
	if err != nil {
		return
	}
	ItemId, err := strconv.Atoi(c.Param("item-id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}
	dltList := h.services.TodoItem.DeleteItem(UserId, ItemId)
	if dltList != nil {
		c.JSON(http.StatusBadRequest, dltList)
		newErrorResponse(c, http.StatusInternalServerError, dltList.Error())

	} else {
		c.JSON(http.StatusOK, "success list delete")
	}
}
