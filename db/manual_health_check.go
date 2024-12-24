package db

import (
	"context"
	"fmt"
	"log"
	"time"
)

func PoolHealthCheckLoop(dbURL string) {
	for {
		healthy := isHealthy()
		if !healthy {
			log.Println("Database pool unhealthy Attempting to reconnect...")
			err := reinitialisePool(dbURL)
			if err != nil {
				log.Printf("Failed to reinitialise the database pool: %v. Will try again in 30 seconds.", err)
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func reinitialisePool(dbURL string) error {
	fmt.Println("Reinitialising database pool...")
	if Pool != nil {
		Pool.Close()
	}
	time.Sleep(3 * time.Second)

	return initialisePool(dbURL)
}

func isHealthy() bool {
	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Failed to acquire connection: %v", err)
		return false
	}
	conn.Release()
	return true
}
