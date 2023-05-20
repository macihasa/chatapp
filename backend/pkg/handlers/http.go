package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/macihasa/chatapp/backend/pkg/models"
)

// Landing is a health check of the server.
func Landing(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, "Server is running!")
}

// Login handles the login of a user and creates a new user record if it does not exist.
func Login(w http.ResponseWriter, r *http.Request) {

	auth0user := new(models.Auth0User)

	defer closeRequestBody(w, r)

	err := json.NewDecoder(r.Body).Decode(auth0user)
	handleHttpServerError(w, "Failed to read request body:", err)

	user := auth0user.CreateUserObject()

	err = user.CreateUserIfNotExist()
	handleHttpServerError(w, "Could not register user:", err)

	writeJSON(w, http.StatusOK, user)

}

// GetFriends returns a list of friends of the user.
func GetFriends(w http.ResponseWriter, r *http.Request) {

	user := new(models.User)

	defer closeRequestBody(w, r)

	err := json.NewDecoder(r.Body).Decode(user)
	handleHttpServerError(w, "Failed to read request body:", err)

	friends, err := user.GetFriends()
	handleHttpServerError(w, "Failed to get friends:", err)

	writeJSON(w, http.StatusOK, friends)

}

func GetNonFriendUsers(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)

	defer closeRequestBody(w, r)

	err := json.NewDecoder(r.Body).Decode(user)
	handleHttpServerError(w, "Failed to read request body:", err)

	users, err := user.GetNonFriendUsers()
	handleHttpServerError(w, "Failed to get users:", err)

	writeJSON(w, http.StatusOK, users)
}

// SendFriendRequest sends a friend request to a user.
func SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	handleHttpServerError(w, "Failed to convert id to int:", err)
	log.Println("ID:", id)

	user := new(models.User)

	defer closeRequestBody(w, r)

	err = json.NewDecoder(r.Body).Decode(user)
	handleHttpServerError(w, "Failed to read request body:", err)
	fmt.Println(user)

	err = user.SendFriendRequest(int64(id))
	handleHttpServerError(w, "Failed to add friend:", err)

	writeJSON(w, http.StatusOK, user)

}

func GetPendingFriendRequests(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)

	defer closeRequestBody(w, r)

	err := json.NewDecoder(r.Body).Decode(user)
	handleHttpServerError(w, "Failed to read request body:", err)

	requests, err := user.GetPendingFriendRequests()
	handleHttpServerError(w, "Failed to get requests:", err)

	writeJSON(w, http.StatusOK, requests)
}

// ------------------ Helper functions ------------------

// writeJSON encodes and sends a json response.
// It also adds neccessary http headers and provided statuscode.
func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// handleHttpServerError checks any errors and logs them both to the client and standard output
func handleHttpServerError(w http.ResponseWriter, msg string, err error) {
	if err != nil {
		log.Println(msg, err)
		http.Error(w, msg+" "+err.Error(), http.StatusInternalServerError)
	}
}

// closeRequestBody closes a request body and handles any httperrors
func closeRequestBody(w http.ResponseWriter, r *http.Request) {
	err := r.Body.Close()
	handleHttpServerError(w, "Failed to close request body:", err)
}
