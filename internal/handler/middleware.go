package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := h.services.Auth.ParseJWT(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("userId", userId)

	if _, err = getUserName(c); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
}

func (h *Handler) logging(c *gin.Context) {
	h.logger.Info(c.Request.URL.Path)
}

func getUserName(c *gin.Context) (int, error) {
	userId, ok := c.Get("userId")
	if !ok {
		return 0, errors.New("username not found")
	}

	return userId.(int), nil
}

func getId(c *gin.Context) (int, error) {
	strId, ok := c.GetQuery("id")
	if ok != true {
		return 0, errors.New("id not found")
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		return 0, errors.New("id is of invalid type")
	}

	return id, nil
}
