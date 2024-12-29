package models

// TODO eventually support having photos in post and comment
type Photo struct {
	ID        int    `json:"id"`
	Url       string `json:"url"`
	Context   string `json:"context"` // comment or post
	ContextID int    `json:"context_id"`
}
