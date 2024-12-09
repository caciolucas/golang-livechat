package components

import (
	"golang-chat/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChannelList struct {
	Channels        []models.Channel
	ActiveChannelID primitive.ObjectID
}

func (c ChannelList) Init() tea.Cmd {
	return nil
}

func (c ChannelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c ChannelList) View() string {
	return "Channel List"
}

