package services

import (
	"context"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
)

type PostCountResDTO struct {
	SubforumID int `json:"subforum_id"`
	PostCount int `json:"post_count"`
}

func FetchSubforumPostCountsbyCountryID(countryID int) ([]PostCountResDTO, error){
	var postCountsRes []PostCountResDTO
	var params []any

	query := `
		SELECT subforums.id, COUNT(posts.id) AS post_count
		FROM subforums
		JOIN posts on subforums.id = posts.subforum_id
	`
	if countryID > 0 {
		query += `
			JOIN post_country pc ON posts.id = pc.post_id
			WHERE pc.country_id = $1
		`
		params = append(params, countryID)
	}

	query += `
		GROUP BY subforums.id
		ORDER BY subforums.id
	`

	rows, err := db.Pool.Query(context.Background(), query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postCount PostCountResDTO
		err := rows.Scan(&postCount.SubforumID, &postCount.PostCount)
		if err != nil {
			return nil, err
		}
		postCountsRes = append(postCountsRes, postCount)
	}
	return postCountsRes, nil
}

// func FetchAllSubforums(countryID *int) ([]SubforumResDTO, error) {
// 	var subforumsRes []SubforumResDTO
// 	var params []any
	
// 	// I'm sure there is a way to do it with one sql query using left joins but the existence of posts without any country makes it tricky but this works for now
// 	subforums, err := models.FetchAllSubforums()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, subforum := range subforums {
// 		subforumsRes = append(subforumsRes, SubforumResDTO{
// 			Subforum:  subforum,
// 			PostCount: 0,
// 			CountryID: countryID,
// 		})

// 	}

// 	query := `
// 		SELECT subforums.id, COUNT(posts.id) AS post_count
// 		FROM subforums
// 		JOIN posts on subforums.id = posts.subforum_id
// 	`

// 	if *countryID > 0{
// 		query += `
// 			JOIN post_country pc ON posts.id = pc.post_id
// 			WHERE pc.country_id = $1
// 		`
// 		params = append(params, countryID)
// 	}

// 	query += `
// 		GROUP BY subforums.id
// 		ORDER BY subforums.id
// 	`

// 	rows, err := db.Pool.Query(context.Background(), query, params...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var PostCount int
// 		var subforumID int
// 		err := rows.Scan(&subforumID, &PostCount)
// 		subforumsRes[subforumID].PostCount = PostCount
// 		if err != nil {
// 			return nil, err
// 		}
// 		subforumsRes[subforumID].PostCount = PostCount
// 	}
// 	return subforumsRes, nil
// }

