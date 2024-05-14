package repository

import (
	"fmt"
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"strconv"
)

type WSRepository struct {
	db *sqlx.DB
}

func NewWSRepository(db *sqlx.DB) *WSRepository {
	return &WSRepository{db: db}
}

func (r *WSRepository) CreateRoom(req entities.CreateRoomReq) (string, error) {
	var id string
	query := "INSERT INTO rooms(name, groupa) VALUES ($1, $2) RETURNING id"
	row := r.db.QueryRow(query, req.Name, req.Group)
	if err := row.Scan(&id); err != nil {
		slog.Error("error", err)
		return "", err
	}
	return id, nil
}

func (r *WSRepository) InsertMessage(msg entities.Message) error {
	var id int
	query := "INSERT INTO messages(text, creator, date) VALUES($1,$2,$3) RETURNING id"
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	row := tx.QueryRow(query, msg.Content, msg.Username, msg.Data)
	if err = row.Scan(&id); err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			slog.Error("error", slog.Any("erorr", err1))
			return err1
		}
		slog.Error("error", slog.Any("erorr", err))
		return err
	}
	query2 := "INSERT INTO rooms_messages(room_id, message_id) VALUES ($1, $2)"
	if _, err = tx.Exec(query2, msg.RoomID, id); err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			slog.Error("error", slog.Any("erorr", err1))
			return err1
		}
		slog.Error("error", slog.Any("erorr", err))
		return err
	}
	return tx.Commit()
}

func (r *WSRepository) GetAllMess(room_id string) ([]entities.Message, error) {
	var msgs []entities.Message
	query := "SELECT text, creator, date FROM messages m INNER JOIN rooms_messages rm ON rm.message_id = m.id INNER JOIN rooms r ON r.id = rm.room_id WHERE r.id = $1"
	if err := r.db.Select(&msgs, query, room_id); err != nil {
		slog.Error("error", err, msgs)
		return []entities.Message{}, err
	}
	result := make([]entities.Message, 0)
	for _, msg := range msgs {
		id, err := r.GetIdUser(msg.Username)
		if err != nil {
			return []entities.Message{}, err
		}
		result = append(result, entities.Message{
			ClientId: id,
			Content:  msg.Content,
			RoomID:   msg.RoomID,
			Username: msg.Username,
			Data:     msg.Data,
		})
	}
	fmt.Println(result)
	return result, nil
}

func (r *WSRepository) GetRooms(group_id string) ([]entities.Rm, error) {
	var rooms []entities.Rm
	query := "SELECT id,name FROM rooms WHERE groupa = $1"
	if err := r.db.Select(&rooms, query, group_id); err != nil {
		slog.Error("error", err)
		return []entities.Rm{}, err
	}
	roomsres := make([]entities.Rm, 0)
	query2 := "SELECT count(*) FROM rooms_users WHERE room_id = $1"
	for _, v := range rooms {
		var a int
		row := r.db.QueryRow(query2, v.Id)
		if err := row.Scan(&a); err != nil {
			return []entities.Rm{}, err
		}
		roomsres = append(roomsres, entities.Rm{
			Id:    v.Id,
			Name:  v.Name,
			Count: fmt.Sprint(a),
		})
	}
	return roomsres, nil
}

func (r *WSRepository) GetRoomId(room string) (string, error) {
	var room_id int
	query := "SELECT id FROM rooms WHERE name = $1"
	row := r.db.QueryRow(query, room)
	if err := row.Scan(&room_id); err != nil {
		return "", err
	}
	return fmt.Sprint(room_id), nil
}

func (r *WSRepository) GetUsername(client_id string) (string, error) {
	var username string
	id, err := strconv.Atoi(client_id)
	if err != nil {
		return "", err
	}
	query := "SELECT name FROM users WHERE id = $1"
	row := r.db.QueryRow(query, id)
	if err = row.Scan(&username); err != nil {
		return "", err
	}
	return username, nil
}

func (r *WSRepository) InsertUser(clientId, roomId string) error {
	query := "INSERT INTO rooms_users(room_id, user_id) VALUES ($1, $2)"
	if _, err := r.db.Exec(query, roomId, clientId); err != nil {
		return err
	}
	return nil
}

func (r *WSRepository) DeleteUser(id, roomid string) error {
	query := "DELETE FROM rooms_users WHERE room_id = $1 AND user_id = $2"
	if _, err := r.db.Exec(query, roomid, id); err != nil {
		return err
	}
	return nil
}

func (r *WSRepository) GetIdUser(username string) (string, error) {
	var id string
	query := "SELECT id FROM users WHERE name = $1"
	row := r.db.QueryRow(query, username)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
