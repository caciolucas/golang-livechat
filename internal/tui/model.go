package tui

import (
	"fmt"
	"golang-chat/internal/models"
	"golang-chat/internal/services"
	channellist "golang-chat/internal/tui/components/channelList"
	"golang-chat/internal/tui/styles"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type model struct {
	channelsList   channellist.Model
	chatHistory    viewport.Model
	messageInput   textarea.Model
	currentChannel *models.Channel
	prompt         textinput.Model

	username         string
	wsConnection     *websocket.Conn
	incomingMessages chan models.Message
}

func (m model) Init() tea.Cmd {
	return nil
}

func InitialModel() model {
	channels, err := services.FetchChannels()
	if err != nil {
		log.Fatalf("Error fetching channels: %v", err)
	}

	channelsList := channellist.Model{Channels: channels}
	channelsList.Focus()
	channelsList.SetCursor(0)

	messageInput := textarea.New()
	messageInput.Placeholder = "Send a message..."
	messageInput.Prompt = "┃ "
	messageInput.CharLimit = 280
	messageInput.SetWidth(0)
	messageInput.SetHeight(4)
	messageInput.ShowLineNumbers = false
	messageInput.MaxHeight = 4

	chatHistory := viewport.New(0, 0)

	prompt := textinput.New()
	m := model{
		channelsList:     channelsList,
		chatHistory:      chatHistory,
		messageInput:     messageInput,
		incomingMessages: make(chan models.Message),
		prompt:           prompt,
		username:         "",
	}

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.username == "" {
		m.prompt.Focus()
		m.prompt.Placeholder = "Enter your username"
		m.messageInput.Blur()
		m.channelsList.Blur()
	}
	LogToFile("===== UPDATE =====\n")
	LogToFile(fmt.Sprintf("MSG: %v\n", msg))
	LogToFile(fmt.Sprintf("TYPE: %T\n", msg))
	LogToFile("==================\n")

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = resizeTUI(m, msg)
		return m, nil

	case NewMessageReceived:
		m.channelsList.SelectedChannel.Messages = append(m.channelsList.SelectedChannel.Messages, msg.Message)
		setViewportContent(m.channelsList.SelectedChannel, &m.chatHistory)

		// Move to the bottom of the chat history
		m.chatHistory.ViewDown()

		return m, readIncomingMessages(m.incomingMessages)

	case tea.KeyMsg:
		if m.messageInput.Focused() {
			if msg.String() == "esc" {
				m.messageInput.Blur()
				m.channelsList.Focus()
			} else if msg.String() == "enter" {
				message := models.Message{
					Username: m.username,
					Message:  m.messageInput.Value(),
					Channel:  m.channelsList.SelectedChannel.ID,
				}
				m.wsConnection.WriteJSON(message)
				m.messageInput.Reset()
			} else if msg.String() == "up" || msg.String() == "down" {
				if msg.String() == "up" {
					m.chatHistory.ViewUp()
				}
				if msg.String() == "down" {
					m.chatHistory.ViewDown()
				}
			} else {
				c, cmd := m.messageInput.Update(msg)
				m.messageInput = c
				return m, cmd
			}

			return m, nil
		}
		if m.channelsList.Focused() {
			if msg.String() == "enter" {
				m.messageInput.Focus()
				m.messageInput.Reset()
				m.channelsList.Blur()

				m.channelsList.SelectedChannel = &m.channelsList.Channels[m.channelsList.Cursor()]
				loadMessages(m.channelsList.SelectedChannel, &m.chatHistory)

				conn, err := services.ConnectChannelWS(m.channelsList.SelectedChannel.ID)
				if err != nil {
					log.Fatalf("Error connecting to channel WS: %v", err)
				}
				m.wsConnection = conn

				go listenChannelWSMessages(&m)
				return m, readIncomingMessages(m.incomingMessages)

			} else {
				c, cmd := m.channelsList.Update(msg)
				m.channelsList = c
				return m, cmd
			}
		}
		if m.prompt.Focused() {
			if msg.String() == "esc" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			} else if msg.String() == "up" || msg.String() == "down" || msg.String() == "left" || msg.String() == "right" {

				LogToFile(fmt.Sprintf("Prompt style margins %v %v %v %v\n", styles.PromptStyle.GetMarginTop(), styles.PromptStyle.GetMarginRight(), styles.PromptStyle.GetMarginBottom(), styles.PromptStyle.GetMarginLeft()))
				if msg.String() == "up" {
					styles.PromptStyle = styles.PromptStyle.Margin(styles.PromptStyle.GetMarginTop()+1, styles.PromptStyle.GetMarginRight())
				}
				if msg.String() == "down" {
					styles.PromptStyle = styles.PromptStyle.Margin(styles.PromptStyle.GetMarginTop()-1, styles.PromptStyle.GetMarginRight())
				}
				if msg.String() == "left" {
					styles.PromptStyle = styles.PromptStyle.Margin(styles.PromptStyle.GetMarginTop(), styles.PromptStyle.GetMarginRight()+1)
				}
				if msg.String() == "right" {
					styles.PromptStyle = styles.PromptStyle.Margin(styles.PromptStyle.GetMarginTop(), styles.PromptStyle.GetMarginRight()-1)
				}

			} else if msg.String() == "enter" {
				if m.username == "" {
					m.username = m.prompt.Value()
					m.prompt.Reset()
					m.channelsList.Focus()
				}
				return m, nil
			} else {
				c, cmd := m.prompt.Update(msg)
				m.prompt = c
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.username == "" {
		return styles.PromptStyle.Render(m.prompt.View())
	}

	channelsList := styles.ChannelListStyle.Render(m.channelsList.View())
	chatHistory := styles.ChatHistoryStyle.Render(m.chatHistory.View())
	messageInput := styles.MessageInputStyle.Render(m.messageInput.View())

	chatHistoryAndInput := lipgloss.JoinVertical(lipgloss.Top, chatHistory, messageInput)
	return lipgloss.JoinHorizontal(lipgloss.Top, channelsList, chatHistoryAndInput)
}
