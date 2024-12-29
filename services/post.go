package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/jackc/pgx/v5"
)

type PostResDTO struct {
	models.Post
	Username  string   `json:"username"`
	Countries []string `json:"countries"`
}

type PostReqDTO struct {
	SubforumID int      `json:"subforum_id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Countries  []string `json:"countries"`
}

// for generic posts (not shop subforums)
func AddPost(post *PostReqDTO, userID int) (PostResDTO, error) {
	var postRes PostResDTO
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return postRes, fmt.Errorf("unable to start add post transaction: %v", err)
	}
	defer tx.Rollback(context.Background())
	postRes, err = addPostInTx(tx, post, userID)

	if err != nil {
		return postRes, err
	}

	err = tx.Commit(context.Background())

	return postRes, err
}

func UpdatePost(post *PostReqDTO, postID int, userID int) (PostResDTO, error) {
	var postRes PostResDTO
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return postRes, fmt.Errorf("unable to start edit post transacation, %v", err)
	}
	defer tx.Rollback(context.Background())
	query := `
		WITH new_post AS (
			UPDATE posts
			SET subforum_id=$1, title=$2, content=$3
			WHERE id=$4 AND user_id=$5
			RETURNING *
		)
		SELECT new_post.*, users.username
		FROM new_post
		JOIN users ON new_post.user_id = users.id
	`

	err = tx.QueryRow(
		context.Background(),
		query,
		post.SubforumID,
		post.Title,
		post.Content,
		postID,
		userID,
	).Scan(
		&postRes.ID,
		&postRes.SubforumID,
		&postRes.UserID,
		&postRes.Title,
		&postRes.Content,
		&postRes.CommentCount,
		&postRes.CreatedAt,
		&postRes.Username,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return postRes, errors.New("user unauthorised to change this post or post not found")
		}
		return postRes, fmt.Errorf("error updating entry in posts")
	}

	// delete old country links
	_, err = tx.Exec(
		context.Background(),
		"DELETE FROM post_country WHERE post_id=$1",
		postID,
	)
	if err != nil {
		return postRes, fmt.Errorf("error deleting old post_country links")
	}

	// adding new country links
	for _, country := range post.Countries {
		var countryID int

		err := tx.QueryRow(context.Background(), "SELECT id FROM countries where name =$1", country).Scan(&countryID)

		if err != nil {
			if err.Error() == "no rows in result set" {
				return postRes, fmt.Errorf("country not part of list")
			} else {
				return postRes, err
			}
		}

		_, err = tx.Exec(
			context.Background(),
			"INSERT INTO post_country (post_id, country_id) VALUES ($1, $2)",
			postRes.ID,
			countryID,
		)

		if err != nil {
			return postRes, fmt.Errorf("error inserting post_country links")
		}
		postRes.Countries = append(postRes.Countries, country)
	}

	err = tx.Commit(context.Background())
	return postRes, err
}

func RemovePost(postID int, userID int) error {
	query := `
		DELETE FROM posts
		WHERE id=$1 AND user_id=$2
	`
	commandTag, err := db.Pool.Exec(context.Background(), query, postID, userID)
	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}
	return err
}

func FetchPostByID(postID int) (PostResDTO, error) {
	var post PostResDTO

	query := `
		SELECT 
			posts.*, 
			users.username, 
			array_agg(countries.name) as country_names
		FROM posts
		JOIN users ON posts.user_id = users.id
		JOIN post_country pc ON posts.id = pc.post_id
		JOIN countries ON pc.country_id = countries.id
	 	WHERE posts.id = $1
		GROUP BY posts.id, users.id
	`
	err := db.Pool.QueryRow(context.Background(), query, postID).Scan(
		&post.ID,
		&post.SubforumID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CommentCount,
		&post.CreatedAt,
		&post.Username,
		&post.Countries,
	)

	return post, err
}

/*
Fetch all posts given params and page number and limit
in the case of countries param, as long as the posts is associated with a country in the countries, it will return the post alongside the list of countries the post is associated with

Eg. The sql query built to find 2nd page posts (where each page is 10 posts) in subforum 3 by user with id 4 that is associated with country Brunei and Cambodia (their country ids are 1 and 2):

	SELECT posts.*,
			users.username,
			array_agg(countries.name) AS country_names
	FROM posts
	JOIN users ON posts.user_id=users.id
	JOIN post_country pc ON posts.id = pc.post_id
	JOIN countries ON pc.country_id = countries.id
	WHERE posts.id IN (
		SELECT DISTINCT pc.post_id
		FROM post_country pc
		JOIN countries
		ON pc.country_id = countries.id
		WHERE pc.country_id IN (1,2)
	)
	AND posts.subforum_id=3
	AND posts.user_id=4
	GROUP BY posts.id, users.id
	ORDER BY posts.created_at DESC
	LIMIT 10 OFFSET 10
*/
func FetchPosts(limit int, offset int, subforumID int, userID int, countryIDs []int) ([]PostResDTO, error) {
	var err error
	var posts []PostResDTO

	query := `
		SELECT posts.*, 
			users.username, 
			array_agg(countries.name) AS country_names
		FROM posts 
		JOIN users ON posts.user_id=users.id
		JOIN post_country pc ON posts.id = pc.post_id
		JOIN countries ON pc.country_id = countries.id
	`

	conditions := []string{}
	placeholdersIndex := 1
	params := []any{}
	if len(countryIDs) > 0 {
		placeholders := []string{}
		for i := range countryIDs {
			placeholders = append(placeholders, fmt.Sprintf("$%d", placeholdersIndex))
			params = append(params, countryIDs[i])
			placeholdersIndex++
		}
		placeholderString := strings.Join(placeholders, ",")
		condition := fmt.Sprintf(`
			posts.id IN (
				SELECT DISTINCT pc.post_id
				FROM post_country pc
				JOIN countries 
				ON pc.country_id = countries.id
				WHERE pc.country_id IN (%s)
			)
		`, placeholderString)
		conditions = append(conditions, condition)
	}

	if subforumID > 0 {
		conditions = append(conditions, fmt.Sprintf("%s=$%d", "posts.subforum_id", placeholdersIndex))
		params = append(params, subforumID)
		placeholdersIndex++
	}

	if userID > 0 {
		conditions = append(conditions, fmt.Sprintf("%s=$%d", "posts.user_id", placeholdersIndex))
		params = append(params, userID)
		placeholdersIndex++
	}

	if len(conditions) > 0 {
		query += fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
	}

	query += fmt.Sprintf(`
			GROUP BY posts.id, users.id
			ORDER BY posts.created_at DESC
			LIMIT $%d OFFSET $%d
		`,
		placeholdersIndex,
		placeholdersIndex+1,
	)

	params = append(params, limit, offset)
	fmt.Println(query)
	fmt.Println(params)

	rows, err := db.Pool.Query(context.Background(), query, params...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var post PostResDTO
		err := rows.Scan(
			&post.ID,
			&post.SubforumID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CommentCount,
			&post.CreatedAt,
			&post.Username,
			&post.Countries,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func addPostInTx(tx pgx.Tx, post *PostReqDTO, userID int) (PostResDTO, error) {
	var postRes PostResDTO

	if tx == nil {
		return postRes, errors.New("needs to be part of transaction due to multiple queries required")
	}
	post_query := `
		WITH new_post AS (
			INSERT INTO posts (subforum_id, user_id, title, content)
			VALUES ($1, $2, $3, $4)
			RETURNING *
		)
		SELECT new_post.*, users.username
		FROM new_post
		JOIN users ON  new_post.user_id = users.id
	`

	err := tx.QueryRow(
		context.Background(),
		post_query,
		post.SubforumID,
		userID,
		post.Title,
		post.Content,
	).Scan(
		&postRes.ID,
		&postRes.SubforumID,
		&postRes.UserID,
		&postRes.Title,
		&postRes.Content,
		&postRes.CommentCount,
		&postRes.CreatedAt,
		&postRes.Username,
	)

	if err != nil {
		return postRes, fmt.Errorf("error inserting into posts: %v", err.Error())
	}

	// handling of countries
	for _, country := range post.Countries {
		var countryID int

		// To check if country is valid
		err := tx.QueryRow(context.Background(), "SELECT id FROM countries where name =$1", country).Scan(&countryID)

		if err != nil {
			if err.Error() == "no rows in result set" {
				return postRes, fmt.Errorf("country not part of list")
			} else {
				return postRes, err
			}
		}

		// Insertion into link table betwen post and country
		_, err = tx.Exec(
			context.Background(),
			"INSERT INTO post_country (post_id, country_id) VALUES ($1, $2)",
			postRes.ID,
			countryID,
		)

		if err != nil {
			return postRes, fmt.Errorf("error inserting post_country links")
		}
		postRes.Countries = append(postRes.Countries, country)
	}

	return postRes, err
}
