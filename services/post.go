package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/jackc/pgx/v5"
)

type PostResDTO struct {
	models.Post
	Username string `json:"username"`
	Countries []string `json:"countries"`
}

type PostReqDTO struct {
	SubforumID int    `json:"subforum_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Countries []string `json:"countries"`
}

// for generic posts (not shop subforums)
func AddPost(post *PostReqDTO, userID int) (PostResDTO, error) {
	var postRes PostResDTO
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return postRes, fmt.Errorf("unable to start post transaction: %v", err)
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

	err := db.Pool.QueryRow(
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
	return postRes, err
}

func RemovePost(postID int, userID int) (int64, error) {
	query := `
		DELETE FROM posts
		WHERE id=$1 AND user_id=$2
	`
	commandTag, err := db.Pool.Exec(context.Background(), query, postID, userID)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func FetchPostByID(postID int) (PostResDTO, error) {
	var post PostResDTO

	query := `
		SELECT posts.*, users.username
		FROM posts
		JOIN users ON posts.user_id = users.id
	 	WHERE posts.id = $1
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
	)

	return post, err
}

func FetchPosts(limit int, offset int, subforumID int, userID int) ([]PostResDTO, error) {
	var (
		posts []PostResDTO
		post  PostResDTO
		rows  pgx.Rows
		err   error
	)

	query := `
		SELECT posts.*, users.username
		FROM posts 
		JOIN users ON posts.user_id=users.id
	`
	params := []interface{}{}
	placeholderIndex := 1

	if subforumID != -1 {
		query += fmt.Sprintf(" AND posts.subforum_id = $%d", placeholderIndex)
		params = append(params, subforumID)
		placeholderIndex++
	}
	if userID != -1 {
		query += fmt.Sprintf(" AND posts.user_id=$%d", placeholderIndex)
		params = append(params, userID)
		placeholderIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", placeholderIndex, placeholderIndex+1)
	params = append(params, limit, offset)

	rows, err = db.Pool.Query(context.Background(), query, params...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.SubforumID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CommentCount,
			&post.CreatedAt,
			&post.Username,
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

	row := tx.QueryRow(
		context.Background(),
		post_query,
		post.SubforumID,
		userID,
		post.Title,
		post.Content,
	)

	err := row.Scan(
		&postRes.ID,
		&postRes.SubforumID,
		&postRes.UserID,
		&postRes.Title,
		&postRes.Content,
		&postRes.CommentCount,
		&postRes.CreatedAt,
		&postRes.Username,
	)

	// handling of countries
	for _, country := range post.Countries {
		var countryID int

		// To check if country is valid
		err := tx.QueryRow(context.Background(), "SELECT id FROM countries where name =$1", country).Scan(&countryID)

		if err != nil {
			if err.Error() == "no rows in result set"{
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
			return postRes, fmt.Errorf("error inserting into post_country link table")
		}
		postRes.Countries = append(postRes.Countries, country)
	}

	return postRes, err
}