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
	// viewing posts/ shops is public
	postRoutes := e.Group("/post")
	routes.PostRoutes(postRoutes)
	//
	// countryRoutes := e.Group("/countries")
	// routes.CountryRoutes(countryRoutes)

	// Protected Routes
	protected := e.Group("/protected")
	protected.Use(middleware.JWT)

	// user profile
	// userRoutes := protected.Group("/user")
	// routes.UserRoutes(userRoutes)

	// Making/Deleting posts/comments requires auth
	protectedPostRoutes := protected.Group("/post")
	routes.ProtectedPostRoutes(protectedPostRoutes)

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
