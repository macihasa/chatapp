package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username" bson:"username"`
	PassWord string             `json:"password" bson:"password"`
	Friends  []string           `json:"friends,omitempty" bson:"friends,omitempty"`
}

func (u UserModel) Register() (primitive.ObjectID, error) {
	var db = MongoDBClient.Database("chat")
	collection := db.Collection("users")

	res, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "username", Value: u.UserName},
		{Key: "password", Value: u.PassWord},
		{Key: "friends", Value: u.Friends},
	})
	if err != nil {
		log.Println("Failed to insert user to DB: ", err)
		return primitive.NilObjectID, err
	}
	log.Println("New user registered: ", res)
	return res.InsertedID.(primitive.ObjectID), nil

}

func (u UserModel) UpdateByID() {
	var db = MongoDBClient.Database("chat")
	collection := db.Collection("users")

	res, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", u.ID}}, bson.D{{"$set", u}})
	if err != nil {
		log.Println("Failed updating user: ", err)
		return
	}

	log.Println("Updated user: ", res)

}

func (u UserModel) GetByID() {
	var db = MongoDBClient.Database("chat")
	collection := db.Collection("users")

	res := collection.FindOne(context.TODO(), u)
	var user UserModel
	err := res.Decode(&user)
	if err != nil {
		log.Println("Couldnt decode res :", err)
	}
	log.Println(user)

}
