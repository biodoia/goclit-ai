// App - Main TUI application with pane layout
package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/biodoia/goclit-ai/internal/providers"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bl "github.com/winder/bubblelayout"
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

	// Layout (bubblelayout)
	layout    bl.BubbleLayout
	layoutIDs LayoutIDs

	// Viewports for panes
	agentsVP viewport.Model
	chatVP   viewport.Model

	// Focus
	focusedPane int // 0=agents, 1=chat

	// Input
	input   textinput.Model
	spinner spinner.Model

	// Data
	messages      []Message
	agents        []AgentItem
	selectedAgent int
	agentRunning  bool

	// Provider
	provider     *providers.Client
	providerErr  string
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

	// Layout with bubblelayout
	layout, ids := NewLayout()

	// Auto-detect provider
	provider, providerErr := providers.AutoDetect()
	errMsg := ""
	if providerErr != nil {
		errMsg = providerErr.Error()
	}

	return App{
		state:     StateIntro,
		introTime: time.Now(),
		layout:    layout,
		layoutIDs: ids,
		agentsVP:  viewport.New(0, 0),
		chatVP:    viewport.New(0, 0),
		focusedPane: 1, // Start on chat
		input:     ti,
		spinner:   s,
		agents:    DefaultAgents(),
		messages: []Message{
			{Role: "system", Content: "Welcome to GOCLIT - The Dream CLI", Time: time.Now()},
		},
		provider:    provider,
		providerErr: errMsg,
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
				a.focusedPane = (a.focusedPane + 1) % 2
			}
		case "shift+tab":
			if a.state == StateMain {
				a.focusedPane = (a.focusedPane + 1) % 2
			}
		case "up", "k":
			if a.focusedPane == 0 && a.selectedAgent > 0 {
				a.selectedAgent--
				a.updateAgentsPane()
			}
		case "down", "j":
			if a.focusedPane == 0 && a.selectedAgent < len(a.agents)-1 {
				a.selectedAgent++
				a.updateAgentsPane()
			}
		case "enter":
			if a.focusedPane == 1 && a.input.Value() != "" {
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
			// Trigger resize to update layout
			return a, func() tea.Msg {
				return a.layout.Resize(a.width, a.height)
			}
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.input.Width = msg.Width - 6
		// Convert to bubblelayout message
		return a, func() tea.Msg {
			return a.layout.Resize(msg.Width, msg.Height)
		}

	case bl.BubbleLayoutMsg:
		// Update component sizes from layout
		if sz, err := msg.Size(a.layoutIDs.Agents); err == nil {
			a.agentsVP.Width = sz.Width - 4
			a.agentsVP.Height = sz.Height - 2
		}
		if sz, err := msg.Size(a.layoutIDs.Chat); err == nil {
			a.chatVP.Width = sz.Width - 4
			a.chatVP.Height = sz.Height - 2
		}
		a.updateAgentsPane()
		a.updateChatPane()

	case tickMsg:
		a.introFrame++
		if a.state == StateIntro {
			elapsed := time.Since(a.introTime)
			if elapsed > 2500*time.Millisecond {
				a.state = StateMain
				// Trigger resize to update layout
				return a, func() tea.Msg {
					return a.layout.Resize(a.width, a.height)
				}
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

	// Update viewports
	if a.state == StateMain {
		var cmd tea.Cmd
		if a.focusedPane == 0 {
			a.agentsVP, cmd = a.agentsVP.Update(msg)
		} else {
			a.chatVP, cmd = a.chatVP.Update(msg)
		}
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
	a.agentsVP.SetContent(content)
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

	a.chatVP.SetContent(strings.Join(lines, "\n"))
	a.chatVP.GotoBottom()
}

type agentResponseMsg struct {
	agent   string
	content string
}

func (a *App) processCommand(cmd string) tea.Cmd {
	a.agentRunning = true
	agent := a.agents[a.selectedAgent]

	// Check if provider is available
	if a.provider == nil {
		return func() tea.Msg {
			return agentResponseMsg{
				agent:   agent.Name,
				content: fmt.Sprintf("⚠️ No provider configured.\n\n%s\n\nSet OPENROUTER_API_KEY, ANTHROPIC_API_KEY, or start Ollama/GoBro.", a.providerErr),
			}
		}
	}

	// Real API call
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Build messages
		messages := []providers.Message{
			{Role: "system", Content: fmt.Sprintf("You are %s, a specialized AI agent. %s", agent.Name, agent.Role)},
			{Role: "user", Content: cmd},
		}

		response, err := a.provider.Chat(ctx, messages)
		if err != nil {
			return agentResponseMsg{
				agent:   agent.Name,
				content: fmt.Sprintf("❌ Error: %v", err),
			}
		}

		return agentResponseMsg{
			agent:   agent.Name,
			content: response,
		}
	}
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

	// Phase 2+: Logo with effects (Gopilot style)
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

	// Antenna flicker (longer phase)
	if progress > 0.3 && progress < 0.6 {
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

	// NO title here - it goes in the panes header

	// Listening (after 0.65)
	if progress > 0.65 {
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

	// Header
	sections = append(sections, a.renderHeader())

	// Panes
	sections = append(sections, a.renderPanes())

	// Input bar
	sections = append(sections, a.renderInputBar())

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Background(BgDark).
		Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (a App) renderHeader() string {
	logo := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		Render("⚡ GOCLIT")

	ver := lipgloss.NewStyle().
		Foreground(Gray500).
		Render(" v0.2.0")

	// Provider status
	providerStatus := ""
	if a.provider != nil {
		providerStatus = lipgloss.NewStyle().
			Foreground(Green).
			Render(" │ " + a.provider.ProviderName() + ":" + a.provider.Model())
	} else {
		providerStatus = lipgloss.NewStyle().
			Foreground(Red).
			Render(" │ No Provider")
	}

	return lipgloss.NewStyle().
		Width(a.width).
		Background(BgHighlight).
		Padding(0, 1).
		Render(logo + ver + providerStatus)
}

func (a App) renderPanes() string {
	// Agents pane
	agentsBorder := Gray700
	if a.focusedPane == 0 {
		agentsBorder = Cyan
	}

	agentsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(White).
		Render("AGENTS")

	if a.focusedPane == 0 {
		agentsTitle = lipgloss.NewStyle().Bold(true).Foreground(Cyan).Render("AGENTS")
	}

	agentsPane := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(agentsBorder).
		Width(a.agentsVP.Width + 4).
		Height(a.agentsVP.Height + 3).
		Render(agentsTitle + "\n\n" + a.agentsVP.View())

	// Chat pane
	chatBorder := Gray700
	if a.focusedPane == 1 {
		chatBorder = Cyan
	}

	chatTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(White).
		Render("CHAT")

	if a.focusedPane == 1 {
		chatTitle = lipgloss.NewStyle().Bold(true).Foreground(Cyan).Render("CHAT")
	}

	chatPane := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(chatBorder).
		Width(a.chatVP.Width + 4).
		Height(a.chatVP.Height + 3).
		Render(chatTitle + "\n\n" + a.chatVP.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, agentsPane, " ", chatPane)
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
