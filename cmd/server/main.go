package main

import (
	"golang-chat/internal/api"
	"golang-chat/internal/livechat"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/channels", api.ListChannels)
	router.GET("/channels/:id/messages", api.ListMessages)
	router.GET("/channels/:id/ws", api.HandleWSConnections)

	go livechat.HandleMessages()
	router.Run(":8080")

}
