package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "middleware.go: no authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "middleware.go: invalid authorization header")
		return
	}

	// parse token
	//fmt.Println("middleware: headerParts", headerParts[1])
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
	//fmt.Println("middleware.go: userId:", userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "middleware.go: no user id")
		return 0, errors.New("middleware.go: no user id")
	}
	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "middleware.go: invalid user id")
		return 0, errors.New("middleware.go: invalid user id")
	}
	return idInt, nil
}
