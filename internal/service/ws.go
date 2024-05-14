package service

import (
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/Futturi/Raspisanie/internal/repository"
)

type WService struct {
	repo repository.WSRepo
}

func NewWService(repo repository.WSRepo) *WService {
	return &WService{repo: repo}
}

func (a *WService) CreateRoom(req entities.CreateRoomReq) (string, error) {
	return a.repo.CreateRoom(req)
}

func (a *WService) InsertMessage(msg entities.Message) error {
	return a.repo.InsertMessage(msg)
}

func (a *WService) GetAllMess(room_id string) ([]entities.Message, error) {
	return a.repo.GetAllMess(room_id)
}

func (a *WService) GetRooms(group_id string) ([]entities.Rm, error) {
	return a.repo.GetRooms(group_id)
}

func (a *WService) GetRoomId(room string) (string, error) {
	return a.repo.GetRoomId(room)
}

func (a *WService) GetUsername(client_id string) (string, error) {
	return a.repo.GetUsername(client_id)
}

func (a *WService) InsertUser(clientId, roomId string) error {
	return a.repo.InsertUser(clientId, roomId)
}

func (a *WService) DeleteUser(id, roomid string) error {
	return a.repo.DeleteUser(id, roomid)
}
