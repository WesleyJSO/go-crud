package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Person struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty"`
	Age  int                `json:"age,omitempty"`
}