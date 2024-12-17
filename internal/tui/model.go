package tui

import (
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
	userPrompt     textinput.Model
	hostPrompt     textinput.Model

	username         string
	host             string
	wsConnection     *websocket.Conn
	incomingMessages chan models.Message
}

func (m model) Init() tea.Cmd {
	return nil
}

func InitialModel() model {
	channelsList := channellist.Model{}
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

	userPrompt := textinput.New()
	userPrompt.Placeholder = "Enter your username"
	hostPrompt := textinput.New()
	hostPrompt.Placeholder = "Enter the host"

	m := model{
		channelsList:     channelsList,
		chatHistory:      chatHistory,
		messageInput:     messageInput,
		incomingMessages: make(chan models.Message),
		userPrompt:       userPrompt,
		hostPrompt:       hostPrompt,
		username:         "",
		host:             "",
	}

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.username == "" {
		m.userPrompt.Focus()
		m.messageInput.Blur()
		m.channelsList.Blur()
	}
	if m.host == "" {
		m.hostPrompt.Focus()
		m.messageInput.Blur()
		m.channelsList.Blur()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = resizeTUI(m, msg)
		return m, nil

	case HostEntered:
		channels, err := services.FetchChannels(m.host)
		if err != nil {
			log.Fatalf("Error fetching channels: %v", err)
		}
		m.channelsList.Channels = channels

		conn, err := services.ConnectChannelWS(m.host)
		if err != nil {
			log.Fatalf("Error connecting to channel WS: %v", err)
		}
		m.wsConnection = conn

		go listenChannelWSMessages(&m)
		return m, nil

	case NewMessageReceived:

		if m.channelsList.SelectedChannel.ID == msg.Message.Channel {
			m.channelsList.SelectedChannel.Messages = append(m.channelsList.SelectedChannel.Messages, msg.Message)
			setViewportContent(m.channelsList.SelectedChannel, &m.chatHistory)
			m.chatHistory.ViewDown()
		}

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
				loadMessages(m.channelsList.SelectedChannel, &m)

				return m, readIncomingMessages(m.incomingMessages)

			} else {
				c, cmd := m.channelsList.Update(msg)
				m.channelsList = c
				return m, cmd
			}
		}
		if m.hostPrompt.Focused() {
			if msg.String() == "esc" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			} else if msg.String() == "enter" {
				if m.host == "" {
					m.host = m.hostPrompt.Value()
					m.hostPrompt.Blur()
					m.userPrompt.Focus()

					return m, hostEntered(m.host)
				}
				return m, nil
			} else {
				c, cmd := m.hostPrompt.Update(msg)
				m.hostPrompt = c
				return m, cmd
			}
		}
		if m.userPrompt.Focused() {
			if msg.String() == "esc" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			} else if msg.String() == "enter" {
				if m.username == "" {
					m.username = m.userPrompt.Value()
					m.userPrompt.Reset()
					m.channelsList.Focus()
				}
				return m, nil
			} else {
				if len(m.userPrompt.Value()) > 20 {
					return m, nil
				}
				c, cmd := m.userPrompt.Update(msg)
				m.userPrompt = c
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.host == "" {
		return styles.PromptStyle.Render(m.hostPrompt.View())
	}
	if m.username == "" {
		return styles.PromptStyle.Render(m.userPrompt.View())
	}

	channelsList := styles.ChannelListStyle.Render(m.channelsList.View())
	chatHistory := styles.ChatHistoryStyle.Render(m.chatHistory.View())
	messageInput := styles.MessageInputStyle.Render(m.messageInput.View())

	chatHistoryAndInput := lipgloss.JoinVertical(lipgloss.Top, chatHistory, messageInput)
	return lipgloss.JoinHorizontal(lipgloss.Top, channelsList, chatHistoryAndInput)
}
