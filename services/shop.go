package services

import (
	"context"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShopReqDTO struct {
	Name string  `json:"name"`
	Type string  `json:"type"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

type ShopResDTO struct {
	models.Shop
	PostTitle string `json:"post_title"`
	Country   string `json:"country"`
}

func AddShop(shop *ShopReqDTO, userID int, postID int) (any, error) {
	return AddShopinTx(nil, shop, userID, postID)
}

func AddShopinTx(tx *pgxpool.Tx, shop *ShopReqDTO, userID int, postID int) (any, error) {
	var shopRes ShopResDTO

	location, err := utils.GetShopLocation(shop.Lat, shop.Lng)
	if err != nil {
		return -1, fmt.Errorf("error in getting location: %v", err)
	}

	countryID, err := FetchCountryIDbyName(location.Country)
	if err != nil {
		return -1, err
	}

	query := `
		WITH new_shop AS (
			INSERT INTO shops (post_id, name, country_id, type, lat, lng, address, map_link) 
			SELECT $1,$2,$3,$4,$5,$6,$7,$8
			WHERE EXISTS( 
				SELECT 1
				FROM posts
				WHERE id = $1 AND user_id = $9
			)
			RETURNING *
		)
		SELECT new_shop.*, posts.title, countries.name
		FROM new_shop
		JOIN posts ON posts.id = new_shop.post_id
		JOIN countries ON countries.id = new_shop.country_id
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(
			context.Background(),
			query,
			postID,
			shop.Name,
			countryID,
			shop.Type,
			shop.Lat,
			shop.Lng,
			&location.Address,
			location.MapLink,
			userID,
		)
	} else {
		row = db.Pool.QueryRow(
			context.Background(),
			query,
			postID,
			shop.Name,
			countryID,
			shop.Type,
			shop.Lat,
			shop.Lng,
			&location.Address,
			location.MapLink,
			userID,
		)
	}

	err = row.Scan(
		&shopRes.ID,
		&shopRes.PostID,
		&shopRes.Name,
		&shopRes.AvgRating,
		&shopRes.CountryID,
		&shopRes.Type,
		&shopRes.Lat,
		&shopRes.Lng,
		&shopRes.Address,
		&shopRes.MapLink,
		&shopRes.PostTitle,
		&shopRes.Country,
	)
	return shopRes, err
}
