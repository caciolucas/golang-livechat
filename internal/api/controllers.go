package api

import (
	"context"
	"golang-chat/internal/database"
	"golang-chat/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListChannels(c *gin.Context) {
	client := database.ConnectDB()

	coll := client.Database("golang-chat").Collection("channels")

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"errora": err.Error()})
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
