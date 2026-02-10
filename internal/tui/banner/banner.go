// Package banner provides animated ASCII banner for goclit
// GOAL: Better than GitHub Copilot CLI banner
// Features: 35+ frames, multicolor gradients, spring physics fly-in, particle effects
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
	fps           = 60 // Smooth 60fps
	frameTime     = time.Second / fps
	totalDuration = 4 * time.Second
)

// Extended color palette for multicolor gradients
var (
	// Primary gradient spectrum (10 colors for smooth transitions)
	spectrum = []lipgloss.Color{
		"#FF6B9D", // Pink
		"#C567FF", // Violet
		"#A855F7", // Purple
		"#8B5CF6", // Indigo
		"#6366F1", // Blue-indigo
		"#3B82F6", // Blue
		"#0EA5E9", // Sky
		"#06B6D4", // Cyan
		"#14B8A6", // Teal
		"#10B981", // Emerald
	}

	// Accent colors
	colorHighlight = lipgloss.Color("#F59E0B") // Amber
	colorGlow      = lipgloss.Color("#22D3EE") // Bright cyan
	colorEye       = lipgloss.Color("#10B981") // Emerald
	colorShadow    = lipgloss.Color("#6366F1") // Indigo
	colorWhite     = lipgloss.Color("#F8FAFC")
	colorDim       = lipgloss.Color("#64748B")
	colorGold      = lipgloss.Color("#FFD700") // Gold
	colorHot       = lipgloss.Color("#FF4500") // Hot
)

// Extended eye frames (10 frames for richer animation)
var eyeFrames = [][]string{
	// 0: Center
	{
		"╭─────────╮",
		"│  ◉   ◉  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// 1: Left
	{
		"╭─────────╮",
		"│ ◉   ◉   │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// 2: Right
	{
		"╭─────────╮",
		"│   ◉   ◉ │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// 3: Blink
	{
		"╭─────────╮",
		"│  ─   ─  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// 4: Wide (excited)
	{
		"╭─────────╮",
		"│  ⊙   ⊙  │",
		"│    ▼    │",
		"│  ╰◡─◡╯  │",
		"╰─────────╯",
	},
	// 5: Squint
	{
		"╭─────────╮",
		"│  ⁻   ⁻  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// 6: Surprised
	{
		"╭─────────╮",
		"│  ⊚   ⊚  │",
		"│    ▽    │",
		"│  ╰   ╯  │",
		"╰─────────╯",
	},
	// 7: Happy
	{
		"╭─────────╮",
		"│  ◠   ◠  │",
		"│    ▼    │",
		"│  ╰◡─◡╯  │",
		"╰─────────╯",
	},
	// 8: Looking up
	{
		"╭─────────╮",
		"│  ∘◉ ∘◉  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
	// 9: Looking down
	{
		"╭─────────╮",
		"│  ◉∘ ◉∘  │",
		"│    ▼    │",
		"│  ╰───╯  │",
		"╰─────────╯",
	},
}

// Main logo ASCII art
var mainLogo = []string{
	"  ██████╗  ██████╗  ██████╗██╗     ██╗████████╗",
	" ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝",
	" ██║  ███╗██║   ██║██║     ██║     ██║   ██║   ",
	" ██║   ██║██║   ██║██║     ██║     ██║   ██║   ",
	" ╚██████╔╝╚██████╔╝╚██████╗███████╗██║   ██║   ",
	"  ╚═════╝  ╚═════╝  ╚═════╝╚══════╝╚═╝   ╚═╝   ",
}

// Extended glitch characters
var glitchChars = []rune{
	'░', '▒', '▓', '█', '▀', '▄', '▌', '▐',
	'■', '□', '▪', '▫', '◢', '◣', '◤', '◥',
	'╔', '╗', '╚', '╝', '║', '═', '╬', '╫',
	'⣿', '⣷', '⣯', '⣟', '⡿', '⢿', '⣻', '⣽',
}

// Sparkle/particle characters
var sparkleChars = []string{
	"✦", "✧", "⋆", "✨", "⭐", "✴", "✵", "❇",
	"★", "☆", "◆", "◇", "❖", "✺", "✹", "✸",
}

// Spring physics for smooth fly-in
type Spring struct {
	position float64
	velocity float64
	target   float64
	omega    float64 // Angular frequency
	zeta     float64 // Damping ratio
}

func (s *Spring) Update(dt float64) {
	// Critically damped spring: smooth, no overshoot
	diff := s.target - s.position
	s.velocity += (diff*s.omega*s.omega - 2*s.zeta*s.omega*s.velocity) * dt
	s.position += s.velocity * dt
}

// Particle for sparkle effects
type Particle struct {
	x, y    float64
	vx, vy  float64
	char    string
	color   lipgloss.Color
	life    int
	maxLife int
	scale   float64
}

// Model for the banner animation
type Model struct {
	width     int
	height    int
	frame     int
	startTime time.Time
	phase     int // 0=matrix, 1=scan-reveal, 2=fly-in, 3=sparkle, 4=rainbow-wave, 5=stable

	// Eye animation
	eyeSpring  Spring
	eyeFrame   int
	eyeHistory []int // For smooth transitions

	// Logo animation
	revealProgress float64
	glitchMap      [][]rune
	waveOffset     float64

	// Particles
	particles []Particle

	// State
	done bool
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

	return Model{
		width:     width,
		height:    height,
		startTime: time.Now(),
		eyeSpring: Spring{
			position: 100, // Start far right (off-screen)
			velocity: 0,
			target:   50, // Target center-ish
			omega:    8,  // Fast response
			zeta:     0.6, // Slight overshoot for bounce
		},
		glitchMap:  glitchMap,
		particles:  make([]Particle, 0),
		eyeHistory: make([]int, 0),
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
		m.done = true
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.frame++
		elapsed := time.Since(m.startTime)
		progress := float64(elapsed) / float64(totalDuration)
		dt := 1.0 / float64(fps)

		// Phase transitions (5 phases)
		switch {
		case progress < 0.12:
			m.phase = 0 // Matrix glitch
		case progress < 0.35:
			m.phase = 1 // Scan-line reveal
		case progress < 0.55:
			m.phase = 2 // Eye fly-in
		case progress < 0.75:
			m.phase = 3 // Sparkle burst
		case progress < 0.92:
			m.phase = 4 // Rainbow wave
		default:
			m.phase = 5 // Stable
		}

		// Update glitch map
		if m.phase == 0 {
			intensity := 1.0 - progress/0.12 // Fade out
			for i := range m.glitchMap {
				for j := range m.glitchMap[i] {
					if rand.Float64() < 0.4*intensity {
						m.glitchMap[i][j] = glitchChars[rand.Intn(len(glitchChars))]
					}
				}
			}
		}

		// Update reveal progress
		if m.phase >= 1 {
			revealStart := 0.12
			revealEnd := 0.35
			m.revealProgress = clamp((progress-revealStart)/(revealEnd-revealStart), 0, 1)
			// Ease-out cubic for smooth reveal
			m.revealProgress = 1 - math.Pow(1-m.revealProgress, 3)
		}

		// Update eye spring (fly-in)
		if m.phase >= 2 {
			m.eyeSpring.Update(dt)
		}

		// Update eye frame animation
		if m.frame%6 == 0 {
			// Mostly center, occasional look around
			choices := []int{0, 0, 0, 0, 1, 2, 8, 9}
			m.eyeFrame = choices[rand.Intn(len(choices))]
		}
		// Occasional blink
		if m.frame%55 == 0 {
			m.eyeFrame = 3
		}
		// Excited during sparkle phase
		if m.phase == 3 && m.frame%20 == 0 {
			m.eyeFrame = 4
		}
		// Happy at the end
		if m.phase == 4 && m.frame%30 == 0 {
			m.eyeFrame = 7
		}

		// Spawn particles during sparkle phase
		if m.phase == 3 && m.frame%2 == 0 {
			for i := 0; i < 3; i++ { // Multiple particles per frame
				m.particles = append(m.particles, Particle{
					x:       m.eyeSpring.position + 5 + rand.Float64()*5,
					y:       10 + rand.Float64()*3,
					vx:      (rand.Float64() - 0.5) * 4,
					vy:      -rand.Float64()*3 - 1,
					char:    sparkleChars[rand.Intn(len(sparkleChars))],
					color:   spectrum[rand.Intn(len(spectrum))],
					life:    0,
					maxLife: 20 + rand.Intn(15),
					scale:   0.5 + rand.Float64()*0.5,
				})
			}
		}

		// Update particles
		alive := make([]Particle, 0)
		for _, p := range m.particles {
			p.x += p.vx * dt * 20
			p.y += p.vy * dt * 20
			p.vy += 0.15 // Gravity
			p.life++
			if p.life < p.maxLife && p.y < float64(m.height) && p.y > 0 {
				alive = append(alive, p)
			}
		}
		m.particles = alive

		// Update wave offset for rainbow effect
		if m.phase >= 4 {
			m.waveOffset += dt * 3
		}

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

	centerX := (m.width - 50) / 2
	if centerX < 0 {
		centerX = 0
	}
	padding := strings.Repeat(" ", centerX)

	b.WriteString("\n\n")

	switch m.phase {
	case 0:
		b.WriteString(m.renderGlitch(padding))
	default:
		b.WriteString(m.renderMain(padding))
	}

	return b.String()
}

func (m Model) renderGlitch(padding string) string {
	var b strings.Builder

	for _, row := range m.glitchMap {
		b.WriteString(padding)
		for j, ch := range row {
			// Use spectrum colors for matrix effect
			colorIdx := (j + m.frame/2) % len(spectrum)
			if rand.Float64() < 0.6 {
				style := lipgloss.NewStyle().Foreground(spectrum[colorIdx])
				b.WriteString(style.Render(string(ch)))
			} else {
				style := lipgloss.NewStyle().Foreground(colorDim)
				b.WriteString(style.Render(string(ch)))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderMain(padding string) string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().Foreground(colorDim).Italic(true)
	b.WriteString(padding + titleStyle.Render("  Welcome to") + "\n\n")

	// Main logo with animated gradient
	revealCol := int(m.revealProgress * float64(len(mainLogo[0])))

	for i, line := range mainLogo {
		b.WriteString(padding)
		runes := []rune(line)
		for j, ch := range runes {
			if j > revealCol && m.phase < 5 {
				// Not yet revealed
				if rand.Float64() < 0.2 {
					style := lipgloss.NewStyle().Foreground(colorDim)
					b.WriteString(style.Render(string(glitchChars[rand.Intn(len(glitchChars))])))
				} else {
					b.WriteString(" ")
				}
			} else {
				// Revealed with animated rainbow gradient
				var colorIdx int
				if m.phase >= 4 {
					// Rainbow wave effect
					wave := math.Sin(float64(j)*0.15 + float64(i)*0.3 + m.waveOffset)
					colorIdx = int((wave+1)/2*float64(len(spectrum)-1)) % len(spectrum)
				} else {
					// Static gradient
					ratio := float64(j) / float64(len(runes))
					colorIdx = int(ratio * float64(len(spectrum)-1))
				}
				if colorIdx < 0 {
					colorIdx = 0
				}
				if colorIdx >= len(spectrum) {
					colorIdx = len(spectrum) - 1
				}
				style := lipgloss.NewStyle().Foreground(spectrum[colorIdx]).Bold(true)
				b.WriteString(style.Render(string(ch)))
			}
		}
		b.WriteString("\n")
	}

	// Tagline with glow
	if m.revealProgress > 0.6 {
		tagStyle := lipgloss.NewStyle().Foreground(colorHighlight).Italic(true)
		b.WriteString(padding + "               " + tagStyle.Render("✨ The Dream CLI ✨") + "\n\n")
	} else {
		b.WriteString("\n\n")
	}

	// Terminal Eye with spring physics
	if m.phase >= 2 {
		eye := eyeFrames[m.eyeFrame]
		eyeStyle := lipgloss.NewStyle().Foreground(colorGlow)
		borderStyle := lipgloss.NewStyle().Foreground(colorShadow)
		pupilStyle := lipgloss.NewStyle().Foreground(colorEye).Bold(true)

		for _, line := range eye {
			spaces := int(m.eyeSpring.position)
			if spaces < 0 {
				spaces = 0
			}
			if spaces > m.width-15 {
				spaces = m.width - 15
			}
			b.WriteString(strings.Repeat(" ", spaces))

			for _, ch := range line {
				switch ch {
				case '◉', '⊙', '●', '∘', '◠', '⊚', '⁻':
					b.WriteString(pupilStyle.Render(string(ch)))
				case '╭', '╮', '╰', '╯', '│', '─':
					b.WriteString(borderStyle.Render(string(ch)))
				default:
					b.WriteString(eyeStyle.Render(string(ch)))
				}
			}
			b.WriteString("\n")
		}
	}

	// Render particles (simplified display at bottom)
	if len(m.particles) > 0 {
		b.WriteString("\n" + padding)
		displayed := 0
		for _, p := range m.particles {
			if p.life < p.maxLife/2 && displayed < 20 {
				style := lipgloss.NewStyle().Foreground(p.color)
				b.WriteString(style.Render(p.char + " "))
				displayed++
			}
		}
		b.WriteString("\n")
	}

	// Version
	if m.phase >= 5 {
		verStyle := lipgloss.NewStyle().Foreground(colorDim)
		b.WriteString("\n" + padding + verStyle.Render("  CLI v0.2.0") + "\n")
	}

	return b.String()
}

func (m Model) renderFinal() string {
	var b strings.Builder

	centerX := (m.width - 50) / 2
	if centerX < 0 {
		centerX = 0
	}
	padding := strings.Repeat(" ", centerX)

	b.WriteString("\n\n")

	// Title
	titleStyle := lipgloss.NewStyle().Foreground(colorDim).Italic(true)
	b.WriteString(padding + titleStyle.Render("  Welcome to") + "\n\n")

	// Final logo with stable rainbow gradient
	for i, line := range mainLogo {
		b.WriteString(padding)
		runes := []rune(line)
		for j, ch := range runes {
			// Subtle wave
			wave := math.Sin(float64(j)*0.12 + float64(i)*0.25)
			colorIdx := int((wave+1)/2*float64(len(spectrum)-1)) % len(spectrum)
			style := lipgloss.NewStyle().Foreground(spectrum[colorIdx]).Bold(true)
			b.WriteString(style.Render(string(ch)))
		}
		b.WriteString("\n")
	}

	// Tagline
	tagStyle := lipgloss.NewStyle().Foreground(colorGold).Italic(true).Bold(true)
	b.WriteString(padding + "               " + tagStyle.Render("✨ The Dream CLI ✨") + "\n\n")

	// Final eye
	eye := eyeFrames[7] // Happy face
	eyeStyle := lipgloss.NewStyle().Foreground(colorGlow)
	borderStyle := lipgloss.NewStyle().Foreground(colorShadow)
	pupilStyle := lipgloss.NewStyle().Foreground(colorEye).Bold(true)

	for _, line := range eye {
		spaces := int(m.eyeSpring.target)
		b.WriteString(strings.Repeat(" ", spaces))
		for _, ch := range line {
			switch ch {
			case '◉', '⊙', '●', '◠', '◡':
				b.WriteString(pupilStyle.Render(string(ch)))
			case '╭', '╮', '╰', '╯', '│', '─':
				b.WriteString(borderStyle.Render(string(ch)))
			default:
				b.WriteString(eyeStyle.Render(string(ch)))
			}
		}
		b.WriteString("\n")
	}

	// Version
	verStyle := lipgloss.NewStyle().Foreground(colorDim)
	b.WriteString("\n" + padding + verStyle.Render("  CLI v0.2.0") + "\n")

	return b.String()
}

func (m Model) Done() bool {
	return m.done
}

// Helper
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
