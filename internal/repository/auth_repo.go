package repository

import (
	"fmt"
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) SignUp(group entities.User) (int, error) {
	var id int
	query := fmt.Sprint("INSERT INTO users(email, password, name, groupa) VALUES($1, $2, $3, $4) RETURNING id")
	row := r.db.QueryRow(query, group.Email, group.Password, group.Name, group.Group)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthRepo) SignIn(user entities.User) (int, string, error) {
	var id int
	var group string
	query := "SELECT id, groupa FROM users WHERE email = $1 AND password = $2"
	row := r.db.QueryRow(query, user.Email, user.Password)
	if err := row.Scan(&id, &group); err != nil {
		return 0, "", err
	}
	return id, group, nil
}
