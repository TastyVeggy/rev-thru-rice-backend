package models

type Photo struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
	Context string `json:"context"` // comment or post
	ContextID int `json:"context_id"`
}
