package models

import (
	"context"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
)

type Subforum struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Image       string `json:"image"`
}

func FetchAllSubforums() ([]Subforum, error) {
	var subforums []Subforum

	query := "SELECT * FROM subforums ORDER BY id"

	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subforum Subforum
		err := rows.Scan(&subforum.ID, &subforum.Name, &subforum.Description, &subforum.Category, &subforum.Image)
		if err != nil {
			return nil, err
		}
		subforums = append(subforums, subforum)
	}
	return subforums, nil
}

func FetchSubforumCategorybyID(ID int) (string, error) {
	var category string

	query := "SELECT category FROM subforums WHERE id=$1"

	err := db.Pool.QueryRow(context.Background(), query, ID).Scan(&category)

	return category, err
}
