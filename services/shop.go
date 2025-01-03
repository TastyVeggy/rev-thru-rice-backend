package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/jackc/pgx/v5"
)

type ShopReqDTO struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Address *string `json:"address"`
	Country string `json:"country"`
}

type ShopResDTO struct {
	models.Shop
	PostTitle string `json:"post_title"`
	Country   string `json:"country"`
}

func AddShop(shop *ShopReqDTO, userID int, postID int) (ShopResDTO, error) {

	return addShopInTx(nil, shop, userID, postID)
}

func UpdateShop(shop *ShopReqDTO, userID int, shopID int) (ShopResDTO, error) {
	var shopRes ShopResDTO


	countryID, err := FetchCountryIDbyName(shop.Country)
	if err != nil {
		return shopRes, err 
	}

	shopRes.Country = shop.Country

	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return shopRes, fmt.Errorf("unable to start edit shop transaction, %v", err)
	}
	defer tx.Rollback(context.Background())
	query := `
		UPDATE shops
		SET name=$1,country_id=$2,lat=$3,lng=$4,address=$5,map_link=$6
		FROM posts
		WHERE shops.post_id = posts.id
			AND shops.id=$7 
			AND posts.user_id=$8
		RETURNING shops.*,posts.title
	`

	err = tx.QueryRow(
		context.Background(),
		query,
		shop.Name,
		countryID,
		shop.Lat,
		shop.Lng,
		&shop.Address,
		utils.GenerateMapLinkFromLatLng(shop.Lat, shop.Lng),
		shopID,
		userID,
	).Scan(
		&shopRes.ID,
		&shopRes.PostID,
		&shopRes.Name,
		&shopRes.AvgRating,
		&shopRes.CountryID,
		&shopRes.Lat,
		&shopRes.Lng,
		&shopRes.Address,
		&shopRes.MapLink,
		&shopRes.PostTitle,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return shopRes, errors.New("no entry in shop")
		}
		return shopRes, fmt.Errorf("error updating entry in shops:  %v", err)
	}

	country_query := `
		UPDATE post_country
		SET country_id=$1
		WHERE post_id=$2
	`
	_, err = tx.Exec(
		context.Background(),
		country_query,
		shopRes.CountryID,
		shopRes.PostID,
	)
	if err != nil {
		return shopRes, fmt.Errorf("error updating entry in link table: %v", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return shopRes, fmt.Errorf("error commiting edit shop transaction: %v", err)
	}
	return shopRes, err
}

func RemoveShop(shopID int, userID int) error {
	query := `
		DELETE FROM shops
		WHERE id=$1 and user_id=$2
	`

	commandTag, err := db.Pool.Exec(context.Background(), query, shopID, userID)
	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}
	return err
}

func FetchShopByID(shopID int) (ShopResDTO, error) {
	var shop ShopResDTO

	query := `
		SELECT shops.*, posts.title, countries.name
		FROM shops
		JOIN posts ON shops.post_id = posts.id
		JOIN countries ON shops.country_id = countries.id
		WHERE shops.id = $1
	`

	err := db.Pool.QueryRow(
		context.Background(),
		query,
		shopID,
	).Scan(
		&shop.ID,
		&shop.PostID,
		&shop.Name,
		&shop.AvgRating,
		&shop.CountryID,
		&shop.Lat,
		&shop.Lng,
		&shop.Address,
		&shop.MapLink,
		&shop.PostTitle,
		&shop.Country,
	)

	return shop, err
}

func addShopInTx(tx pgx.Tx, shop *ShopReqDTO,  userID int, postID int) (ShopResDTO, error) {

	var shopRes ShopResDTO

	countryID, err := FetchCountryIDbyName(shop.Country)
	if err != nil {
		return shopRes, err
	}

	shopRes.Country = shop.Country

	query := `
		WITH new_shop AS (
			INSERT INTO shops (post_id, name, country_id, lat, lng, address, map_link) 
			SELECT $1,$2,$3,$4,$5,$6,$7
			WHERE EXISTS( 
				SELECT 1
				FROM posts
				WHERE id = $1 AND user_id = $8
			)
			RETURNING *
		)
		SELECT new_shop.*, posts.title
		FROM new_shop
		JOIN posts ON posts.id = new_shop.post_id
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(
			context.Background(),
			query,
			postID,
			shop.Name,
			countryID,
			shop.Lat,
			shop.Lng,
			&shop.Address,
			utils.GenerateMapLinkFromLatLng(shop.Lat, shop.Lng),
			userID,
		)
	} else {
		row = db.Pool.QueryRow(
			context.Background(),
			query,
			postID,
			shop.Name,
			countryID,
			shop.Lat,
			shop.Lng,
			&shop.Address,
			utils.GenerateMapLinkFromLatLng(shop.Lat, shop.Lng),
			userID,
		)
	}

	err = row.Scan(
		&shopRes.ID,
		&shopRes.PostID,
		&shopRes.Name,
		&shopRes.AvgRating,
		&shopRes.CountryID,
		&shopRes.Lat,
		&shopRes.Lng,
		&shopRes.Address,
		&shopRes.MapLink,
		&shopRes.PostTitle,
	)
	return shopRes, err
}
