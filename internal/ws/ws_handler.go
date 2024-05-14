package ws

import (
	"fmt"
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/Futturi/Raspisanie/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"strings"
)

const (
	userHeader  = "user_id"
	groupHeader = "group_id"
)

type Handler struct {
	hub     *Hub
	service *service.Serivce
}

func NewHandler(hub *Hub, service *service.Serivce) *Handler {
	return &Handler{hub: hub, service: service}
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req entities.CreateRoomReq
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "smth wrong"})
		slog.Error("error with binding body", slog.Any("error", err))
		return
	}
	req.Group = c.GetString("group_id")
	id, err := h.service.CreateRoom(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "smth wrong"})
		slog.Error("error with creating room", slog.Any("error", err))
		return
	}
	h.hub.Rooms[id] = &Room{
		ID:      id,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}
	c.JSON(http.StatusOK, req)
}

var upr = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) CheckIdentity(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		c.AbortWithStatus(http.StatusUnauthorized)
		slog.Error("error with token, its null")
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
	c.Set(userHeader, fmt.Sprint(userId))
	c.Set(groupHeader, group)
}

// коннект к руме в параме /api/join/рума
func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upr.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("error with establish conn", slog.Any("error", err))
		return
	}
	h.CheckIdentity(c)
	fmt.Println(c.GetString(userHeader))
	room := c.Param("roomId")
	roomId, err := h.service.GetRoomId(room)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "smth wrong")
		slog.Error("error", err)
		return
	}
	clientId := c.GetString(userHeader)
	fmt.Println(clientId, c.GetString("group_id"))
	username, err := h.service.GetUsername(clientId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("error with finding name", slog.Any("error", err))
		return
	}
	cl := &Client{
		Conn:     conn,
		Message:  make(chan *entities.Message, 10),
		ID:       clientId,
		RoomID:   roomId,
		Username: username,
	}
	h.hub.Register <- cl
	err = h.service.InsertUser(clientId, roomId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("error with inserting user name", slog.Any("error", err))
		return
	}
	go cl.writeMessage(h.service)
	cl.readMessage(h.hub, h.service)
}

func (h *Handler) GetRooms(c *gin.Context) {
	group_id := c.GetString("group_id")
	rooms, err := h.service.GetRooms(group_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "smth wrong"})
		slog.Error("error", slog.Any("error", err))
	}
	for _, v := range rooms {
		h.hub.Rooms[v.Id] = &Room{
			ID:      v.Id,
			Name:    v.Name,
			Clients: make(map[string]*Client),
			Count:   v.Count,
		}
	}
	c.JSON(http.StatusOK, rooms)
}
