package models

import "time"

type Rating struct {
	ID        int       `json:"id"`
	ShopID    int       `json:"shop_id"`
	UserID    int       `json:"user_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
