// Package tui - Main TUI layout with panes
// Flow: Black → Logo animation → Panes appear → Title
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Pane represents a UI panel
type Pane struct {
	Title    string
	Content  string
	Width    int
	Height   int
	Focused  bool
	Style    lipgloss.Style
}

// Layout represents the full TUI layout
type Layout struct {
	Width      int
	Height     int
	TitleBar   string
	LeftPane   Pane  // Agents/Navigation
	MainPane   Pane  // Chat/Output
	RightPane  Pane  // Context/Status (optional)
	StatusBar  string
	ShowRight  bool
}

// Color roles
var (
	primaryColor   = lipgloss.Color("39")  // Cyan
	secondaryColor = lipgloss.Color("245") // Gray
	accentColor    = lipgloss.Color("213") // Pink
	successColor   = lipgloss.Color("82")  // Green
	warningColor   = lipgloss.Color("220") // Yellow
	errorColor     = lipgloss.Color("196") // Red
	borderColor    = lipgloss.Color("240") // Dark gray
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(lipgloss.Color("235")).
			Padding(0, 2)

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)

	focusedPaneStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor)

	statusStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Background(lipgloss.Color("235"))
)

// NewLayout creates a new TUI layout
func NewLayout(width, height int) *Layout {
	leftWidth := width / 4
	mainWidth := width - leftWidth - 2

	return &Layout{
		Width:  width,
		Height: height,
		TitleBar: "GOCLIT",
		LeftPane: Pane{
			Title:   "🤖 Agents",
			Width:   leftWidth,
			Height:  height - 4,
			Focused: false,
		},
		MainPane: Pane{
			Title:   "💬 Chat",
			Width:   mainWidth,
			Height:  height - 4,
			Focused: true,
		},
		StatusBar: "Ready | ultrawork | Ctrl+C to quit",
		ShowRight: false,
	}
}

// Render returns the full layout as a string
func (l *Layout) Render() string {
	var sb strings.Builder

	// Title bar
	title := titleStyle.Width(l.Width).Render(
		fmt.Sprintf(" %s  ✨ Agents are listening...", l.TitleBar),
	)
	sb.WriteString(title)
	sb.WriteString("\n")

	// Panes
	leftStyle := paneStyle
	mainStyle := paneStyle
	if l.LeftPane.Focused {
		leftStyle = focusedPaneStyle
	}
	if l.MainPane.Focused {
		mainStyle = focusedPaneStyle
	}

	// Left pane content
	leftContent := l.renderPaneContent(l.LeftPane)
	leftBox := leftStyle.
		Width(l.LeftPane.Width).
		Height(l.LeftPane.Height).
		Render(leftContent)

	// Main pane content
	mainContent := l.renderPaneContent(l.MainPane)
	mainBox := mainStyle.
		Width(l.MainPane.Width).
		Height(l.MainPane.Height).
		Render(mainContent)

	// Join horizontally
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, mainBox)
	sb.WriteString(panes)
	sb.WriteString("\n")

	// Status bar
	status := statusStyle.Width(l.Width).Render(" " + l.StatusBar)
	sb.WriteString(status)

	return sb.String()
}

func (l *Layout) renderPaneContent(p Pane) string {
	var sb strings.Builder
	
	// Pane title
	titleLine := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Render(p.Title)
	sb.WriteString(titleLine)
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", p.Width-4))
	sb.WriteString("\n\n")

	// Content
	if p.Content != "" {
		sb.WriteString(p.Content)
	}

	return sb.String()
}

// SetAgentsList sets the left pane with agent list
func (l *Layout) SetAgentsList(agents []string) {
	var content strings.Builder
	for i, agent := range agents {
		if i == 0 {
			content.WriteString("► ")
		} else {
			content.WriteString("  ")
		}
		content.WriteString(agent)
		content.WriteString("\n")
	}
	l.LeftPane.Content = content.String()
}

// SetChatContent sets the main pane content
func (l *Layout) SetChatContent(content string) {
	l.MainPane.Content = content
}

// SetStatus sets the status bar
func (l *Layout) SetStatus(status string) {
	l.StatusBar = status
}

// AnimatePanesIn shows panes appearing with animation
func AnimatePanesIn(width, height int) {
	layout := NewLayout(width, height)
	
	// Set initial content
	layout.SetAgentsList([]string{
		"⚙️  Sisyphus",
		"🔨 Hephaestus",
		"🔮 Oracle",
		"📚 Librarian",
		"🎨 Frontend",
		"⚡ Backend",
		"🔧 DevOps",
	})
	
	layout.SetChatContent("Welcome to goclit!\n\nType a message or use:\n• ultrawork <task>\n• /agents\n• /help")

	// Animation: reveal line by line
	lines := strings.Split(layout.Render(), "\n")
	
	fmt.Print("\033[2J\033[H") // Clear
	
	for i, line := range lines {
		fmt.Printf("\033[%d;0H%s", i+1, line)
		time.Sleep(20 * time.Millisecond)
	}
}

// PlayFullStartup plays the complete startup sequence
func PlayFullStartup(width, height int) {
	// Phase 1: Logo animation
	PlayFullIntro(width, height)
	
	// Phase 2: Panes appear
	AnimatePanesIn(width, height)
}

// DefaultAgents returns the default agent list
func DefaultAgents() []string {
	return []string{
		"⚙️  Sisyphus",
		"🔨 Hephaestus", 
		"🔮 Oracle",
		"📚 Librarian",
		"🎨 Frontend Eng",
		"⚡ Backend Eng",
		"🔧 DevOps Eng",
	}
}
