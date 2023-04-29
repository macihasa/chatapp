package routes

import (
	"github.com/gorilla/mux"
	"github.com/macihasa/chatapp/backend/pkg/handlers"
)

func InitMux() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Landing)
	r.HandleFunc("/ws", handlers.WsNewClient)
	return r
}
