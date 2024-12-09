package main

import (
	// "fmt"
	"golang-chat/internal/api"
	"golang-chat/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("Error upgrading connection:", err.Error())
	}

	defer conn.Close()

	clients[conn] = true

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			delete(clients, conn)
			return
		}

		broadcast <- msg

	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		println("Message received: ", msg.Message)
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing JSON: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	router := gin.Default()

	router.GET("/channels", api.ListChannels)
	router.GET("/channels/:id/messages", api.ListMessages)

	router.Run(":8080")
	// http.HandleFunc("/ws", handleConnections)
	//
	// fmt.Println("Starting server on :8080")
	//
	// go handleMessages()
	//
	// err := http.ListenAndServe(":8080", nil)
	//
	// if err != nil {
	// 	log.Fatal("Error starting server:", err.Error())
	// }
}
