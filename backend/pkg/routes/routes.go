package routes

import (
	"github.com/gorilla/mux"
	"github.com/macihasa/chatapp/backend/pkg/handlers"
)

func InitMux() *mux.Router {
	r := mux.NewRouter()
	// HTTP
	r.HandleFunc("/", handlers.Landing).Methods("GET")
	r.HandleFunc("/users/register", handlers.RegisterNewUser).Methods("POST")

	// WEBSOCKET
	r.HandleFunc("/ws", handlers.WsNewClient)
	return r
}
