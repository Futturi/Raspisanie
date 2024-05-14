package handler

import (
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) GetRasp(c *gin.Context) {
	gr := c.GetString(groupHeader)
	var group entities.Group
	err := c.BindJSON(&group)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "incorrect body",
		})
		slog.Error("error with body", slog.Any("error", err))
		return
	}

	result, err := h.service.GetRasp(group, gr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": "smth wrong",
		})
		slog.Error("error", slog.Any("error", err))
	}
	c.JSONP(http.StatusOK, result)
}
