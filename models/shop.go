package models

type Shop struct {
	ID        int     `json:"id"`
	PostID    int     `json:"post_id"`
	Name      string  `json:"name"`
	AvgRating float32 `json:"avg_rating"`
	CountryID int     `json:"country_id"`
	Type    int     `json:"type"`
	Lat float64 `json:"lat"`
	Long float64 `json:"long"`
	Address string `json:"address"`
}
