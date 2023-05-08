package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type UserModel struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Auth0ID  string             `json:"auth0id,omitempty" bson:"auth0id,"`
	UserName string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Friends  []string           `json:"friends,omitempty" bson:"friends,omitempty"`
}

// FindbyID finds a user based on an Auth0ID and decodes the result into the usermodel
func (u *UserModel) FindByID() error {
	var db = MongoDBClient.Database("chat")
	collection := db.Collection("users")

	err := collection.FindOne(context.TODO(), bson.D{{"auth0id", u.Auth0ID}}).Decode(&u)
	if err != nil {
		// Early return without logging to register user if no user was found
		if err == mongo.ErrNoDocuments {
			return err
		}
		// Handle and log any other error
		return fmt.Errorf("failed to find user: %v", err.Error())
	}
	return nil
}

func (u *UserModel) Register() error {
	var db = MongoDBClient.Database("chat")
	collection := db.Collection("users")
	res, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "auth0id", Value: u.Auth0ID},
		{Key: "username", Value: u.UserName},
		{Key: "email", Value: u.Email},
		{Key: "friends", Value: u.Friends},
	})
	if err != nil {
		log.Println("Failed to insert user to DB: ", err)
		return err
	}
	u.ID = res.InsertedID.(primitive.ObjectID)
	log.Println("New user registered: ", *u)
	return nil

}

func (u *UserModel) UpdateByID() {
	var db = MongoDBClient.Database("chat")
	collection := db.Collection("users")

	res, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", u.ID}}, bson.D{{"$set", u}})
	if err != nil {
		log.Println("Failed updating user: ", err)
		return
	}

	log.Println("Updated user: ", res)

}

func (au Auth0User) CreateUserObject() *UserModel {
	user := new(UserModel)
	user.Auth0ID = au.Sub
	user.Email = au.Email
	user.UserName = au.Nickname
	return user
}
