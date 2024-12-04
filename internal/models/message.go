package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Message  string             `json:"message" bson:"message"`
	Channel  primitive.ObjectID `json:"channel" bson:"channel"`
}
