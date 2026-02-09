// App - Main TUI application
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppState represents the current state
type AppState int

const (
	StateIntro AppState = iota
	StateMain
	StateAgentSelect
	StateRunning
)

// Message represents a chat message
type Message struct {
	Role    string // "user", "assistant", "system"
	Content string
	Time    time.Time
	Agent   string // Which agent responded
}

// App is the main TUI model
type App struct {
	// Dimensions
	width  int
	height int

	// State
	state     AppState
	intro     IntroModel
	introTime time.Time

	// UI Components
	input   textinput.Model
	spinner spinner.Model

	// Chat
	messages    []Message
	currentTask string

	// Agents
	agents       []AgentInfo
	activeAgent  int
	agentRunning bool

	// Focus
	focusPane int // 0=agents, 1=chat, 2=input
}

// AgentInfo describes an agent
type AgentInfo struct {
	Icon   string
	Name   string
	Role   string
	Color  lipgloss.Color
	Active bool
}

// Initialize default agents
var defaultAgents = []AgentInfo{
	{"⚙️", "Sisyphus", "Discipline", Purple, false},
	{"🔨", "Hephaestus", "Autonomy", Blue, false},
	{"🔮", "Oracle", "Knowledge", Cyan, false},
	{"📚", "Librarian", "Search", Green, false},
	{"🎨", "Frontend", "UI/UX", Pink, false},
	{"⚡", "Backend", "Server", Orange, false},
	{"🔧", "DevOps", "Infra", Yellow, false},
}

func NewApp() App {
	// Text input
	ti := textinput.New()
	ti.Placeholder = "Ask anything or type a command..."
	ti.CharLimit = 500
	ti.Width = 60

	// Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	return App{
		state:     StateIntro,
		intro:     NewIntro(80, 24),
		introTime: time.Now(),
		input:     ti,
		spinner:   s,
		agents:    defaultAgents,
		messages: []Message{
			{
				Role:    "system",
				Content: "Welcome to GOCLIT - The Dream CLI",
				Time:    time.Now(),
			},
		},
		focusPane: 2, // Start focused on input
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		textinput.Blink,
		a.spinner.Tick,
		tick(),
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if a.state == StateMain && !a.agentRunning {
				return a, tea.Quit
			}
		case "tab":
			// Cycle focus between panes
			a.focusPane = (a.focusPane + 1) % 3
		case "enter":
			if a.focusPane == 2 && a.input.Value() != "" {
				// Submit message
				userMsg := a.input.Value()
				a.messages = append(a.messages, Message{
					Role:    "user",
					Content: userMsg,
					Time:    time.Now(),
				})
				a.input.Reset()
				a.currentTask = userMsg

				// Simulate agent response
				cmds = append(cmds, a.processCommand(userMsg))
			}
		case "up", "down":
			if a.focusPane == 0 {
				// Navigate agents
				if msg.String() == "up" && a.activeAgent > 0 {
					a.activeAgent--
				}
				if msg.String() == "down" && a.activeAgent < len(a.agents)-1 {
					a.activeAgent++
				}
			}
		}

		// Skip intro on any key
		if a.state == StateIntro {
			a.state = StateMain
			a.input.Focus()
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.intro.width = msg.Width
		a.intro.height = msg.Height
		a.input.Width = msg.Width - 10

	case tickMsg:
		if a.state == StateIntro {
			// Check if intro animation is done
			if time.Since(a.introTime) > animationDuration {
				a.state = StateMain
				a.input.Focus()
			}
			newIntro, cmd := a.intro.Update(msg)
			a.intro = newIntro.(IntroModel)
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, tick())

	case spinner.TickMsg:
		var cmd tea.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case agentResponseMsg:
		a.messages = append(a.messages, Message{
			Role:    "assistant",
			Content: msg.content,
			Time:    time.Now(),
			Agent:   msg.agent,
		})
		a.agentRunning = false
	}

	// Update text input if focused
	if a.focusPane == 2 {
		var cmd tea.Cmd
		a.input, cmd = a.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

type agentResponseMsg struct {
	agent   string
	content string
}

func (a App) processCommand(cmd string) tea.Cmd {
	a.agentRunning = true

	// Simulate processing
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		agent := a.agents[a.activeAgent]

		// Simple command parsing
		response := fmt.Sprintf("Agent %s received: %s\n\n⚠️ Provider not configured yet.", agent.Name, cmd)

		if strings.HasPrefix(cmd, "ultrawork") {
			response = "🚀 ULTRAWORK MODE\n━━━━━━━━━━━━━━━━━━\nAll agents coordinating...\n\n⚠️ Connect a provider first:\ngoclit config --provider claude"
		}

		return agentResponseMsg{
			agent:   agent.Name,
			content: response,
		}
	})
}

func (a App) View() string {
	if a.width == 0 || a.height == 0 {
		return ""
	}

	if a.state == StateIntro {
		return a.intro.View()
	}

	return a.renderMain()
}

func (a App) renderMain() string {
	var b strings.Builder

	// Header
	header := a.renderHeader()
	b.WriteString(header + "\n")

	// Calculate dimensions
	contentHeight := a.height - 5 // Header + status bar + padding
	agentsPaneWidth := 24
	chatPaneWidth := a.width - agentsPaneWidth - 4

	// Panes
	agentsPane := a.renderAgentsPane(agentsPaneWidth, contentHeight)
	chatPane := a.renderChatPane(chatPaneWidth, contentHeight)

	panes := lipgloss.JoinHorizontal(lipgloss.Top, agentsPane, " ", chatPane)
	b.WriteString(panes + "\n")

	// Input bar
	inputBar := a.renderInputBar()
	b.WriteString(inputBar)

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Background(BgDark).
		Render(b.String())
}

func (a App) renderHeader() string {
	logo := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		Render("⚡ GOCLIT")

	ver := lipgloss.NewStyle().
		Foreground(Gray500).
		Render(" v0.2.0")

	status := ""
	if a.agentRunning {
		status = a.spinner.View() + " Running..."
	}

	right := lipgloss.NewStyle().
		Foreground(Cyan).
		Render(status)

	gap := a.width - lipgloss.Width(logo+ver) - lipgloss.Width(right) - 4
	if gap < 0 {
		gap = 1
	}

	return lipgloss.NewStyle().
		Width(a.width).
		Background(BgHighlight).
		Padding(0, 1).
		Render(logo + ver + strings.Repeat(" ", gap) + right)
}

func (a App) renderAgentsPane(width, height int) string {
	var content strings.Builder

	title := lipgloss.NewStyle().
		Foreground(White).
		Bold(true).
		Render("AGENTS")

	content.WriteString(title + "\n\n")

	for i, agent := range a.agents {
		isActive := i == a.activeAgent
		isFocused := a.focusPane == 0

		// Icon with color
		iconStyle := lipgloss.NewStyle().MarginRight(1)
		if isActive && isFocused {
			iconStyle = iconStyle.Background(agent.Color).Foreground(Black)
		}
		icon := iconStyle.Render(agent.Icon)

		// Name
		nameStyle := lipgloss.NewStyle().Foreground(Gray300)
		if isActive {
			nameStyle = nameStyle.Foreground(White).Bold(true)
		}
		name := nameStyle.Render(agent.Name)

		// Cursor
		cursor := "  "
		if isActive {
			cursor = lipgloss.NewStyle().Foreground(Cyan).Render("▸ ")
		}

		content.WriteString(cursor + icon + name + "\n")
	}

	// Add help at bottom
	content.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().Foreground(Gray500).Italic(true)
	content.WriteString(helpStyle.Render("↑↓ select\n⏎  activate"))

	style := PanelStyle.Width(width).Height(height)
	if a.focusPane == 0 {
		style = ActivePanelStyle.Width(width).Height(height)
	}

	return style.Render(content.String())
}

func (a App) renderChatPane(width, height int) string {
	var content strings.Builder

	title := lipgloss.NewStyle().
		Foreground(White).
		Bold(true).
		Render("CHAT")

	content.WriteString(title + "\n\n")

	// Render messages (last N that fit)
	maxMessages := height - 6
	start := 0
	if len(a.messages) > maxMessages {
		start = len(a.messages) - maxMessages
	}

	for _, msg := range a.messages[start:] {
		switch msg.Role {
		case "user":
			prefix := lipgloss.NewStyle().Foreground(Blue).Bold(true).Render("You: ")
			text := lipgloss.NewStyle().Foreground(White).Render(msg.Content)
			content.WriteString(prefix + text + "\n\n")

		case "assistant":
			agentStyle := lipgloss.NewStyle().Foreground(Cyan).Bold(true)
			if msg.Agent != "" {
				content.WriteString(agentStyle.Render(msg.Agent+": "))
			}
			text := lipgloss.NewStyle().Foreground(Gray300).Render(msg.Content)
			content.WriteString(text + "\n\n")

		case "system":
			sysStyle := lipgloss.NewStyle().Foreground(Gray500).Italic(true)
			content.WriteString(sysStyle.Render("• "+msg.Content) + "\n\n")
		}
	}

	style := PanelStyle.Width(width).Height(height)
	if a.focusPane == 1 {
		style = ActivePanelStyle.Width(width).Height(height)
	}

	return style.Render(content.String())
}

func (a App) renderInputBar() string {
	prompt := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		Render("❯ ")

	inputStyle := InputStyle.Width(a.width - 6)
	if a.focusPane == 2 {
		inputStyle = inputStyle.BorderForeground(Cyan)
	}

	return lipgloss.NewStyle().
		Width(a.width).
		Padding(0, 1).
		Render(prompt + inputStyle.Render(a.input.View()))
}

// Run starts the TUI application
func Run() error {
	p := tea.NewProgram(
		NewApp(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
