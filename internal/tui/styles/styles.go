package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var DocStyle = lipgloss.NewStyle().Margin(1, 2)

// BASE COLORS AND STYLES
var border = lipgloss.NormalBorder()
var primaryColor = lipgloss.Color("4")
var secondaryColor = lipgloss.Color("2")
var componentStyle = lipgloss.NewStyle().Border(border).
	BorderForeground(primaryColor)

// COMPONENTS STYLES
var ChannelListStyle = componentStyle
var ChatHistoryStyle = componentStyle
var MessageInputStyle = componentStyle.Height(4)

// CHANNEL LIST STYLES
var ActiveChannelStyle = lipgloss.NewStyle().
	Foreground(primaryColor).
	Bold(true)
var HoveredChannelStyle = lipgloss.NewStyle().
	Foreground(secondaryColor).
	Bold(true)

// MESSAGE INPUT STYLES
var MessageInput = lipgloss.NewStyle().Foreground(primaryColor)
var Prompt = lipgloss.NewStyle().Foreground(primaryColor)
