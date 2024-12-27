package services

import (
	"context"
	"errors"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/jackc/pgx/v5"
)

func FetchCountryIDbyName(name string) (int, error) {
	var CountryID int
	query := "SELECT id FROM countries WHERE name = $1"

	err := db.Pool.QueryRow(context.Background(), query, name).Scan(&CountryID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return CountryID, errors.New("country not part of list")
		}
		return CountryID, err
	}
	return CountryID, nil

}
