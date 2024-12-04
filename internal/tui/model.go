package tui

import (
	"encoding/json"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"golang-chat/internal/models"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

type model struct {
	channelsList list.Model
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

func InitialModel() model {
	channels := fetchChannels()

	items := make([]list.Item, len(channels))
	for i := range channels {
		items[i] = channels[i]
	}

	terminalWidth, terminalHeigh, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal("Error getting terminal size: ", err)
	}

	channelList := list.New(items, list.NewDefaultDelegate(), terminalWidth, terminalHeigh-1)
	channelList.Title = "Channels"
	channelList.SetShowStatusBar(true)
	channelList.Styles.TitleBar = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	// Return the model with the list
	return model{
		channelsList: channelList,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		m.channelsList, _ = m.channelsList.Update(msg)
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return m.channelsList.View()
}
