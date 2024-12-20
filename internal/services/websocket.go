package services

import (
	"encoding/json"
	"fmt"
	"golang-chat/internal/models"

	"github.com/gorilla/websocket"
)

func ConnectChannelWS(host string) (*websocket.Conn, error) {
	url := fmt.Sprintf("ws://%s/channels/ws", host)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to WebSocket server: %w", err)
	}

	return conn, nil
}

func SendMessage(conn *websocket.Conn, message models.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}

	return nil
}
