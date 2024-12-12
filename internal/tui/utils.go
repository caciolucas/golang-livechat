package tui

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang-chat/internal/models"
	"golang-chat/internal/services"
	"golang-chat/internal/tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

type NewMessageReceived struct {
	Message models.Message
}

func resizeTUI(m model, msg tea.WindowSizeMsg) model {
	frameWidth, frameHeight := styles.DocStyle.GetFrameSize()
	terminalWidth, terminalHeight := msg.Width, msg.Height

	listWidth := (terminalWidth - frameWidth) / 3
	listHeight := terminalHeight - frameHeight
	styles.ChannelListStyle = styles.ChannelListStyle.Width(listWidth).Height(listHeight)

	chatHistoryWidth := (2 * (terminalWidth - frameWidth)) / 3
	chatHistoryHeight := terminalHeight - frameHeight - 6
	m.chatHistory.Width = chatHistoryWidth
	m.chatHistory.Height = chatHistoryHeight
	styles.ChatHistoryStyle = styles.ChatHistoryStyle.Width(chatHistoryWidth).Height(chatHistoryHeight)

	m.messageInput.SetWidth((2 * (terminalWidth - frameWidth)) / 3)

	styles.PromptStyle = styles.PromptStyle.Width(24).Margin(
		(terminalHeight-frameHeight)/2,
		((terminalWidth-frameWidth)-24)/2,
	)
	LogToFile(fmt.Sprintf("Terminal width: %d, Terminal height: %d\n", terminalWidth, terminalHeight))
	LogToFile(fmt.Sprintf("Frame width: %d, Frame height: %d\n", frameWidth, frameHeight))

	return m
}

func formatMessages(messages []models.Message) string {
	var builder strings.Builder
	for _, msg := range messages {
		builder.WriteString(fmt.Sprintf("%s: %s\n\n", msg.Username, msg.Message))
	}
	return builder.String()
}

func setViewportContent(channel *models.Channel, vp *viewport.Model) {
	vp.SetContent(formatMessages(channel.Messages))
}

func loadMessages(channel *models.Channel, vp *viewport.Model) {
	if channel == nil {
		return
	}

	messages, err := services.FetchMessages(channel.ID)
	channel.Messages = messages

	if err != nil {
		log.Fatalf("Error fetching messages: %v", err)
	}

	setViewportContent(channel, vp)
}

func listenChannelWSMessages(m *model) {
	for {
		mt, msg, err := m.wsConnection.ReadMessage()
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		if mt == websocket.TextMessage {
			var message models.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Fatalf("Error unmarshalling message: %v", err)
				continue
			}

			m.incomingMessages <- message
		}
	}
}

func readIncomingMessages(incoming chan models.Message) tea.Cmd {
	return func() tea.Msg {
		msg := <-incoming
		return NewMessageReceived{Message: msg}
	}
}

// NOTE: Debug function, remove before production
func LogToFile(data string) {
	fileName := "debug.log"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		fmt.Println("Error writing to file: ", err)
		os.Exit(1)
	}

}
