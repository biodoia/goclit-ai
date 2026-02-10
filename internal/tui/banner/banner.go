// Package banner provides animated ASCII banner for goclit
// GOAL: Better than GitHub Copilot CLI banner
package banner

import (
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps           = 30
	frameTime     = time.Second / fps
	totalDuration = 4 * time.Second
)

// ANSI Color roles (semantic, not hardcoded)
var (
	// Primary gradient: purple → cyan → blue
	colorPrimary   = lipgloss.Color("#A855F7") // Purple
	colorSecondary = lipgloss.Color("#06B6D4") // Cyan
	colorAccent    = lipgloss.Color("#3B82F6") // Blue
	colorHighlight = lipgloss.Color("#F59E0B") // Amber
	colorGlow      = lipgloss.Color("#22D3EE") // Bright cyan
	colorEye       = lipgloss.Color("#10B981") // Emerald (eyes)
	colorShadow    = lipgloss.Color("#6366F1") // Indigo shadow
	colorWhite     = lipgloss.Color("#F8FAFC")
	colorDim       = lipgloss.Color("#64748B")
)

// Logo frames - Terminal Eye with different expressions
var eyeFrames = [][]string{
	// Frame 0: Eyes center
	{
		"╭─────────╮",
		"│  ◉   ◉  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// Frame 1: Eyes left
	{
		"╭─────────╮",
		"│ ◉   ◉   │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// Frame 2: Eyes right
	{
		"╭─────────╮",
		"│   ◉   ◉ │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// Frame 3: Blink
	{
		"╭─────────╮",
		"│  ─   ─  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// Frame 4: Wide eyes (excited)
	{
		"╭─────────╮",
		"│  ⊙   ⊙  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
}

// Main logo ASCII
var mainLogo = []string{
	"  ██████╗  ██████╗  ██████╗██╗     ██╗████████╗",
	" ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝",
	" ██║  ███╗██║   ██║██║     ██║     ██║   ██║   ",
	" ██║   ██║██║   ██║██║     ██║     ██║   ██║   ",
	" ╚██████╔╝╚██████╔╝╚██████╗███████╗██║   ██║   ",
	"  ╚═════╝  ╚═════╝  ╚═════╝╚══════╝╚═╝   ╚═╝   ",
}

// Glitch characters for matrix effect
var glitchChars = []rune{'░', '▒', '▓', '█', '▀', '▄', '▌', '▐', '■', '□', '▪', '▫'}

// Sparkle characters
var sparkleChars = []string{"✦", "✧", "⋆", "✨", "⭐", "✴", "✵", "❇", "★", "☆"}

// Model for the banner animation
type Model struct {
	width       int
	height      int
	frame       int
	startTime   time.Time
	phase       int // 0=glitch, 1=reveal, 2=eye-fly-in, 3=sparkle, 4=stable
	eyeX        float64
	eyeTargetX  float64
	eyeFrame    int
	glitchMap   [][]rune
	revealed    [][]bool
	sparkles    []Sparkle
	done        bool
}

type Sparkle struct {
	x, y    float64
	vx, vy  float64
	char    string
	life    int
	maxLife int
}

type tickMsg time.Time

func New(width, height int) Model {
	// Initialize glitch map
	glitchMap := make([][]rune, 12)
	for i := range glitchMap {
		glitchMap[i] = make([]rune, 60)
		for j := range glitchMap[i] {
			glitchMap[i][j] = glitchChars[rand.Intn(len(glitchChars))]
		}
	}

	// Initialize reveal mask
	revealed := make([][]bool, 12)
	for i := range revealed {
		revealed[i] = make([]bool, 60)
	}

	return Model{
		width:      width,
		height:     height,
		startTime:  time.Now(),
		eyeX:       80, // Start off-screen right
		eyeTargetX: 50,
		glitchMap:  glitchMap,
		revealed:   revealed,
		sparkles:   make([]Sparkle, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(frameTime, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Skip animation on any key
		m.done = true
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.frame++
		elapsed := time.Since(m.startTime)
		progress := float64(elapsed) / float64(totalDuration)

		// Phase transitions
		if progress < 0.15 {
			m.phase = 0 // Glitch
		} else if progress < 0.4 {
			m.phase = 1 // Reveal (scan line)
		} else if progress < 0.7 {
			m.phase = 2 // Eye fly-in
		} else if progress < 0.9 {
			m.phase = 3 // Sparkle
		} else {
			m.phase = 4 // Stable
		}

		// Update glitch
		if m.phase == 0 {
			for i := range m.glitchMap {
				for j := range m.glitchMap[i] {
					if rand.Float64() < 0.3 {
						m.glitchMap[i][j] = glitchChars[rand.Intn(len(glitchChars))]
					}
				}
			}
		}

		// Update reveal (scan line effect)
		if m.phase >= 1 {
			revealProgress := (progress - 0.15) / 0.25
			revealCol := int(revealProgress * 60)
			for i := range m.revealed {
				for j := 0; j < revealCol && j < 60; j++ {
					m.revealed[i][j] = true
				}
			}
		}

		// Update eye position (smooth fly-in)
		if m.phase >= 2 {
			// Ease-out animation
			m.eyeX += (m.eyeTargetX - m.eyeX) * 0.15
		}

		// Update eye frame (looking around)
		if m.frame%10 == 0 {
			if rand.Float64() < 0.3 {
				m.eyeFrame = rand.Intn(4) // Don't blink too often
			}
		}
		// Occasional blink
		if m.frame%45 == 0 {
			m.eyeFrame = 3 // Blink
		}

		// Spawn sparkles
		if m.phase == 3 && m.frame%3 == 0 {
			m.sparkles = append(m.sparkles, Sparkle{
				x:       m.eyeX + 5,
				y:       4,
				vx:      (rand.Float64() - 0.5) * 2,
				vy:      -rand.Float64() * 1.5,
				char:    sparkleChars[rand.Intn(len(sparkleChars))],
				life:    0,
				maxLife: 15 + rand.Intn(10),
			})
		}

		// Update sparkles
		newSparkles := make([]Sparkle, 0)
		for _, s := range m.sparkles {
			s.x += s.vx
			s.y += s.vy
			s.vy += 0.1 // Gravity
			s.life++
			if s.life < s.maxLife {
				newSparkles = append(newSparkles, s)
			}
		}
		m.sparkles = newSparkles

		// Check if done
		if progress >= 1.0 {
			m.done = true
		}

		return m, tick()
	}

	return m, nil
}

func (m Model) View() string {
	if m.done {
		return m.renderFinal()
	}

	var b strings.Builder

	// Calculate center position
	centerX := (m.width - 60) / 2
	if centerX < 0 {
		centerX = 0
	}
	padding := strings.Repeat(" ", centerX)

	// Render based on phase
	switch m.phase {
	case 0: // Glitch phase
		b.WriteString(m.renderGlitch(padding))
	case 1, 2, 3, 4: // Reveal and beyond
		b.WriteString(m.renderMain(padding))
	}

	return b.String()
}

func (m Model) renderGlitch(padding string) string {
	var b strings.Builder

	// Add some empty lines for vertical centering
	b.WriteString("\n\n")

	glitchStyle := lipgloss.NewStyle().Foreground(colorSecondary)
	dimStyle := lipgloss.NewStyle().Foreground(colorDim)

	for i, row := range m.glitchMap {
		b.WriteString(padding)
		for j, ch := range row {
			if rand.Float64() < 0.7 {
				b.WriteString(glitchStyle.Render(string(ch)))
			} else {
				b.WriteString(dimStyle.Render(string(ch)))
			}
			_ = j
		}
		b.WriteString("\n")
		_ = i
	}

	return b.String()
}

func (m Model) renderMain(padding string) string {
	var b strings.Builder

	b.WriteString("\n\n")

	// Title style with gradient effect
	titleStyle := lipgloss.NewStyle().Foreground(colorDim).Italic(true)
	b.WriteString(padding + titleStyle.Render("  Welcome to") + "\n\n")

	// Render main logo with color gradient
	for i, line := range mainLogo {
		b.WriteString(padding)
		for j, ch := range line {
			if !m.revealed[i][j] && m.phase < 4 {
				// Not yet revealed - show glitch or nothing
				if rand.Float64() < 0.3 {
					b.WriteString(lipgloss.NewStyle().Foreground(colorDim).Render(string(glitchChars[rand.Intn(len(glitchChars))])))
				} else {
					b.WriteString(" ")
				}
			} else {
				// Revealed - show with gradient color
				ratio := float64(j) / float64(len(line))
				var color lipgloss.Color
				if ratio < 0.33 {
					color = colorPrimary // Purple
				} else if ratio < 0.66 {
					color = colorSecondary // Cyan
				} else {
					color = colorAccent // Blue
				}
				style := lipgloss.NewStyle().Foreground(color).Bold(true)
				b.WriteString(style.Render(string(ch)))
			}
		}
		b.WriteString("\n")
	}

	// Tagline
	tagStyle := lipgloss.NewStyle().Foreground(colorHighlight).Italic(true)
	b.WriteString(padding + "                    " + tagStyle.Render("The Dream CLI") + "\n\n")

	// Render Terminal Eye at current position
	if m.phase >= 2 {
		eyeLines := eyeFrames[m.eyeFrame]
		eyeStyle := lipgloss.NewStyle().Foreground(colorGlow)
		eyeBorderStyle := lipgloss.NewStyle().Foreground(colorShadow)
		eyePupilStyle := lipgloss.NewStyle().Foreground(colorEye).Bold(true)

		for _, line := range eyeLines {
			spaces := int(m.eyeX)
			if spaces < 0 {
				spaces = 0
			}
			if spaces > m.width-15 {
				spaces = m.width - 15
			}
			b.WriteString(strings.Repeat(" ", spaces))

			// Color the eye parts
			for _, ch := range line {
				switch ch {
				case '◉', '⊙', '●':
					b.WriteString(eyePupilStyle.Render(string(ch)))
				case '╭', '╮', '╰', '╯', '│', '─':
					b.WriteString(eyeBorderStyle.Render(string(ch)))
				default:
					b.WriteString(eyeStyle.Render(string(ch)))
				}
			}
			b.WriteString("\n")
		}
	}

	// Render sparkles
	if m.phase >= 3 {
		sparkleStyle := lipgloss.NewStyle().Foreground(colorHighlight)
		for _, s := range m.sparkles {
			if s.y >= 0 && s.y < float64(m.height) && s.x >= 0 && s.x < float64(m.width) {
				// This is simplified - real implementation would need cursor positioning
				opacity := 1.0 - float64(s.life)/float64(s.maxLife)
				if opacity > 0.5 {
					b.WriteString(sparkleStyle.Render(s.char))
				}
			}
		}
	}

	return b.String()
}

func (m Model) renderFinal() string {
	var b strings.Builder

	centerX := (m.width - 60) / 2
	if centerX < 0 {
		centerX = 0
	}
	padding := strings.Repeat(" ", centerX)

	b.WriteString("\n\n")

	// Title
	titleStyle := lipgloss.NewStyle().Foreground(colorDim).Italic(true)
	b.WriteString(padding + titleStyle.Render("  Welcome to") + "\n\n")

	// Logo with full gradient
	for _, line := range mainLogo {
		b.WriteString(padding)
		for j, ch := range line {
			ratio := float64(j) / float64(len(line))
			// Smooth gradient: purple → cyan → blue
			var color lipgloss.Color
			if ratio < 0.33 {
				color = colorPrimary
			} else if ratio < 0.66 {
				color = colorSecondary
			} else {
				color = colorAccent
			}
			style := lipgloss.NewStyle().Foreground(color).Bold(true)
			b.WriteString(style.Render(string(ch)))
		}
		b.WriteString("\n")
	}

	// Tagline
	tagStyle := lipgloss.NewStyle().Foreground(colorHighlight).Italic(true)
	b.WriteString(padding + "                    " + tagStyle.Render("The Dream CLI") + "\n\n")

	// Final eye position
	eyeLines := eyeFrames[0]
	eyeStyle := lipgloss.NewStyle().Foreground(colorGlow)
	eyeBorderStyle := lipgloss.NewStyle().Foreground(colorShadow)
	eyePupilStyle := lipgloss.NewStyle().Foreground(colorEye).Bold(true)

	for _, line := range eyeLines {
		spaces := int(m.eyeTargetX)
		b.WriteString(strings.Repeat(" ", spaces))
		for _, ch := range line {
			switch ch {
			case '◉', '⊙', '●':
				b.WriteString(eyePupilStyle.Render(string(ch)))
			case '╭', '╮', '╰', '╯', '│', '─':
				b.WriteString(eyeBorderStyle.Render(string(ch)))
			default:
				b.WriteString(eyeStyle.Render(string(ch)))
			}
		}
		b.WriteString("\n")
	}

	// Version
	verStyle := lipgloss.NewStyle().Foreground(colorDim)
	b.WriteString("\n" + padding + verStyle.Render("  CLI Version 0.1.0") + "\n")

	return b.String()
}

func (m Model) Done() bool {
	return m.done
}

// Helper for wave effect
func wave(x, t float64) float64 {
	return math.Sin(x*0.5 + t*3)
}
