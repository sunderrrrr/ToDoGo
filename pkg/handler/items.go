package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createItem(c *gin.Context) {

}

func (h *Handler) getAllItems(c *gin.Context) {
	id, _ := c.Get(userCtx)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) getItemById(c *gin.Context) {

}

func (h *Handler) updateItem(c *gin.Context) {

}

func (h *Handler) deleteItem(c *gin.Context) {

}
