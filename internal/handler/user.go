package handler

import (
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) GetUser(c *gin.Context) {
	id := c.GetString(userHeader)
	gr, err := h.service.GetUser(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "smth wrong",
		})
		slog.Error("error", err)
		return
	}
	c.JSON(http.StatusOK, gr)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.GetString(userHeader)
	var us entities.UpdateUser
	if err := c.BindJSON(&us); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "incorrect body",
		})
		slog.Error("error")
		return
	}
	err := h.service.UpdateUser(id, us)
	if err != nil {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "smth wrong",
			})
			slog.Error("error", err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
