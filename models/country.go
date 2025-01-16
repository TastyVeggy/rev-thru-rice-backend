package models

import (
	"context"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
)

type Country struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func FetchAllCountries() ([]Country, error) {
	var countries []Country

	query := "SELECT * FROM countries ORDER BY name"

	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var country Country
		err := rows.Scan(&country.ID, &country.Name, &country.Description)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}
