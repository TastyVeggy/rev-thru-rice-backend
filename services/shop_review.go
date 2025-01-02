package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/models"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
)

type ShopReviewReqDTO struct {
	Post   PostReqDTO   `json:"post"`
	Shop   ShopReqDTO   `json:"shop"`
	Rating RatingReqDTO `json:"rating"`
}

type ShopReviewResDTO struct {
	Post   PostResDTO   `json:"post"`
	Shop   ShopResDTO   `json:"shop"`
	Rating RatingResDTO `json:"rating"`
}

// Making a shop review involves an atomic transaction of
// adding a post, adding a shop and adding a rating
func AddShopReview(shopReview *ShopReviewReqDTO, userID int, subforumID int) (ShopReviewResDTO, error) {
	var shopReviewRes ShopReviewResDTO

	subforumCategory, err := models.FetchSubforumCategorybyID(subforumID)
	if err != nil {
		return shopReviewRes, fmt.Errorf("unable to determine if subforum is shop review, %v", err)
	}
	if subforumCategory != "Review"{
		return shopReviewRes, errors.New("cannot add shop review to non shop review subforums")
	}

	if len(shopReview.Post.Countries) > 0 {
		return shopReviewRes, errors.New("should not have any countries in shop post request, country is determined via lat long")
	}

	// reverse geocode location to determine the location of shop for the country associated with the shop post
	location, err := utils.GetShopLocation(shopReview.Shop.Lat, shopReview.Shop.Lng)
	if err != nil {
		return shopReviewRes, fmt.Errorf("error in getting location: %v", err)
	}

	shopReview.Post.Countries = append(shopReview.Post.Countries, location.Country)

	// Begin the process of adding post, shop and rating
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return shopReviewRes, fmt.Errorf("unable to start shop post transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	shopReviewRes.Post, err = addPostInTx(tx, &shopReview.Post, userID, subforumID)
	if err != nil {
		return shopReviewRes, fmt.Errorf("error adding post: %v", err)
	}

	shopReviewRes.Shop, err = addShopInTx(tx, &shopReview.Shop, location, userID, shopReviewRes.Post.ID)
	if err != nil {
		return shopReviewRes, fmt.Errorf("error adding shop: %v", err)
	}

	shopReviewRes.Rating, err = addRatingInTx(tx, &shopReview.Rating, shopReviewRes.Shop.ID, userID)
	if err != nil {
		return shopReviewRes, fmt.Errorf("error adding rating: %v", err)
	}

	err = tx.Commit(context.Background())
	return shopReviewRes, err
}
