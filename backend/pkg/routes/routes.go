package routes

import (
	"github.com/gorilla/mux"
	"github.com/macihasa/chatapp/backend/pkg/handlers"
	"github.com/macihasa/chatapp/backend/pkg/middleware"
)

func InitMux() *mux.Router {
	r := mux.NewRouter()

	// HTTP - OPTIONS are allowed to handle preflight requests.
	r.HandleFunc("/", middleware.AuthMiddleware(handlers.Landing)).Methods("GET", "OPTIONS")
	r.HandleFunc("/users/register", handlers.RegisterNewUser).Methods("POST", "OPTIONS")

	// WEBSOCKET
	r.HandleFunc("/ws", handlers.WsNewClient)
	return r
}
