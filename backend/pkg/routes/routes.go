package routes

import (
	"github.com/gorilla/mux"
	"github.com/macihasa/chatapp/backend/pkg/handlers"
	"github.com/macihasa/chatapp/backend/pkg/middleware"
)

func InitMux() *mux.Router {
	r := mux.NewRouter()

	// HTTP: OPTIONS are allowed to handle preflight requests.
	r.HandleFunc("/", middleware.AuthMiddleware(handlers.Landing)).Methods("GET", "OPTIONS")
	r.HandleFunc("/users/login/", middleware.AuthMiddleware(handlers.Login)).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/friends/", middleware.AuthMiddleware(handlers.GetFriends)).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/friends/add/{id}", middleware.AuthMiddleware(handlers.SendFriendRequest)).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/friends/pending/", middleware.AuthMiddleware(handlers.GetPendingFriendRequests)).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/friends/accept/{id}", middleware.AuthMiddleware(handlers.AcceptFriendRequest)).Methods("POST", "OPTIONS")

	// WEBSOCKET:
	r.HandleFunc("/ws", handlers.WsNewClient)
	return r
}
