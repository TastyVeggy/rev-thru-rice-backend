package models

import (
	"context"
	"fmt"
	"time"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/jackc/pgx/v5"
)

type Post struct {
	ID           int       `json:"id"`
	SubforumID   int       `json:"subforum_id"`
	UserID       int       `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type PostResDTO struct {
	Post
	Username string `json:"username"`
}

type PostReqDTO struct {
	SubforumID int    `json:"subforum_id"`
	UserID     int    `json:"user_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
}

func AddPost(post *PostReqDTO) error {
	query := `
		INSERT INTO posts (subforum_id, user_id, title, content)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Pool.Exec(
		context.Background(),
		query,
		post.SubforumID,
		post.UserID,
		post.Title,
		post.Content,
	)

	if err != nil {
		return err
	}
	return nil
}

func EditPost(postID int, post *PostReqDTO) (int64, error) {
	query := `
		UPDATE posts
		SET subforum_id=$1, title=$2, content=$3
		WHERE id=$4 AND user_id=$5
	`
	commandTag, err := db.Pool.Exec(context.Background(), query, post.SubforumID, post.Title, post.Content, postID, post.UserID)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
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
	if err != nil {
		return post, err
	}

	return post, nil
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

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", placeholderIndex, placeholderIndex + 1)
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
