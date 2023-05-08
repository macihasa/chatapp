package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/macihasa/chatapp/backend/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// Landing is a health check of the server.
func Landing(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, "{healthcheck: Server responding!}")
}

// After auth0 authentication is done, login checks wether or not a DB record has been created for the user.
// If yes: return credentials. If no: create DB record and return credentials
func Login(w http.ResponseWriter, r *http.Request) {
	defer closeRequestBody(w, r)

	auth0User := new(models.Auth0User)
	err := json.NewDecoder(r.Body).Decode(auth0User)
	httpServerError(w, "Failed to decode request body:", err)

	user := auth0User.CreateUserObject()

	err = user.FindByID()
	if err == mongo.ErrNoDocuments {
		// Create record in DB if no user was found
		err = user.Register()
		httpServerError(w, "Failed to register user", err)
	}

	writeJSON(w, http.StatusOK, user)

}

// Helper functions --

// writeJSON encodes and sends a json response.
// It also adds neccessary http headers and provided statuscode.
func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// httpServerError checks any errors and logs them both to the client and standard output
func httpServerError(w http.ResponseWriter, msg string, err error) {
	if err != nil {
		log.Println(msg, err)
		http.Error(w, msg+" "+err.Error(), http.StatusInternalServerError)
	}
}

// closeRequestBody closes a request body and handles any httperrors
func closeRequestBody(w http.ResponseWriter, r *http.Request) {
	err := r.Body.Close()
	httpServerError(w, "Failed to close request body:", err)
}
