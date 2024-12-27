package services

import (
	"context"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShopReqDTO struct {
	Name string  `json:"name"`
	Type string  `json:"type"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}


func AddShop(shop *ShopReqDTO, userID int, postID int) error {
	return AddShopinTx(nil, shop, userID, postID)
}

func AddShopinTx(tx *pgxpool.Tx, shop *ShopReqDTO, userID int, postID int) error {
	var err error

	location, err := utils.GetShopLocation(shop.Lat, shop.Lng) 
	if err != nil {
		return fmt.Errorf("error in getting location: %v", err)
	}

	countryID, err := FetchCountryIDbyName(location.Country)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO shops (post_id, name, country_id, type, lat, lng, address, map_link) 
		SELECT $1,$2,$3,$4,$5,$6,$7,$8
		WHERE EXISTS( 
			SELECT 1
			FROM posts
			WHERE id = $1 AND user_id = $9
		)
	`
	if tx != nil {
		_, err = tx.Exec(
			context.Background(),
			query,
			postID,
			shop.Name,
			countryID,
			shop.Type,
			shop.Lat,
			shop.Lng,
			location.Address,
			location.MapLink,
			userID,
		)
	} else {
		_, err = db.Pool.Exec(
			context.Background(),
			query,
			postID,
			shop.Name,
			countryID,
			shop.Type,
			shop.Lat,
			shop.Lng,
			location.Address,
			location.MapLink,
			userID,
		)
	}

	if err != nil {
		return err
	}
	return nil
}