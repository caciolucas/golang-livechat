# Live Chat TUI

A terminal-based live chat application built in Go, featuring a Text User Interface (TUI) for seamless user interaction, real-time communication using WebSockets, and MongoDB for efficient data storage. This project was created as an experiment while exploring Golang, WebSocket protocols, and TUI design.

## Features

- **Multiple Channels:** Each server supports multiple chat channels, enabling users to join and communicate in specific topics.
- **Real-time Updates:** Messages are delivered and displayed in real-time using WebSockets.
- **Terminal UI:** An intuitive terminal-based interface built with Bubble Tea, Bubble, and Lipgloss for style.
- **Persistent Storage:** MongoDB is used to store user data, messages, and channels.
- **Scalable Design:** Built with Gonic Gin and Gorilla WebSocket for a robust backend.

---

## Technologies Used

### Backend
- **[Gorilla WebSocket](https://github.com/gorilla/websocket):** Enables WebSocket support for real-time communication.
- **[Gonic Gin](https://github.com/gin-gonic/gin):** Provides a fast and flexible HTTP server for handling API requests.
- **MongoDB:** Stores persistent data for users, messages, and channels.

### Frontend (TUI)
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea):** A Go framework for building terminal applications.
- **[Bubble](https://github.com/charmbracelet/bubbles):** A library of common TUI components for Bubble Tea.
- **[Lipgloss](https://github.com/charmbracelet/lipgloss):** Adds stylish formatting to the terminal interface.
