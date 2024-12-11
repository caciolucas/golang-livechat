package services

import (
	"context"
	"fmt"
	"golang-chat/internal/database"
	"golang-chat/internal/models"
)

func SaveMessage(message models.Message) error {
	client := database.ConnectDB()

	coll := client.Database("golang-chat").Collection("messages")
	_, err := coll.InsertOne(context.Background(), message)
	if err != nil {
		return fmt.Errorf("error saving message: %w", err)
	}
	return nil
}
