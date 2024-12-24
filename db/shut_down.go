package db

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func waitForShutdown() {
	// Capture process termination signal
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Block until termination signal is received
	<-signalChannel

	// Close the database pool connection
	if Pool != nil {
		Pool.Close()
	}

	log.Println("Shutting down database pool connection")
}
