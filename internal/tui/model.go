package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"golang-chat/internal/models"
)

type model struct {
	channelsList   list.Model
	chatHistory    viewport.Model
	messageInput   textarea.Model
	currentChannel *models.Channel
}

func (m model) Init() tea.Cmd {
	return nil
}

func InitialModel() model {
	channels := fetchChannels()

	items := make([]list.Item, len(channels))
	for i := range channels {
		items[i] = channels[i]
	}

	channelsList := list.New(
		items, list.NewDefaultDelegate(), 0, 0,
	)
	channelsList.SetShowHelp(false)
	channelsList.SetShowStatusBar(false)
	channelsList.Title = "Available channels"

	messageInput := textarea.New()
	messageInput.Placeholder = "Type your message here..."
	messageInput.Prompt = ""
	messageInput.ShowLineNumbers = false
	messageInput.SetHeight(4)

	chatHistory := viewport.New(0, 0)
	m := model{channelsList: channelsList, chatHistory: chatHistory, messageInput: messageInput}

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if selectedChannel, ok := m.channelsList.SelectedItem().(models.Channel); ok {
				m.currentChannel = &selectedChannel
				messages := fetchMessages(selectedChannel.ID)
				m.chatHistory.SetContent(formatMessages(messages))
			}
			m.messageInput.Focus()
		}
	case tea.WindowSizeMsg:
		m = resizeTUI(m, msg)
	}

	var cmd tea.Cmd
	m.channelsList, cmd = m.channelsList.Update(msg)
	return m, cmd
}

func (m model) View() string {
	channelsList := channelListStyle.Render(m.channelsList.View())
	chatHistory := chatHistoryStyle.Render(m.chatHistory.View())
	messageInput := messageInputStyle.Render(m.messageInput.View())

	chatHistoryAndInput := lipgloss.JoinVertical(lipgloss.Top, chatHistory, messageInput)
	return lipgloss.JoinHorizontal(lipgloss.Top, channelsList, chatHistoryAndInput)
}
