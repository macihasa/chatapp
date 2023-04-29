package main

import (
	"log"
	"net/http"

	"github.com/macihasa/chatapp/backend/pkg/handlers"
	"github.com/macihasa/chatapp/backend/pkg/routes"
)

func main() {
	r := routes.InitMux()
	go handlers.StartDistributionHub()
	log.Printf("Server listening on port: 5000...\n")
	err := http.ListenAndServe("localhost:5000", r)
	if err != nil {
		log.Fatal(err)
	}
}
