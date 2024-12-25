package db

import (
	"context"
	"log"
	"time"

	"github.com/TastyVeggy/rev-thru-rice-backend/db/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool *pgxpool.Pool
)

func Config(dbURL string) *pgxpool.Config {
	const defaultMaxConns = int32(10)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("You may want to check if .env file is present with DATABASE_URL variable defined. Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}

func initialisePool(dbURL string) error {
	var err error
	Pool, err = pgxpool.NewWithConfig(context.Background(), Config(dbURL))
	if err != nil {
		return err
	}
	log.Println("Database pool initialised")
	return nil
}

func InitPool(dbURL string, wantToCreateTables string, seedDataDir string) {
	err := initialisePool(dbURL)
	if err != nil {
		log.Fatalf("Could not initialise database: %v", err)
	}

	if wantToCreateTables == "TRUE" {
		migrations.InitDB(Pool, true, seedDataDir)
	} else {
		migrations.InitDB(Pool, false, seedDataDir)
	}
	go poolHealthCheckLoop(dbURL)
	go waitForShutdown()
}
