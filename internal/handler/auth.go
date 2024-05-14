package handler

import (
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) SignUp(c *gin.Context) {
	var user entities.User

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "smth wrong",
		})
		slog.Error("error with request", slog.Any("error", err))
		return
	}
	id, err := h.service.SignUp(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": "smth wrong",
		})
		slog.Error("error with request", slog.Any("error", err))
		return
	}
	c.JSON(http.StatusOK, map[string]int{
		"id": id,
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	var user entities.User

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "smth wrong",
		})
		slog.Error("error with request", slog.Any("error", err))
		return
	}

	token, err := h.service.SignIn(user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": "smth wrong",
		})
		slog.Error("error with request", slog.Any("error", err))
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
