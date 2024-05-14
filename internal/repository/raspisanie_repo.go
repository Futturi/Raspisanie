package repository

import (
	"encoding/json"
	"fmt"
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/jmoiron/sqlx"
)

type RaspRepo struct {
	db *sqlx.DB
}

func NewRaspRepo(db *sqlx.DB) *RaspRepo {
	return &RaspRepo{db: db}
}

func (r *RaspRepo) GetRasp(group entities.Group, gr string) (entities.Raspisanie, error) {
	str := ""
	var pairs []entities.Raspisanie
	fmt.Println(group)
	query := fmt.Sprint("SELECT pairs FROM pairs WHERE groups = $1 AND week = $2 and day = $3")
	row := r.db.QueryRow(query, gr, group.Week, group.Day)
	if err := row.Scan(&str); err != nil {
		return entities.Raspisanie{}, err
	}
	err := json.Unmarshal([]byte(str), &pairs)
	if err != nil {
		return entities.Raspisanie{}, err
	}
	return pairs[group.Day], nil
}
