package service

import (
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/Futturi/Raspisanie/internal/repository"
)

type User_Service struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *User_Service {
	return &User_Service{repo: repo}
}

func (a *User_Service) GetUser(id string) (entities.User, error) {
	us, err := a.repo.GetUser(id)
	if err != nil {
		return entities.User{}, err
	}
	us.Password = ""
	return us, err
}

func (a *User_Service) UpdateUser(id string, user entities.UpdateUser) error {
	if user.Password != nil {
		pas := hashPass(*user.Password)
		user.Password = &pas
	}
	return a.repo.UpdateUser(id, user)
}
