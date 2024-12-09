package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var channelListStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
var chatHistoryStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
var messageInputStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(3)
