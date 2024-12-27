package main

import (
	"net/http"
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

	// Public Routes
	// login/logout/signup
	authRoutes := e.Group("/auth")
	routes.AuthRoutes(authRoutes)
	// viewing posts is public
	postRoutes := e.Group("/posts")
	routes.PostRoutes(postRoutes)
	// viewing comments is public
	commentRoutes := e.Group("/comments")
	routes.CommentRoutes(commentRoutes)

	// countryRoutes := e.Group("/countries")
	// routes.CountryRoutes(countryRoutes)

	// Protected Routes
	protected := e.Group("/protected")
	protected.Use(middleware.JWT)


	// Making/Deleting posts/comments requires auth
	protectedPostRoutes := protected.Group("/posts")
	routes.ProtectedPostRoutes(protectedPostRoutes)

	protectedCommentsRoutes := protected.Group("/comments")
	routes.ProtectedCommentRoutes(protectedCommentsRoutes)

	// Profile settings
	// protectedUserRoutes := protected.Group("/users")
	// routes.ProtectedUserRoutes(protectedUserRoutes)

	// Making/Deleting shops/ratings require auth
	// protectedShopRoutes := protected.Group("/shop")
	// routes.ProtectedShopRoutes(protectedShopRoutes)


	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + port))

}

func initialiseApp() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")

	wantToCreateTables := utils.GetEnvWithDefault("CREATE_TABLES", "FALSE")
	seedDataDir := os.Getenv("SEED_DATA_DIR")

	db.InitPool(dbURL, wantToCreateTables, seedDataDir)
	utils.InitValidator()
}
