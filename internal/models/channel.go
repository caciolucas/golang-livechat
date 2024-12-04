package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name"`
}

// Implement list.Item interface for models.Channel
func (c Channel) Title() string       { return c.Name }
func (c Channel) Description() string { return fmt.Sprintf("ID: %s", c.ID.Hex()) }
func (c Channel) FilterValue() string { return c.Name }
