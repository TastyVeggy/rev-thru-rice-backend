package services

import (
	"context"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/jackc/pgx/v5"
)

type CommentReqDTO struct {
	Content string `json:"content"`
}

type CommentResDTO struct {
	models.Comment
	PostTitle string `json:"post_title"`
	Username  string `json:"username"`
}

func AddComment(comment *CommentReqDTO, userID int, postID int) (CommentResDTO, error) {
	var commentRes CommentResDTO
	query := `
		WITH new_comment AS (
			INSERT INTO comments (post_id, user_id, content)
			VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT new_comment.*, posts.title, users.username
		FROM new_comment
		JOIN users ON users.id = new_comment.user_id
		JOIN posts ON posts.id = new_comment.post_id
	`

	err := db.Pool.QueryRow(
		context.Background(),
		query,
		postID,
		userID,
		comment.Content,
	).Scan(
		&commentRes.ID,
		&commentRes.PostID,
		&commentRes.UserID,
		&commentRes.Content,
		&commentRes.CreatedAt,
		&commentRes.PostTitle,
		&commentRes.Username,
	)
	return commentRes, err
}

func EditComment(comment *CommentReqDTO, userID int, commentID int) (CommentResDTO, error) {
	var commentRes CommentResDTO
	query := `
		WITH new_comment AS (
			UPDATE comments
			SET content=$1
			WHERE id=$2 AND user_id=$3
			RETURNING *
		)
		SELECT new_comment.*, posts.title, users.username
		FROM new_comment
		JOIN  
		users ON users.id = new_comment.user_id
		JOIN 
		posts ON posts.id = new_comment.post_id
	`
	err := db.Pool.QueryRow(
		context.Background(),
		query,
		comment.Content,
		commentID,
		userID,
	).Scan(
		&commentRes.ID,
		&commentRes.PostID,
		&commentRes.UserID,
		&commentRes.Content,
		&commentRes.CreatedAt,
		&commentRes.PostTitle,
		&commentRes.Username,
	)
	return commentRes, err
}

func RemoveComment(commentID int, userID int) (int64, error) {
	query := `
		DELETE FROM comments
		WHERE id=$1 AND user_id=$2
	`
	commandTag, err := db.Pool.Exec(context.Background(), query, commentID, userID)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func FetchComments(limit int, offset int, postID int, userID int) ([]CommentResDTO, error) {
	var (
		comments []CommentResDTO
		comment  CommentResDTO
		rows     pgx.Rows
		err      error
	)

	query := `
		SELECT comments.*, posts.title, users.username
		FROM comments
		JOIN
		posts on comments.post_id = posts.id
		JOIN
		users on comments.user_id = users.id
	`
	params := []interface{}{}
	placeholderIndex := 1

	if postID != -1 {
		query += fmt.Sprintf(" AND comments.post_id = $%d", placeholderIndex)
		params = append(params, postID)
		placeholderIndex++
	}
	if userID != -1 {
		query += fmt.Sprintf(" AND comments.user_id=$%d", placeholderIndex)
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
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.PostTitle,
			&comment.Username,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
