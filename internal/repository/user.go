package repository

import (
	"fmt"
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/jmoiron/sqlx"
	"strings"
)

type User_Repo struct {
	db *sqlx.DB
}

func NewUser_Repo(db *sqlx.DB) *User_Repo {
	return &User_Repo{db: db}
}

func (r *User_Repo) GetUser(id string) (entities.User, error) {
	var us entities.User
	query := "SELECT email, password, name, groupa FROM users WHERE id = $1"
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&us); err != nil {
		return entities.User{}, err
	}
	return us, nil
}

func (r *User_Repo) UpdateUser(id string, user entities.UpdateUser) error {
	args := make([]any, 0)
	setVal := make([]string, 0)
	argId := 1
	if user.Password != nil {
		setVal = append(setVal, fmt.Sprintf("password=$%d", argId))
		args = append(args, user.Password)
		argId++
	}
	if user.Email != nil {
		setVal = append(setVal, fmt.Sprintf("email=$%d", argId))
		args = append(args, user.Email)
		argId++
	}
	if user.Name != nil {
		setVal = append(setVal, fmt.Sprintf("name=$%d", argId))
		args = append(args, user.Name)
		argId++
	}
	if user.Group != nil {
		setVal = append(setVal, fmt.Sprintf("groupa=$%d", argId))
		args = append(args, user.Group)
		argId++
	}
	strQuery := strings.Join(setVal, ",")
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strQuery, argId)
	args = append(args, id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
