package handler

import (
	"ToDoGo/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Данный файл содержит в себе функции для хендлеров, которые затрагивают сами списки(НЕ ИХ ЭЛЕМЕНТЫ)
// При добавлении нового хендлера

// @Summary Create todo list
// @Security ApiKeyAuth
// @Tags lists
// @Description create todo list
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body models.TodoList true "list info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	Userid, err := getUserId(c)
	if err != nil {
		return
	}
	var input models.ToDo
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("list.go: invalid input body: %s", err.Error()))
		return
	}

	//call service method
	id, err := h.services.TodoList.Create(Userid, input) // return created list id
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

type getListResponse struct {
	Data []models.ToDo `json:"data"`
}

// @Summary Get All Lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all lists
// @ID get-all-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllListsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	Userid, err := getUserId(c)
	if err != nil {
		return
	}
	lists, err := h.services.TodoList.GetAllLists(Userid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(lists)
	c.JSON(http.StatusOK, getListResponse{Data: lists})
}

// @Summary Get List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description get list by id
// @ID get-list-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ListItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id [get]
func (h *Handler) getListById(c *gin.Context) {
	Userid, err := getUserId(c)
	if err != nil {
		return
	}
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}
	list, err := h.services.TodoList.GetListById(Userid, i)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	Userid, err := getUserId(c)
	if err != nil {
		return
	}
	listId := c.Param("id")
	listIdInt, err := strconv.Atoi(listId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}
	var input models.ToDo
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("list.go: invalid input body: %s", err.Error()))
		return
	}

	updateList := h.services.TodoList.UpdateList(Userid, listIdInt, input)
	if updateList != nil {
		newErrorResponse(c, http.StatusInternalServerError, updateList.Error())
		c.JSON(http.StatusBadRequest, updateList.Error())
	} else {
		c.JSON(http.StatusOK, "success list update")
	}

}

func (h *Handler) deleteList(c *gin.Context) {
	Userid, err := getUserId(c)
	if err != nil {
		return
	}
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}
	dltList := h.services.TodoList.DeleteList(Userid, i)
	if dltList != nil {
		c.JSON(http.StatusBadRequest, dltList)
		newErrorResponse(c, http.StatusInternalServerError, dltList.Error())

	} else {
		c.JSON(http.StatusOK, "success list delete")
	}
}
