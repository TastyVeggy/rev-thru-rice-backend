package services

import (
	"context"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/jackc/pgx/v5"
)

type RatingReqDTO struct {
	Score int `json:"score"`
}

type RatingResDTO struct {
	models.Rating
	ShopName string `json:"shop_name"`
	Username string `json:"username"`
}

func AddRating(rating *RatingReqDTO, shopID int, userID int) (RatingResDTO, error) {
	return AddRatingInTx(nil, rating, shopID, userID)
}

func AddRatingInTx(tx pgx.Tx, rating *RatingReqDTO, shopID int, userID int) (RatingResDTO, error) {
	var ratingRes RatingResDTO
	query := `
		WITH new_rating AS (
			INSERT INTO ratings (shop_id, user_id, score)
			VALUES ($1,$2,$3)
			RETURNING *
		)
		SELECT new_rating.*, shops.name, users.username
		FROM new_rating
		JOIN shops on new_rating.shop_id = shops.id
		JOIN users on new_rating.user_id = users.id
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(
			context.Background(),
			query,
			shopID,
			userID,
			rating.Score,
		)
	} else {
		row = db.Pool.QueryRow(
			context.Background(),
			query,
			shopID,
			userID,
			rating.Score,
		)
	}

	err := row.Scan(
		&ratingRes.ID,
		&ratingRes.ShopID,
		&ratingRes.UserID,
		&ratingRes.Score,
		&ratingRes.CreatedAt,
		&ratingRes.ShopName,
		&ratingRes.Username,
	)

	return ratingRes, err

}
