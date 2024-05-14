package handler

import (
	"github.com/Futturi/Raspisanie/internal/service"
	"github.com/Futturi/Raspisanie/internal/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler struct {
	service *service.Serivce
}

func NewHandler(service *service.Serivce) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(wsHan *ws.Handler) http.Handler {
	serv := gin.Default()

	serv.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	auth := serv.Group("/auth")
	{
		auth.POST("/signin", h.SignIn)
		auth.POST("/signup", h.SignUp)
	}
	api := serv.Group("/api", h.CheckIdentity)
	{
		api.POST("/", h.GetRasp)
		api.POST("/createroom", wsHan.CreateRoom)
		api.GET("/joinroom/:roomId", wsHan.JoinRoom)
		api.GET("/rooms", wsHan.GetRooms)
	}
	return serv
}
