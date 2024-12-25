package models

type Shop struct {
	ID         int    `json:"id"`
	PostID int	`json:"post_id"`
	Name string `json:"name"`
	Location string `json:"location"`
	Address int `json:"address"`
	AvgRating float32 `json:"avg_rating"`
	CountryID int `json:"country_id"`
}