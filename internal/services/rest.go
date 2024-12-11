package services

import (
	"encoding/json"
	"fmt"
	"golang-chat/internal/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchMessages(channelID primitive.ObjectID) ([]models.Message, error) {
	url := fmt.Sprintf("http://localhost:8080/channels/%s/messages", channelID.Hex())
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %w", err)
	}
	defer response.Body.Close()

	var messages []models.Message
	if err := json.NewDecoder(response.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return messages, nil
}

func FetchChannels() ([]models.Channel, error) {
	response, err := http.Get("http://localhost:8080/channels")
	if err != nil {
		return nil, fmt.Errorf("error fetching channels: %w", err)
	}
	defer response.Body.Close()

	var channels []models.Channel
	if err := json.NewDecoder(response.Body).Decode(&channels); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return channels, nil
}
