package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

const (
	userHeader  = "user_id"
	groupHeader = "group_id"
)

func (h *Handler) CheckIdentity(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	headerParts := strings.Split(auth, " ")
	if len(headerParts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid header"})
		c.AbortWithStatus(http.StatusUnauthorized)
		slog.Error("error with header", slog.Any("header", headerParts))
	}
	userId, group, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		slog.Error("error with id", slog.Any("error", err))
	}
	c.SetCookie(userHeader, fmt.Sprint(userId), 7200, "/", "localhost", false, false)
	c.SetCookie(groupHeader, group, 7200, "/", "localhost", false, false)
	c.Set(userHeader, fmt.Sprint(userId))
	fmt.Println(group)
	c.Set(groupHeader, group)
}
