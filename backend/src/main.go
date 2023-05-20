package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/macihasa/chatapp/backend/pkg/handlers"
	"github.com/macihasa/chatapp/backend/pkg/models"
	"github.com/macihasa/chatapp/backend/pkg/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load environment variables")
	}

	// Start router, a distributionhub and connect to the DB
	r := routes.InitMux()
	go handlers.StartDistributionHub()

	models.StartMYSQL()
	log.Println("Successfully connected to DB")

	log.Println("Server is running on localhost:5000..")
	err = http.ListenAndServe("localhost:5000", r)
	if err != nil {
		log.Fatal(err)
	}
}
