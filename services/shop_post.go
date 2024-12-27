package services

import (
	"context"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
)

type ShopPostReqDTO struct {
	Post   PostReqDTO   `json:"post"`
	Shop   ShopReqDTO   `json:"shop"`
	Rating RatingReqDTO `json:"rating"`
}

type ShopPostResDTO struct {
	Post   PostResDTO   `json:"post"`
	Shop   ShopResDTO   `json:"shop"`
	Rating RatingResDTO `json:"rating"`
}

// Making a shop post involves an atomic transaction of
// adding a post, adding a shop and adding a rating
func AddShopPost(shopPost *ShopPostReqDTO, userID int) (ShopPostResDTO, error) {
	var shopPostRes ShopPostResDTO

	if len(shopPost.Post.Countries) > 0{
		return shopPostRes, fmt.Errorf("should not have any countries in shop post request, country is determined via lat long")
	}

	// reverse geocode location to determine the location of shop for the country associated with the shop post
	location, err := utils.GetShopLocation(shopPost.Shop.Lat, shopPost.Shop.Lng)
	if err != nil {
		return shopPostRes, fmt.Errorf("error in getting location: %v", err)
	}

	shopPost.Post.Countries = append(shopPost.Post.Countries, location.Country)


	// Begin the process of adding post, shop and rating
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return shopPostRes, fmt.Errorf("unable to start shop post transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	shopPostRes.Post, err = AddPostInTx(tx, &shopPost.Post, userID)
	if err != nil {
		return shopPostRes, fmt.Errorf("error adding post: %v", err)
	}

	shopPostRes.Shop, err = AddShopInTx(tx, &shopPost.Shop, location, userID, shopPostRes.Post.ID)
	if err != nil {
		return shopPostRes, fmt.Errorf("error adding shop: %v", err)
	}

	shopPostRes.Rating, err = AddRatingInTx(tx, &shopPost.Rating, shopPostRes.Shop.ID, userID)
	if err != nil {
		return shopPostRes, fmt.Errorf("error adding rating: %v", err)
	}

	err = tx.Commit(context.Background())
	return shopPostRes, err
}