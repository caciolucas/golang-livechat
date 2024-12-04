package api

import (
	"context"
	"golang-chat/internal/database"
	"golang-chat/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func ListChannels(c *gin.Context) {
	client := database.ConnectDB()

	coll := client.Database("golang-chat").Collection("channels")

	// Find all channels in the collection, cast to a slice of Channels and return
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
