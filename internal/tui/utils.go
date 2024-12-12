package tui

import (
	"encoding/json"
	"fmt"
	"log"
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

type HostEntered struct {
	Host string
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

func loadMessages(channel *models.Channel, m *model) {
	if channel == nil {
		return
	}

	messages, err := services.FetchMessages(channel.ID, m.host)
	channel.Messages = messages

	if err != nil {
		log.Fatalf("Error fetching messages: %v", err)
	}

	setViewportContent(channel, &m.chatHistory)
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

func hostEntered(host string) tea.Cmd {
	return func() tea.Msg {
		return HostEntered{Host: host}
	}
}
func readIncomingMessages(incoming chan models.Message) tea.Cmd {
	return func() tea.Msg {
		msg := <-incoming
		return NewMessageReceived{Message: msg}
	}
}
