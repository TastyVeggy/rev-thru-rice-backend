package services

import (
	"context"
	"errors"

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
	return addRatingInTx(nil, rating, shopID, userID)
}

func FetchRatingByShopandUser(shopID int, userID int) (RatingResDTO, error) {
	var ratingRes RatingResDTO

	query := `
		SELECT ratings.*, shops.name, users.username 
		FROM ratings
		JOIN shops on ratings.shop_id = shops.id
		JOIN users on ratings.user_id = users.id
		WHERE shop_id = $1 AND user_id = $2
	`

	err := db.Pool.QueryRow(
		context.Background(),
		query,
		shopID,
		userID,
	).Scan(
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

func FetchRatingByPostandUser(postID int, userID int) (RatingResDTO, error) {
	var ratingRes RatingResDTO

	query := `
		SELECT ratings.*, shops.name, users.username 
		FROM ratings
		JOIN shops on ratings.shop_id = shops.id
		JOIN posts on posts.id=shops.post_id
		JOIN users on ratings.user_id = users.id
		WHERE shops.post_id = $1 AND ratings.user_id = $2
	`

	err := db.Pool.QueryRow(
		context.Background(),
		query,
		postID,
		userID,
	).Scan(
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

func UpdateRating(rating *RatingReqDTO, userID int, shopID int) (RatingResDTO, error) {
	var ratingRes RatingResDTO
	query := `
		WITH new_rating AS (
			UPDATE ratings
			SET score = $1
			WHERE user_id=$2 AND shop_id=$3
			RETURNING *
		)
		SELECT new_rating.*, shops.name, users.username
		FROM new_rating
		JOIN shops on new_rating.shop_id = shops.id
		JOIN users on new_rating.user_id = users.id
	`

	err := db.Pool.QueryRow(
		context.Background(),
		query,
		rating.Score,
		userID,
		shopID,
	).Scan(
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

func RemoveRating(shopID int, userID int) error {
	// Important condition: Cannot remove rating for the shop if user is the one that made the shop post (every shop is tied to one and only one post, and comes with at least one rating from the one who make the post)
	query := `
		DELETE FROM ratings
		USING shops, posts
		WHERE shops.id = ratings.shop_id
			AND posts.id = shops.post_id
			AND posts.user_id <> $1
			AND ratings.shop_id = $2
			AND ratings.user_id = $1 
	`
	commandTag, err := db.Pool.Exec(context.Background(), query, userID, shopID)
	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}
	return err
}

func addRatingInTx(tx pgx.Tx, rating *RatingReqDTO, shopID int, userID int) (RatingResDTO, error) {
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
