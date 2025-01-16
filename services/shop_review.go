package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/models"
)

type ShopReviewReqDTO struct {
	Post   PostReqDTO   `json:"post"`
	Shop   ShopReqDTO   `json:"shop"`
	Rating RatingReqDTO `json:"rating"`
}

type ShopReviewResDTO struct {
	Post PostResDTO `json:"post"`
	Shop ShopResDTO `json:"shop"`
}

type ShopReviewCreationResDTO struct {
	Post   PostResDTO   `json:"post"`
	Shop   ShopResDTO   `json:"shop"`
	Rating RatingResDTO `json:"rating"`
}

// Making a shop review involves an atomic transaction of
// adding a post, adding a shop and adding a rating
func AddShopReview(shopReview *ShopReviewReqDTO, userID int, subforumID int) (ShopReviewCreationResDTO, error) {
	var shopReviewRes ShopReviewCreationResDTO

	subforumCategory, err := models.FetchSubforumCategorybyID(subforumID)
	if err != nil {
		return shopReviewRes, fmt.Errorf("unable to determine if subforum is shop review, %v", err)
	}
	if subforumCategory != "Review" {
		return shopReviewRes, errors.New("cannot add shop review to non shop review subforums")
	}

	if len(shopReview.Post.Countries) > 0 {
		return shopReviewRes, errors.New("in request, countries should be blank for post, the country attribute should be filled in for the shop value")
	}

	shopReview.Post.Countries = append(shopReview.Post.Countries, shopReview.Shop.Country)

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

	shopReviewRes.Shop, err = addShopInTx(tx, &shopReview.Shop, userID, shopReviewRes.Post.ID)
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

func FetchShopReviewByPostID(postID int) (ShopReviewResDTO, error) {
	var shopReviewRes ShopReviewResDTO

	query := `
		SELECT 
			posts.*, 
			shops.id,
			shops.name,
			shops.avg_rating,
			shops.country_id,
			shops.lat,
			shops.lng,
			shops.address,
			shops.map_link,
			countries.name,
			users.username 
		FROM posts
		JOIN users ON posts.user_id = users.id
		JOIN shops ON shops.post_id = posts.id
		JOIN countries ON countries.id = shops.country_id
		WHERE posts.id = $1
	`

	err := db.Pool.QueryRow(
		context.Background(),
		query,
		postID,
	).Scan(
		&shopReviewRes.Post.ID,
		&shopReviewRes.Post.SubforumID,
		&shopReviewRes.Post.UserID,
		&shopReviewRes.Post.Title,
		&shopReviewRes.Post.Content,
		&shopReviewRes.Post.CommentCount,
		&shopReviewRes.Post.CreatedAt,
		&shopReviewRes.Shop.ID,
		&shopReviewRes.Shop.Name,
		&shopReviewRes.Shop.AvgRating,
		&shopReviewRes.Shop.CountryID,
		&shopReviewRes.Shop.Lat,
		&shopReviewRes.Shop.Lng,
		&shopReviewRes.Shop.Address,
		&shopReviewRes.Shop.MapLink,
		&shopReviewRes.Shop.Country,
		&shopReviewRes.Post.Username,
	)
	if err != nil {
		return shopReviewRes, err
	}

	if shopReviewRes.Post.Username == nil {
		deletedUsername := "[deleted]"
		shopReviewRes.Post.Username = &deletedUsername
	}

	shopReviewRes.Post.Countries = []string{shopReviewRes.Shop.Country}
	shopReviewRes.Shop.PostTitle = shopReviewRes.Post.Title
	shopReviewRes.Shop.PostID = shopReviewRes.Post.ID

	return shopReviewRes, err
}
