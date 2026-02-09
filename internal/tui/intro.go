// Intro animation - Logo with flicker effect
package tui

import (
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	animationDuration = 2500 * time.Millisecond
	fps               = 30
	frameTime         = time.Second / fps
)

// Robot logo ASCII art
var robotLogo = []string{
	"      ★      ",
	"   ▄▄▄▄▄▄▄   ",
	"   █ ◉ ◉ █   ",
	"   █  ▼  █   ",
	"   █ ╰─╯ █   ",
	"   ▀▀▀▀▀▀▀   ",
}

var robotLogoNoAntenna = []string{
	"             ",
	"   ▄▄▄▄▄▄▄   ",
	"   █ ◉ ◉ █   ",
	"   █  ▼  █   ",
	"   █ ╰─╯ █   ",
	"   ▀▀▀▀▀▀▀   ",
}

const goclitText = "G O C L I T"
const tagline = "The Dream CLI"
const listening = "✨ Agents are listening... ✨"
const version = "v0.2.0"

// IntroModel handles the intro animation
type IntroModel struct {
	width     int
	height    int
	startTime time.Time
	progress  float64
	done      bool
	showPanes bool
	frame     int
}

type tickMsg time.Time

func NewIntro(width, height int) IntroModel {
	return IntroModel{
		width:     width,
		height:    height,
		startTime: time.Now(),
	}
}

func (m IntroModel) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		tea.EnterAltScreen,
	)
}

func tick() tea.Cmd {
	return tea.Tick(frameTime, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m IntroModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.done {
			m.done = true
			m.showPanes = true
			return m, nil
		}
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.frame++
		elapsed := time.Since(m.startTime)
		m.progress = float64(elapsed) / float64(animationDuration)

		if m.progress >= 1.0 {
			m.progress = 1.0
			m.done = true
			m.showPanes = true
			return m, nil
		}

		return m, tick()
	}

	return m, nil
}

func (m IntroModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	if m.showPanes {
		return m.renderWithPanes()
	}

	return m.renderAnimation()
}

func (m IntroModel) renderAnimation() string {
	// Phases:
	// 0.00-0.15: Black screen
	// 0.15-0.35: Logo fades in with flicker
	// 0.35-0.55: Antenna flickers
	// 0.55-0.75: GOCLIT appears letter by letter
	// 0.75-0.90: Tagline fades in
	// 0.90-1.00: "Agents are listening" appears

	var content strings.Builder

	// Background style
	bgStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(BgDark).
		Align(lipgloss.Center, lipgloss.Center)

	// Phase 1: Black screen (0.00-0.15)
	if m.progress < 0.15 {
		return bgStyle.Render("")
	}

	// Phase 2: Logo with flicker (0.15-0.35)
	logoVisible := true
	if m.progress >= 0.15 && m.progress < 0.35 {
		// Random flicker during fade-in
		flickerChance := 0.3 - (m.progress-0.15)*1.5 // Decreases over time
		if rand.Float64() < flickerChance {
			logoVisible = false
		}
	}

	// Phase 3: Antenna flicker (0.35-0.55)
	antennaOn := true
	if m.progress >= 0.35 && m.progress < 0.55 {
		// Antenna flickers rapidly
		flickerSpeed := 8.0
		antennaOn = int(m.progress*flickerSpeed*10)%2 == 0
	}

	// Render logo
	if logoVisible && m.progress >= 0.15 {
		logoColor := Gradient(clamp((m.progress-0.15)/0.4, 0, 1))
		logoStyle := lipgloss.NewStyle().
			Foreground(logoColor).
			Bold(true)

		var logo []string
		if antennaOn {
			logo = robotLogo
		} else {
			logo = robotLogoNoAntenna
		}

		for _, line := range logo {
			// Glitch effect during flicker phase
			if m.progress < 0.35 && rand.Float64() < 0.1 {
				line = glitchLine(line)
			}
			content.WriteString(logoStyle.Render(line) + "\n")
		}
		content.WriteString("\n")
	}

	// Phase 4: GOCLIT letter by letter (0.55-0.75)
	if m.progress >= 0.55 {
		titleProgress := clamp((m.progress-0.55)/0.20, 0, 1)
		visibleChars := int(titleProgress * float64(len(goclitText)))

		titleStyle := lipgloss.NewStyle().
			Foreground(White).
			Bold(true)

		visibleText := goclitText[:visibleChars]

		// Cursor blink at end
		if visibleChars < len(goclitText) && m.frame%10 < 5 {
			visibleText += "█"
		}

		content.WriteString(titleStyle.Render(visibleText) + "\n\n")
	}

	// Phase 5: Tagline (0.75-0.90)
	if m.progress >= 0.75 {
		taglineOpacity := easeOutExpo(clamp((m.progress-0.75)/0.15, 0, 1))
		taglineStyle := lipgloss.NewStyle().
			Foreground(Gray500).
			Italic(true)

		if taglineOpacity > 0.3 {
			content.WriteString(taglineStyle.Render(tagline) + "\n")
			content.WriteString(lipgloss.NewStyle().Foreground(Cyan).Render(version) + "\n\n")
		}
	}

	// Phase 6: Agents listening (0.90-1.00)
	if m.progress >= 0.90 {
		listenOpacity := easeOutExpo(clamp((m.progress-0.90)/0.10, 0, 1))
		listenStyle := lipgloss.NewStyle().
			Foreground(Cyan)

		if listenOpacity > 0.3 {
			// Sparkle animation
			sparkles := []string{"✨", "⚡", "💫", "🌟"}
			sparkleIdx := m.frame / 5 % len(sparkles)
			text := strings.Replace(listening, "✨", sparkles[sparkleIdx], 2)
			content.WriteString(listenStyle.Render(text))
		}
	}

	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content.String())
	return bgStyle.Render(centered)
}

func glitchLine(line string) string {
	// Add random glitch characters
	runes := []rune(line)
	glitchChars := []rune{'█', '▓', '▒', '░', '▀', '▄', '│', '─'}

	for i := range runes {
		if rand.Float64() < 0.2 && runes[i] != ' ' {
			runes[i] = glitchChars[rand.Intn(len(glitchChars))]
		}
	}
	return string(runes)
}

func (m IntroModel) renderWithPanes() string {
	var b strings.Builder

	// Header with logo (smaller, left-aligned)
	headerLogo := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		Render("⚡ GOCLIT")

	headerVersion := lipgloss.NewStyle().
		Foreground(Gray500).
		Render(" " + version)

	header := lipgloss.NewStyle().
		Width(m.width).
		Background(BgHighlight).
		Padding(0, 2).
		Render(headerLogo + headerVersion)

	b.WriteString(header + "\n\n")

	// Calculate pane dimensions
	paneWidth := (m.width - 6) / 2
	paneHeight := m.height - 8

	// Left pane - Agents
	agentsPane := m.renderAgentsPane(paneWidth, paneHeight)

	// Right pane - Chat/Output
	chatPane := m.renderChatPane(paneWidth, paneHeight)

	// Join panes horizontally
	panes := lipgloss.JoinHorizontal(lipgloss.Top, agentsPane, "  ", chatPane)
	b.WriteString(panes)

	// Status bar at bottom
	statusBar := m.renderStatusBar()
	b.WriteString("\n" + statusBar)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(BgDark).
		Render(b.String())
}

func (m IntroModel) renderAgentsPane(width, height int) string {
	agents := []struct {
		icon  string
		name  string
		color lipgloss.Color
	}{
		{"⚙️", "Sisyphus", Purple},
		{"🔨", "Hephaestus", Blue},
		{"🔮", "Oracle", Cyan},
		{"📚", "Librarian", Green},
		{"🎨", "Frontend", Pink},
		{"⚡", "Backend", Orange},
		{"🔧", "DevOps", Yellow},
	}

	var content strings.Builder
	titleStyle := lipgloss.NewStyle().
		Foreground(White).
		Bold(true)

	content.WriteString(titleStyle.Render("AGENTS") + "\n\n")

	for _, a := range agents {
		badge := AgentBadge(a.color).Render(a.icon)
		name := lipgloss.NewStyle().Foreground(Gray300).Render(a.name)
		content.WriteString(badge + " " + name + "\n")
	}

	return ActivePanelStyle.
		Width(width).
		Height(height).
		Render(content.String())
}

func (m IntroModel) renderChatPane(width, height int) string {
	var content strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(White).
		Bold(true)

	content.WriteString(titleStyle.Render("CHAT") + "\n\n")

	welcomeStyle := lipgloss.NewStyle().
		Foreground(Gray500).
		Italic(true)

	content.WriteString(welcomeStyle.Render("Type a command or ask anything...\n"))
	content.WriteString(welcomeStyle.Render("Try: ultrawork \"build a REST API\""))

	return PanelStyle.
		Width(width).
		Height(height).
		Render(content.String())
}

func (m IntroModel) renderStatusBar() string {
	left := lipgloss.NewStyle().
		Foreground(Gray500).
		Render("Press ? for help • q to quit")

	right := lipgloss.NewStyle().
		Foreground(Cyan).
		Render("⚡ Provider: claude")

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 0 {
		gap = 0
	}

	return StatusStyle.
		Width(m.width).
		Render(left + strings.Repeat(" ", gap) + right)
}

// Easing functions

func easeOutExpo(t float64) float64 {
	if t >= 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Run starts the intro animation
func RunIntro() error {
	p := tea.NewProgram(
		NewIntro(80, 24),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
