package tui

import (
	"encoding/json"
	"fmt"
	"golang-chat/internal/models"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func formatMessages(messages []models.Message) string {
	var builder strings.Builder
	for _, msg := range messages {
		builder.WriteString(fmt.Sprintf("%s: %s\n\n", msg.Username, msg.Message))
	}
	return builder.String()
}

func fetchMessages(channelID primitive.ObjectID) []models.Message {
	response, err := http.Get(fmt.Sprintf("http://localhost:8080/channels/%s/messages", channelID.Hex()))

	if err != nil {
		fmt.Println("Error fetching messages: ", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	var messages []models.Message

	err = json.NewDecoder(response.Body).Decode(&messages)

	if err != nil {
		fmt.Println("Error decoding JSON: ", err)
		os.Exit(1)
	}

	return messages
}

func fetchChannels() []models.Channel {
	response, err := http.Get("http://localhost:8080/channels")
	if err != nil {
		fmt.Println("Error fetching channels: ", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	var channels []models.Channel

	err = json.NewDecoder(response.Body).Decode(&channels)

	if err != nil {
		fmt.Println("Error decoding JSON: ", err)
		os.Exit(1)
	}

	return channels
}

func resizeTUI(m model, msg tea.WindowSizeMsg) model {
	frameWidth, frameHeight := docStyle.GetFrameSize()
	terminalWidth, terminalHeight := msg.Width, msg.Height

	listWidth := (terminalWidth - frameWidth) / 3
	listHeight := terminalHeight - frameHeight
	m.channelsList.SetSize(listWidth, listHeight)
	channelListStyle = channelListStyle.Width(listWidth).Height(listHeight)

	chatHistoryWidth := (2 * (terminalWidth - frameWidth)) / 3
	chatHistoryHeight := terminalHeight - frameHeight - 6
	m.chatHistory.Width = chatHistoryWidth
	m.chatHistory.Height = chatHistoryHeight
	chatHistoryStyle.Width(chatHistoryWidth).Height(chatHistoryHeight)

	m.messageInput.SetWidth((2 * (terminalWidth - frameWidth)) / 3)

	// NOTE: For debugging porpuses, set the channel title as the dimensions of each element and used vars
	m.channelsList.Title = fmt.Sprintf("Term: %dx%d, List: %dx%d, Chat: %dx%d, Frame: _x%d", terminalWidth, terminalHeight, listWidth, listHeight, chatHistoryWidth, chatHistoryHeight, frameHeight)

	return m
}
