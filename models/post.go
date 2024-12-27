package models

import (
	"time"
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
