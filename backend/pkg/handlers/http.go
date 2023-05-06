package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/macihasa/chatapp/backend/pkg/models"
)

// handleServerError checks any errors and logs them both to the client and standard output
func handleServerError(w http.ResponseWriter, msg string, err error) {
	if err != nil {
		log.Println(msg, err)
		http.Error(w, msg+" "+err.Error(), http.StatusInternalServerError)
	}
}

// Landing is a health check of the server.
func Landing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is live!", r.RemoteAddr)
}

// RegisterNewUser creates a new user database with the credentials passed in the request.
func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	var user models.UserModel

	defer func() {
		err := r.Body.Close()
		handleServerError(w, "Couldn't close request body: ", err)
	}()

	bs, err := io.ReadAll(r.Body)
	handleServerError(w, "Failed to read request body: ", err)

	err = json.Unmarshal(bs, &user)
	handleServerError(w, "Failed to unmarshal body: ", err)

	id, err := user.Register()
	handleServerError(w, "Could not register user: ", err)

	// Set user id to objectId created by the database
	user.ID = id

	w.Write([]byte(fmt.Sprintln("New user registered: ", user)))

}
