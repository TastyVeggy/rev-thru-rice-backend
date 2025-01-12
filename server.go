package main

import (
	"log"
	"os"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/middleware"
	"github.com/TastyVeggy/rev-thru-rice-backend/routes"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
)

func main() {
	initialiseApp()

	port := utils.GetEnvWithDefault("PORT", "8080")

	e := echo.New()

	e.Use(middleware.CORS(os.Getenv("FRONTEND_URLS")))

	// Public Routes
	// login/logout/signup
	authRoutes := e.Group("/auth")
	routes.AuthRoutes(authRoutes)
	// viewing subforum is public 
	subforumRoutes := e.Group("/subforums")
	routes.SubforumRoutes(subforumRoutes)
	// viewing countries is public
	countryRoutes := e.Group("/countries")
	routes.CountryRoutes(countryRoutes)
	// viewing posts is public
	postRoutes := e.Group("/posts")
	routes.PostRoutes(postRoutes)
	// viewing comments is public
	commentRoutes := e.Group("/comments")
	routes.CommentRoutes(commentRoutes)
	//viewing shops is public
	shopRoutes := e.Group("/shops")
	routes.ShopRoutes(shopRoutes)


	// Protected Routes
	protected := e.Group("/protected")
	protected.Use(middleware.JWT)
	// Creation of posts in a subforum requires auth
	protectedSubforumRoutes := protected.Group("/subforums")
	routes.ProtectedSubforumRoutes(protectedSubforumRoutes)
	// Editing/ Deleting posts/ Making comments/ Making a shop post requires auth
	protectedPostRoutes := protected.Group("/posts")
	routes.ProtectedPostRoutes(protectedPostRoutes)
	//Editing/Deleting comments requires auth
	protectedCommentsRoutes := protected.Group("/comments")
	routes.ProtectedCommentRoutes(protectedCommentsRoutes)
	// Creating rating/ Editing shops/rating Deleting shops/ratings require auth
	protectedShopRoutes := protected.Group("/shops")
	routes.ProtectedShopRoutes(protectedShopRoutes)
	// Profile settings
	protectedUserRoutes := protected.Group("/users")
	routes.ProtectedUserRoutes(protectedUserRoutes)

	if os.Getenv("GO_ENV") == "production"{
		certFile := os.Getenv("SSL_DIR")+"localhost.pem"
		keyFile := os.Getenv("SSL_DIR")+"localhost-key.pem"
		err := e.StartTLS(":"+port, certFile, keyFile)
		if err != nil {
			log.Fatal("Error starting server: ", err)
		}
	} else {
		e.Logger.Fatal(e.Start(":" + port))
	}

}

func initialiseApp() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")

	wantToCreateTables := utils.GetEnvWithDefault("CREATE_TABLES", "FALSE")
	seedDataDir := os.Getenv("SEED_DATA_DIR")

	db.InitPool(dbURL, wantToCreateTables, seedDataDir)
	utils.InitValidator()
}
