package livechat

import (
	"golang-chat/internal/models"
	"golang-chat/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan models.Message)

func HandleMessages() {
	for msg := range Broadcast {
		services.SaveMessage(msg)
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing JSON: %v", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
