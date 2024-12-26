package models

import (
	"context"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
)

type Subforum struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func FetchSubforums() ([]Subforum, error) {
	var subforums []Subforum

	query := "SELECT * FROM subforums"

	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subforum Subforum
		err := rows.Scan(&subforum.ID, &subforum.Name, &subforum.Description)
		if err != nil {
			return nil, err
		}
		subforums = append(subforums, subforum)
	}
	return subforums, nil
}

// func GetSubForumID(forumName string) (int, error) {
// 	var result int

// 	query := "SELECT id FROM subforums WHERE name = $1"
// 	err := db.Pool.QueryRow(context.Background(), query, forumName).Scan(&result)

// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return result, errors.New("subforum does not exist")
// 		}
// 		return result, err
// 	}
// 	return result, nil
// }
