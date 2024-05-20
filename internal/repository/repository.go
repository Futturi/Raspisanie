package repository

import (
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Raspisanie
	Auth
	WSRepo
	UserRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Raspisanie: NewRaspRepo(db), Auth: NewAuthRepo(db), WSRepo: NewWSRepository(db), UserRepo: NewUser_Repo(db)}
}

type Raspisanie interface {
	GetRasp(group entities.Group, gr string) (entities.Raspisanie, error)
}

type Auth interface {
	SignUp(group entities.User) (int, error)
	SignIn(user entities.User) (int, string, error)
}

type WSRepo interface {
	CreateRoom(req entities.CreateRoomReq) (string, error)
	InsertMessage(msg entities.Message) error
	GetAllMess(room_id string) ([]entities.Message, error)
	GetRooms(group_id string) ([]entities.Rm, error)
	GetRoomId(room string) (string, error)
	GetUsername(client_id string) (string, error)
	InsertUser(clientId, roomId string) error
	DeleteUser(id, roomid string) error
	GetIdUser(username string) (string, error)
}

type UserRepo interface {
	GetUser(id string) (entities.User, error)
	UpdateUser(id string, user entities.UpdateUser) error
}
