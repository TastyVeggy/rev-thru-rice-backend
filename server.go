package main

import (
	"net/http"
	"os"

	"github.com/TastyVeggy/rev-thru-rice-backend/db"
	"github.com/TastyVeggy/rev-thru-rice-backend/routes"
	"github.com/TastyVeggy/rev-thru-rice-backend/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	port := utils.GetEnvWithDefault("PORT", "8080")
	wantToCreateTables := utils.GetEnvWithDefault("CREATE_TABLES", "FALSE")
	seedDataDir := os.Getenv("SEED_DATA_DIR")

	db.InitPool(dbURL, wantToCreateTables, seedDataDir)
	utils.InitValidator()

	e := echo.New()
	authRoutes := e.Group("/api/auth")
	routes.AuthRoutes(authRoutes)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + port))

}
