package ws

import (
	"fmt"
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/Futturi/Raspisanie/internal/service"
	"github.com/gorilla/websocket"
	"log/slog"
	"time"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *entities.Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

func (c *Client) writeMessage(service *service.Serivce) {
	defer func() {
		c.Conn.Close()
	}()
	for {
		message, ok := <-c.Message
		if !ok {
			slog.Error("error with msg")
			return
		}
		if message.ClientId == c.ID {
			err := service.InsertMessage(*message)
			fmt.Println(*message)
			if err != nil {
				slog.Error("error", slog.Any("erorr", err))
				return
			}
		}
		err := c.Conn.WriteJSON(*message)
		fmt.Println(message)
		if err != nil {
			slog.Error("error", slog.Any("erorr", err))
			return
		}
	}
}

func (c *Client) readMessage(hub *Hub, service *service.Serivce) {
	defer func() {
		err := service.DeleteUser(c.ID, c.RoomID)
		if err != nil {
			slog.Error("error", err)
		}
		hub.Unregister <- c
		c.Conn.Close()
	}()
	messages, err := service.GetAllMess(c.RoomID)
	if err != nil {
		slog.Error("error with reading saved message", slog.Any("error", err))
		return
	}

	err = c.Conn.WriteJSON(messages)
	if err != nil {
		slog.Error("error with writing messages", slog.Any("error", err))
		return
	}
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("error", slog.Any("error", err))
			}
			slog.Error("error", slog.Any("error", err))
			break
		}
		msg := &entities.Message{
			ClientId: c.ID,
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
			Data:     int(time.Now().Unix()),
		}
		hub.Broadcast <- msg
	}
}
