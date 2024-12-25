package models

type CommentPhoto struct {
	CommentID int `json:"comment_id"`
	PhotoID   int `json:"photo_id"`
}

type PostPhoto struct {
	PostID  int `json:"post_id"`
	PhotoID int `json:"photo_id"`
}

type PostCountry struct {
	PostID    int `json:"post_id"`
	CountryID int `json:"country_id"`
}