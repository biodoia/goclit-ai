// Intro animation - BETTER than Copilot CLI
// Features: Glitch reveal, scan line, eye fly-in, sparkles, gradient
package tui

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps       = 30
	frameTime = time.Second / fps
)

// Animation timeline (milliseconds) - precise timing like Copilot
const (
	// Phase 0: Glitch matrix
	tGlitchStart = 0
	tGlitchEnd   = 400

	// Phase 1: Scan line reveal logo
	tRevealStart = 400
	tRevealEnd   = 1200

	// Phase 2: Eye fly-in
	tEyeStart = 800
	tEyeEnd   = 1800

	// Phase 3: Sparkle celebration
	tSparkleStart = 1600
	tSparkleEnd   = 2400

	// Phase 4: Stable (hold)
	tStableStart = 2400
	tTotalTime   = 3000 // 3 seconds total
)

// Color palette - semantic roles
var (
	colPrimary   = lipgloss.Color("#A855F7") // Purple
	colSecondary = lipgloss.Color("#06B6D4") // Cyan
	colAccent    = lipgloss.Color("#3B82F6") // Blue
	colHighlight = lipgloss.Color("#F59E0B") // Amber
	colGlow      = lipgloss.Color("#22D3EE") // Bright cyan
	colEye       = lipgloss.Color("#10B981") // Emerald
	colShadow    = lipgloss.Color("#6366F1") // Indigo
	colWhite     = lipgloss.Color("#F8FAFC")
	colDim       = lipgloss.Color("#64748B")
)

// Main GOCLIT logo
var mainLogo = []string{
	"  ██████╗  ██████╗  ██████╗██╗     ██╗████████╗",
	" ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝",
	" ██║  ███╗██║   ██║██║     ██║     ██║   ██║   ",
	" ██║   ██║██║   ██║██║     ██║     ██║   ██║   ",
	" ╚██████╔╝╚██████╔╝╚██████╗███████╗██║   ██║   ",
	"  ╚═════╝  ╚═════╝  ╚═════╝╚══════╝╚═╝   ╚═╝   ",
}

// Terminal Eye frames
var eyeFrames = [][]string{
	{ // Center
		"╭─────────╮",
		"│  ◉   ◉  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	{ // Left
		"╭─────────╮",
		"│ ◉   ◉   │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	{ // Right
		"╭─────────╮",
		"│   ◉   ◉ │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	{ // Blink
		"╭─────────╮",
		"│  ─   ─  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	{ // Wide (excited)
		"╭─────────╮",
		"│  ⊙   ⊙  │",
		"│    ▼    │",
		"│  ╰◡─◡╯  │",
		"╰─────────╯",
	},
}

// Glitch characters
var glitchChars = []rune{'░', '▒', '▓', '█', '▀', '▄', '▌', '▐', '■', '□', '▪', '▫', '◢', '◣', '◤', '◥'}

// Sparkle characters
var sparkles = []string{"✦", "✧", "⋆", "✨", "⭐", "✴", "★", "☆", "·", "°"}

const tagline = "The Dream CLI"
const version = "v0.2.0"

// Sparkle particle
type Sparkle struct {
	x, y   float64
	vx, vy float64
	char   string
	life   int
}

// IntroModel handles the intro animation
type IntroModel struct {
	width      int
	height     int
	startTime  time.Time
	frame      int
	phase      int // 0=glitch, 1=scanReveal, 2=eyeFlyIn, 3=sparkle, 4=done
	eyeX       float64
	eyeFrame   int
	revealCol  int
	glitchMap  [][]rune
	particles  []Sparkle
	done       bool
	showPanes  bool
}

type tickMsg time.Time

func NewIntro(width, height int) IntroModel {
	// Initialize glitch
	glitch := make([][]rune, 8)
	for i := range glitch {
		glitch[i] = make([]rune, 50)
		for j := range glitch[i] {
			glitch[i][j] = glitchChars[rand.Intn(len(glitchChars))]
		}
	}

	return IntroModel{
		width:     width,
		height:    height,
		startTime: time.Now(),
		eyeX:      70, // Start off-screen
		glitchMap: glitch,
		particles: make([]Sparkle, 0),
	}
}

func (m IntroModel) Init() tea.Cmd {
	return tea.Batch(tick(), tea.EnterAltScreen)
}

func tick() tea.Cmd {
	return tea.Tick(frameTime, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m IntroModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Skip on any key
		m.done = true
		m.showPanes = true
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.frame++
		elapsed := time.Since(m.startTime)
		ms := elapsed.Milliseconds() // Precise milliseconds

		// Phase transitions based on exact timing
		switch {
		case ms < tGlitchEnd:
			m.phase = 0 // Glitch
		case ms < tRevealEnd:
			m.phase = 1 // Scan reveal
		case ms < tEyeEnd:
			m.phase = 2 // Eye fly-in
		case ms < tSparkleEnd:
			m.phase = 3 // Sparkle celebration
		default:
			m.phase = 4 // Stable
		}

		// Update glitch characters
		if m.phase == 0 {
			for i := range m.glitchMap {
				for j := range m.glitchMap[i] {
					if rand.Float64() < 0.4 {
						m.glitchMap[i][j] = glitchChars[rand.Intn(len(glitchChars))]
					}
				}
			}
		}

		// Update scan reveal (linear from tRevealStart to tRevealEnd)
		if ms >= tRevealStart {
			revealProgress := float64(ms-tRevealStart) / float64(tRevealEnd-tRevealStart)
			if revealProgress > 1 {
				revealProgress = 1
			}
			m.revealCol = int(revealProgress * 50)
		}

		// Update eye fly-in (ease-out cubic)
		if ms >= tEyeStart {
			eyeProgress := float64(ms-tEyeStart) / float64(tEyeEnd-tEyeStart)
			if eyeProgress > 1 {
				eyeProgress = 1
			}
			// Ease-out cubic: 1 - (1-t)^3
			eased := 1 - (1-eyeProgress)*(1-eyeProgress)*(1-eyeProgress)
			targetX := 52.0
			startX := 70.0
			m.eyeX = startX + (targetX-startX)*eased
		}

		// Update eye animation
		if m.frame%8 == 0 {
			choices := []int{0, 0, 0, 1, 2} // Mostly center
			m.eyeFrame = choices[rand.Intn(len(choices))]
		}
		if m.frame%50 == 0 {
			m.eyeFrame = 3 // Blink
		}
		if m.phase == 3 && m.frame%30 == 0 {
			m.eyeFrame = 4 // Excited!
		}

		// Spawn sparkles
		if m.phase == 3 && m.frame%2 == 0 {
			m.particles = append(m.particles, Sparkle{
				x:    m.eyeX + 5 + rand.Float64()*3,
				y:    8 + rand.Float64()*2,
				vx:   (rand.Float64() - 0.5) * 3,
				vy:   -rand.Float64()*2 - 0.5,
				char: sparkles[rand.Intn(len(sparkles))],
				life: 0,
			})
		}

		// Update sparkles
		alive := make([]Sparkle, 0)
		for _, p := range m.particles {
			p.x += p.vx
			p.y += p.vy
			p.vy += 0.15 // Gravity
			p.life++
			if p.life < 20 && p.y < float64(m.height) && p.y > 0 {
				alive = append(alive, p)
			}
		}
		m.particles = alive

		// Check done (after total time)
		if ms >= tTotalTime {
			m.done = true
			m.showPanes = true
		}

		return m, tick()
	}

	return m, nil
}

func (m IntroModel) View() string {
	var b strings.Builder

	// Center padding
	padX := (m.width - 50) / 2
	if padX < 0 {
		padX = 0
	}
	pad := strings.Repeat(" ", padX)

	b.WriteString("\n\n")

	switch m.phase {
	case 0:
		// GLITCH PHASE - Matrix-style reveal
		b.WriteString(m.renderGlitch(pad))

	default:
		// MAIN RENDER
		b.WriteString(m.renderMain(pad))
	}

	return b.String()
}

func (m IntroModel) renderGlitch(pad string) string {
	var b strings.Builder

	style1 := lipgloss.NewStyle().Foreground(colSecondary)
	style2 := lipgloss.NewStyle().Foreground(colPrimary)
	dimStyle := lipgloss.NewStyle().Foreground(colDim)

	for _, row := range m.glitchMap {
		b.WriteString(pad)
		for _, ch := range row {
			r := rand.Float64()
			if r < 0.5 {
				b.WriteString(style1.Render(string(ch)))
			} else if r < 0.8 {
				b.WriteString(style2.Render(string(ch)))
			} else {
				b.WriteString(dimStyle.Render(string(ch)))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m IntroModel) renderMain(pad string) string {
	var b strings.Builder

	// "Welcome to" title
	titleStyle := lipgloss.NewStyle().Foreground(colDim).Italic(true)
	b.WriteString(pad + titleStyle.Render("  Welcome to") + "\n\n")

	// GOCLIT logo with gradient + reveal effect
	for i, line := range mainLogo {
		b.WriteString(pad)
		for j, ch := range line {
			if j > m.revealCol && m.phase < 4 {
				// Not revealed yet
				if rand.Float64() < 0.2 {
					b.WriteString(lipgloss.NewStyle().Foreground(colDim).Render(string(glitchChars[rand.Intn(len(glitchChars))])))
				} else {
					b.WriteString(" ")
				}
			} else {
				// Revealed with gradient
				ratio := float64(j) / float64(len(line))
				var col lipgloss.Color
				switch {
				case ratio < 0.33:
					col = colPrimary
				case ratio < 0.66:
					col = colSecondary
				default:
					col = colAccent
				}
				b.WriteString(lipgloss.NewStyle().Foreground(col).Bold(true).Render(string(ch)))
			}
		}
		b.WriteString("\n")
		_ = i
	}

	// Tagline
	tagStyle := lipgloss.NewStyle().Foreground(colHighlight).Italic(true)
	if m.revealCol > 30 || m.phase >= 4 {
		b.WriteString(pad + "                " + tagStyle.Render(tagline) + "\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Terminal Eye
	if m.phase >= 2 {
		eye := eyeFrames[m.eyeFrame]
		eyeStyle := lipgloss.NewStyle().Foreground(colGlow)
		borderStyle := lipgloss.NewStyle().Foreground(colShadow)
		pupilStyle := lipgloss.NewStyle().Foreground(colEye).Bold(true)

		for _, line := range eye {
			spaces := int(m.eyeX)
			if spaces < 0 {
				spaces = 0
			}
			if spaces > m.width-15 {
				spaces = m.width - 15
			}
			b.WriteString(strings.Repeat(" ", spaces))

			for _, ch := range line {
				switch ch {
				case '◉', '⊙', '─':
					b.WriteString(pupilStyle.Render(string(ch)))
				case '╭', '╮', '╰', '╯', '│':
					b.WriteString(borderStyle.Render(string(ch)))
				default:
					b.WriteString(eyeStyle.Render(string(ch)))
				}
			}
			b.WriteString("\n")
		}
	}

	// Sparkles overlay (simplified - just add at bottom)
	if m.phase >= 3 && len(m.particles) > 0 {
		sparkleStyle := lipgloss.NewStyle().Foreground(colHighlight)
		b.WriteString("\n" + pad)
		for _, p := range m.particles {
			if p.life < 15 {
				b.WriteString(sparkleStyle.Render(p.char + " "))
			}
		}
		b.WriteString("\n")
	}

	// Version (only when stable)
	if m.phase >= 4 {
		verStyle := lipgloss.NewStyle().Foreground(colDim)
		b.WriteString("\n" + pad + verStyle.Render("  CLI "+version) + "\n")
	}

	return b.String()
}

func (m IntroModel) Done() bool {
	return m.done
}

func (m IntroModel) ShowPanes() bool {
	return m.showPanes
}

// Helper functions
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Gradient is defined in theme.go
