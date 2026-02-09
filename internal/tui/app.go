// App - Main TUI application with pane layout
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
)

// Message represents a chat message
type Message struct {
	Role    string
	Content string
	Time    time.Time
	Agent   string
}

// App is the main TUI model
type App struct {
	width  int
	height int

	state     AppState
	introTime time.Time
	introFrame int

	// Pane layout
	panes *PaneLayout

	// Input
	input   textinput.Model
	spinner spinner.Model

	// Data
	messages     []Message
	agents       []AgentItem
	selectedAgent int
	agentRunning bool
}

func NewApp() App {
	// Text input
	ti := textinput.New()
	ti.Placeholder = "Ask anything or type a command..."
	ti.CharLimit = 500
	ti.Width = 60
	ti.Focus()

	// Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	return App{
		state:     StateIntro,
		introTime: time.Now(),
		panes:     NewPaneLayout(),
		input:     ti,
		spinner:   s,
		agents:    DefaultAgents(),
		messages: []Message{
			{Role: "system", Content: "Welcome to GOCLIT - The Dream CLI", Time: time.Now()},
		},
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
			if a.state == StateMain {
				a.panes.FocusNext()
			}
		case "shift+tab":
			if a.state == StateMain {
				a.panes.FocusPrev()
			}
		case "up", "k":
			if a.panes.focusedPane == PaneAgents && a.selectedAgent > 0 {
				a.selectedAgent--
				a.updateAgentsPane()
			}
		case "down", "j":
			if a.panes.focusedPane == PaneAgents && a.selectedAgent < len(a.agents)-1 {
				a.selectedAgent++
				a.updateAgentsPane()
			}
		case "enter":
			if a.panes.focusedPane == PaneChat && a.input.Value() != "" {
				userMsg := a.input.Value()
				a.messages = append(a.messages, Message{
					Role:    "user",
					Content: userMsg,
					Time:    time.Now(),
				})
				a.input.Reset()
				a.updateChatPane()
				cmds = append(cmds, a.processCommand(userMsg))
			}
		}

		// Skip intro on any key
		if a.state == StateIntro {
			a.state = StateMain
			a.panes.SetSize(a.width, a.height-3) // Reserve space for input
			a.updateAgentsPane()
			a.updateChatPane()
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.input.Width = msg.Width - 6
		if a.state == StateMain {
			a.panes.SetSize(a.width, a.height-3)
			a.updateAgentsPane()
			a.updateChatPane()
		}

	case tickMsg:
		a.introFrame++
		if a.state == StateIntro {
			elapsed := time.Since(a.introTime)
			if elapsed > 2500*time.Millisecond {
				a.state = StateMain
				a.panes.SetSize(a.width, a.height-3)
				a.updateAgentsPane()
				a.updateChatPane()
			}
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
		a.updateChatPane()
	}

	// Update panes
	if a.state == StateMain {
		cmd := a.panes.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update text input
	var cmd tea.Cmd
	a.input, cmd = a.input.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) updateAgentsPane() {
	content := RenderAgentList(a.agents, a.selectedAgent)
	a.panes.SetContent(PaneAgents, content)
}

func (a *App) updateChatPane() {
	var lines []string

	for _, msg := range a.messages {
		switch msg.Role {
		case "user":
			prefix := lipgloss.NewStyle().Foreground(Blue).Bold(true).Render("You: ")
			text := lipgloss.NewStyle().Foreground(White).Render(msg.Content)
			lines = append(lines, prefix+text, "")

		case "assistant":
			agentStyle := lipgloss.NewStyle().Foreground(Cyan).Bold(true)
			if msg.Agent != "" {
				lines = append(lines, agentStyle.Render(msg.Agent+":"))
			}
			text := lipgloss.NewStyle().Foreground(Gray300).Render(msg.Content)
			lines = append(lines, text, "")

		case "system":
			sysStyle := lipgloss.NewStyle().Foreground(Gray500).Italic(true)
			lines = append(lines, sysStyle.Render("• "+msg.Content), "")
		}
	}

	a.panes.SetContent(PaneChat, strings.Join(lines, "\n"))
}

type agentResponseMsg struct {
	agent   string
	content string
}

func (a *App) processCommand(cmd string) tea.Cmd {
	a.agentRunning = true
	agent := a.agents[a.selectedAgent]

	return tea.Tick(time.Millisecond*800, func(t time.Time) tea.Msg {
		response := fmt.Sprintf("Received: %s\n\n⚠️ Provider not configured. Run: goclit config --provider <name>", cmd)

		if strings.HasPrefix(strings.ToLower(cmd), "ultrawork") {
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
		return a.renderIntro()
	}

	return a.renderMain()
}

func (a App) renderIntro() string {
	// Animated intro sequence
	progress := float64(time.Since(a.introTime)) / float64(2500*time.Millisecond)

	var content strings.Builder

	// Phase 1: Black screen (0-0.12)
	if progress < 0.12 {
		return lipgloss.NewStyle().
			Width(a.width).
			Height(a.height).
			Background(BgDark).
			Render("")
	}

	// Phase 2+: Logo with effects
	logoColor := Gradient(clamp((progress-0.12)/0.5, 0, 1))
	logoStyle := lipgloss.NewStyle().
		Foreground(logoColor).
		Bold(true)

	// Flicker during early phase
	showLogo := true
	if progress < 0.3 && a.introFrame%4 == 0 {
		showLogo = false
	}

	logo := []string{
		"      ★      ",
		"   ▄▄▄▄▄▄▄   ",
		"   █ ◉ ◉ █   ",
		"   █  ▼  █   ",
		"   █ ╰─╯ █   ",
		"   ▀▀▀▀▀▀▀   ",
	}

	// Antenna flicker
	if progress > 0.3 && progress < 0.5 {
		if a.introFrame%3 == 0 {
			logo[0] = "             "
		}
	}

	if showLogo {
		for _, line := range logo {
			content.WriteString(logoStyle.Render(line) + "\n")
		}
	} else {
		for range logo {
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")

	// Title (letter by letter after 0.5)
	if progress > 0.5 {
		title := "G O C L I T"
		titleProgress := clamp((progress-0.5)/0.2, 0, 1)
		visibleChars := int(titleProgress * float64(len(title)))

		titleStyle := lipgloss.NewStyle().
			Foreground(White).
			Bold(true)

		content.WriteString(titleStyle.Render(title[:visibleChars]))
		if visibleChars < len(title) && a.introFrame%8 < 4 {
			content.WriteString("█")
		}
		content.WriteString("\n\n")
	}

	// Tagline (after 0.7)
	if progress > 0.7 {
		tagStyle := lipgloss.NewStyle().Foreground(Gray500).Italic(true)
		content.WriteString(tagStyle.Render("The Dream CLI") + "\n")
		content.WriteString(lipgloss.NewStyle().Foreground(Cyan).Render("v0.2.0") + "\n\n")
	}

	// Listening (after 0.85)
	if progress > 0.85 {
		sparkles := []string{"✨", "⚡", "💫", "🌟"}
		s := sparkles[a.introFrame/3%len(sparkles)]
		listenStyle := lipgloss.NewStyle().Foreground(Cyan)
		content.WriteString(listenStyle.Render(s+" Agents are listening... "+s))
	}

	centered := lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, content.String())
	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Background(BgDark).
		Render(centered)
}

func (a App) renderMain() string {
	var sections []string

	// Panes
	sections = append(sections, a.panes.View())

	// Input bar
	sections = append(sections, a.renderInputBar())

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Background(BgDark).
		Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (a App) renderInputBar() string {
	prompt := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		Render("❯ ")

	status := ""
	if a.agentRunning {
		status = a.spinner.View() + " "
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(Cyan).
		Width(a.width - 6).
		Padding(0, 1)

	return lipgloss.NewStyle().
		Width(a.width).
		Padding(0, 1).
		Render(status + prompt + inputStyle.Render(a.input.View()))
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
