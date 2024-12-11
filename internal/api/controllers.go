package api

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-chat/internal/database"
	"golang-chat/internal/livechat"
	"golang-chat/internal/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListChannels(c *gin.Context) {
	client := database.ConnectDB()

	coll := client.Database("golang-chat").Collection("channels")

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var channels []models.Channel
	if err = cursor.All(context.Background(), &channels); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, channels)
}

func ListMessages(c *gin.Context) {
	client := database.ConnectDB()

	coll := client.Database("golang-chat").Collection("messages")

	channelId, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	filter := bson.M{"channel": channelId}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		c.JSON(500, gin.H{"errora": err.Error()})
		return
	}

	var messages []models.Message
	if err = cursor.All(context.Background(), &messages); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, messages)
}

func HandleWSConnections(c *gin.Context) {
	conn, err := livechat.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Error upgrading connection:", err.Error())
	}
	defer conn.Close()

	livechat.Clients[conn] = true

	channelId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Invalid channel ID: %v", err)
		return
	}

	fmt.Println("New client connected to channel:", channelId)

	for {
		mt, msg, err := conn.ReadMessage()

		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(livechat.Clients, conn)
			break
		}
		if mt == websocket.TextMessage {
			var message models.Message
			json.Unmarshal(msg, &message)
			message.Channel = channelId
			livechat.Broadcast <- message
		}
	}
	fmt.Println("Client disconnected")
}
