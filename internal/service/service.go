package service

import (
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/Futturi/Raspisanie/internal/repository"
)

type Serivce struct {
	Raspisanie
	Auth
	WS
	UserService
}

func NewSerivce(repo *repository.Repository) *Serivce {
	return &Serivce{Raspisanie: NewRaspService(repo.Raspisanie), Auth: NewAuthService(repo.Auth), WS: NewWService(repo), UserService: NewUserService(repo.UserRepo)}
}

type Raspisanie interface {
	GetRasp(group entities.Group, gr string) (entities.Raspisanie, error)
}

type Auth interface {
	SignUp(group entities.User) (int, error)
	SignIn(user entities.User) (string, error)
	ParseToken(header string) (int, string, error)
}

type WS interface {
	CreateRoom(req entities.CreateRoomReq) (string, error)
	InsertMessage(msg entities.Message) error
	GetAllMess(room_id string) ([]entities.Message, error)
	GetRooms(group_id string) ([]entities.Rm, error)
	GetRoomId(room string) (string, error)
	GetUsername(client_id string) (string, error)
	InsertUser(clientId, roomId string) error
	DeleteUser(id, roomid string) error
}

type UserService interface {
	GetUser(id string) (entities.User, error)
	UpdateUser(id string, user entities.UpdateUser) error
}
