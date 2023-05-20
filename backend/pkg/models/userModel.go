package models

import (
	"fmt"
	"log"
	"time"
)

type Auth0User struct {
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Nickname      string `json:"nickname"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	UpdatedAt     string `json:"updated_at"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Sub           string `json:"sub"`
}

type User struct {
	ID          int64     `json:"id,omitempty"`
	Auth0ID     string    `json:"auth0id),omitempty"`
	UserName    string    `json:"username"`
	Email       string    `json:"email"`
	DateCreated time.Time `json:"datecreated"`
}

type FriendRequest struct {
	ID          int64     `json:"_id,omitempty"`
	FromID      int64     `json:"fromid"`
	ToID        int64     `json:"toid"`
	Status      string    `json:"status"`
	DateCreated time.Time `json:"datecreated"`
}

type Friends struct {
	FriendshipID int        `json:"friendshipId"`
	UserID       int        `json:"userId"`
	FriendID     int        `json:"friendId"`
	DateAdded    *time.Time `json:"dateAdded,omitempty"`
}


// CreateUserIfNotExist checks if a user exists by it's auth0id and creates a DB record if it doesn't.
// Also, the DB user information is appended into the user object calling the method
func (u *User) CreateUserIfNotExist() error {
	// Insert user to DB if it doesn't exist
	stmt := `INSERT INTO users (auth0id, username, email) SELECT ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM users WHERE auth0id = ?)`

	log.Printf("%+v\n", u)
	_, err := DB.Exec(stmt, u.Auth0ID, u.UserName, u.Email, u.Auth0ID)

	if err != nil {
		log.Println("Failed to insert user to DB: ", err)
		return err
	}

	// Get user info from DB
	stmt = `SELECT * FROM users WHERE auth0id = ?`
	sqlrow := DB.QueryRow(stmt, u.Auth0ID)

	err = sqlrow.Scan(&u.ID, &u.Auth0ID, &u.UserName, &u.Email, &u.DateCreated)
	if err != nil {
		log.Println("Failed to get user info from DB: ", err)
		return err
	}

	log.Println("Successfully got userinfo: ", *u)
	return nil
}

// Get friends returns all friends of the user listed in the friends table
func (u *User) GetFriends() ([]User, error) {
	stmt := `SELECT * FROM users WHERE id IN (SELECT friendid FROM friends WHERE userid = ?)`
	sqlrows, err := DB.Query(stmt, u.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get friends from db: %v", err)
	}

	var users []User
	for sqlrows.Next() {
		var user User
		err = sqlrows.Scan(&user.ID, &user.Auth0ID, &user.UserName, &user.Email, &user.DateCreated)
		if err != nil {
			return nil, fmt.Errorf("failed to get user info from DB: %v", err)
		}
		users = append(users, user)
	}
	
	return users, nil
}

// GetNonFriendUsers finds all users that are not friends with the user
func (u *User) GetNonFriendUsers() ([]User, error) {
	stmt := `SELECT * FROM users WHERE id NOT IN (SELECT toid FROM friendrequests WHERE fromid = ?) AND id != ?`
	sqlrows, err := DB.Query(stmt, u.ID, u.ID)

	if err != nil {
		log.Println("Failed to get non friend users from DB: ", err)
		return nil, err
	}

	var users []User
	for sqlrows.Next() {
		var user User
		err = sqlrows.Scan(&user.ID, &user.Auth0ID, &user.UserName, &user.Email, &user.DateCreated)
		if err != nil {
			log.Println("Failed to get user info from DB: ", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetPendingFriendRequests finds all pending friend requests for the user
func (u *User) GetPendingFriendRequests() ([]FriendRequest, error) {
	stmt := `SELECT * FROM friendrequests WHERE toid = ? AND status = ?`
	sqlrows, err := DB.Query(stmt, u.ID, "pending")

	if err != nil {
		log.Println("Failed to get pending friend requests from DB: ", err)
		return nil, err
	}

	var requests []FriendRequest
	for sqlrows.Next() {
		var request FriendRequest
		err = sqlrows.Scan(&request.ID, &request.FromID, &request.ToID, &request.Status, &request.DateCreated)
		if err != nil {
			log.Println("Failed to get friend request from DB: ", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// SendFriendRequest sends a friend request from the user to another user
func (u *User) SendFriendRequest(toID int64) error {
	stmt := `INSERT INTO friendrequests (fromid, toid, status) SELECT ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM friendrequests WHERE fromid = ? AND toid = ?)`

	_, err := DB.Exec(stmt, u.ID, toID, "pending", u.ID, toID)

	fmt.Printf("Sent friend request from %d to %d\n", u.ID, toID)

	if err != nil {
		log.Println("Failed to insert friend request to DB: ", err)
		return err
	}

	return nil
}

// CreateUserObject creates a user object based on the Auth0User object
func (au Auth0User) CreateUserObject() *User {
	user := new(User)
	user.Auth0ID = au.Sub
	user.Email = au.Email
	user.UserName = au.Nickname
	return user
}
