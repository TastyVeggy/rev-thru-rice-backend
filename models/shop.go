package models

type Shop struct {
	ID        int      `json:"id"`
	PostID    int      `json:"post_id"`
	Name      string   `json:"name"`
	AvgRating *float32 `json:"avg_rating"`
	CountryID int      `json:"country_id"`
	Lat       float64  `json:"lat"`
	Lng       float64  `json:"lng"`
	Address   *string  `json:"address"`
	MapLink   string   `json:"map_link"`
}
