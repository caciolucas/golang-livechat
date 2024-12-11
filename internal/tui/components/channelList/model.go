package channellist

import (
	"fmt"
	"golang-chat/internal/models"
	"golang-chat/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Channels        []models.Channel
	SelectedChannel *models.Channel
	cursor          int
	focus           bool
}

func (c Model) Init() tea.Cmd {
	return nil
}

func (c Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return c, tea.Quit
		case "j", "down":
			if c.cursor < len(c.Channels)-1 {
				c.cursor++
			}
		case "k", "up":
			if c.cursor > 0 {
				c.cursor--
			}
		case "enter":
			c.SelectedChannel = &c.Channels[c.cursor]
			c.focus = false
		}
	}
	return c, nil
}

func (c Model) View() string {
	view := "Channels\n"
	for i, channel := range c.Channels {
		style := lipgloss.NewStyle()
		hoveredIndicator, activeIndicator := " ", " "

		if c.SelectedChannel != nil && c.SelectedChannel.ID == channel.ID {
			style = styles.ActiveChannelStyle
		}
		if i == c.cursor && c.focus {
			style = styles.HoveredChannelStyle
			hoveredIndicator = " →"
		}

		view += style.Render(fmt.Sprintf("\n%s%s%s", hoveredIndicator, activeIndicator, channel.Name))
	}
	return fmt.Sprintf(" %s\n ", view)
}

func (c *Model) Focus() {
	c.focus = true
}

func (c *Model) Blur() {
	c.focus = false
}

func (c *Model) Focused() bool {
	return c.focus
}

func (c *Model) SetCursor(cursor int) {
	c.cursor = cursor
}

func (c *Model) Cursor() int {
	return c.cursor
}
