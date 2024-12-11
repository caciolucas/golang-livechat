package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Messages []Message          `json:"messages"`
}

func (c Channel) String() string {
	return fmt.Sprintf("\"%s\" (%s)", c.Name, c.ID.Hex())
}
