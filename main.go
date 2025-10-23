package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"linksupply.io/vmconnect/api"
	"linksupply.io/vmconnect/system"
)

var (
	config *system.Config
)

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment vars")
	}

	config = system.NewConfig()
	_ = system.NewDataSource(&config.Db, config.LogLevel == "INFO")
}

func main() {
	log.Printf("Starting application with config: %s", config)

	// WaitGroup to wait for the goroutines to finish
	wg := &sync.WaitGroup{}

	// Start Application
	StartApplicationBasedOnEntryPoint(wg)

	// Wait for the goroutines to finish
	wg.Wait()

	log.Print("Application stopped.")
}

func StartApplicationBasedOnEntryPoint(wg *sync.WaitGroup) {
	// Setup the API server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Print("Starting the API server.")
		apiServer := api.NewServer(config)
		apiServer.StartServer()
	}()
}
