package services

import (
	"encoding/json"
	"fmt"
	"golang-chat/internal/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchMessages(channelID primitive.ObjectID, host string) ([]models.Message, error) {
	url := fmt.Sprintf("http://%s/channels/%s/messages", host, channelID.Hex())
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

func FetchChannels(host string) ([]models.Channel, error) {
	response, err := http.Get(fmt.Sprintf("http://%s/channels", host))
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
