// Intro Enhanced - COPILOT-QUALITY multicolor animation
// Features: True rainbow gradient, character fly-in, wave shimmer, particle burst
package tui

import (
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Extended color palette - True rainbow gradient
var rainbowColors = []lipgloss.Color{
	"#FF0080", // Magenta
	"#FF0040", // Red-Magenta
	"#FF0000", // Red
	"#FF4000", // Orange-Red
	"#FF8000", // Orange
	"#FFC000", // Gold
	"#FFFF00", // Yellow
	"#80FF00", // Yellow-Green
	"#00FF00", // Green
	"#00FF80", // Cyan-Green
	"#00FFFF", // Cyan
	"#0080FF", // Sky Blue
	"#0000FF", // Blue
	"#4000FF", // Blue-Purple
	"#8000FF", // Purple
	"#C000FF", // Violet
	"#FF00FF", // Magenta (loop)
}

// Neon cyberpunk palette
var neonPalette = []lipgloss.Color{
	"#FF1493", // Deep Pink
	"#00CED1", // Dark Turquoise
	"#9400D3", // Dark Violet
	"#00FF7F", // Spring Green
	"#FF4500", // Orange Red
	"#1E90FF", // Dodger Blue
	"#FFD700", // Gold
	"#FF69B4", // Hot Pink
}

// Character reveal data for fly-in effect
type CharReveal struct {
	char     rune
	x, y     int      // final position
	currX    float64  // current animated position
	currY    float64
	delay    float64  // start delay (0-1)
	revealed bool
	color    lipgloss.Color
}

// Enhanced particle with more properties
type EnhancedParticle struct {
	x, y     float64
	vx, vy   float64
	ax, ay   float64 // acceleration
	char     string
	color    lipgloss.Color
	life     int
	maxLife  int
	size     float64 // for scaling effects
	rotation float64
}

// EnhancedIntroModel - Copilot-quality animation
type EnhancedIntroModel struct {
	width, height int
	startTime     time.Time
	frame         int

	// Character fly-in state
	chars      []CharReveal
	charsReady bool

	// Wave shimmer state
	wavePhase float64

	// Particles
	particles []EnhancedParticle

	// Animation phases
	phase int // 0=dark, 1=matrix, 2=char-fly-in, 3=shimmer, 4=stable
	done  bool

	// Eye animation
	eyeFrame   int
	eyeX, eyeY float64
	eyeScale   float64

	// Glitch map
	glitchMap [][]rune

	// Rainbow offset (for animated gradient)
	rainbowOffset float64
}

// Phase timing (milliseconds)
const (
	pDarkEnd     = 300
	pMatrixEnd   = 800
	pFlyInEnd    = 2000
	pShimmerEnd  = 2800
	pTotalEnd    = 3500
)

// Enhanced sparkle characters
var enhancedSparkles = []string{
	"✦", "✧", "⋆", "✨", "⭐", "✴", "★", "☆", "✶", "✷", "✸", "✹", 
	"·", "°", "•", "◦", "○", "◌", "◎", "●", "◐", "◑", "◒", "◓",
	"⁕", "⁎", "⁑", "⁂", "※", "❋", "❊", "❉", "❈", "❇", "❆", "❅",
}

// Flying characters for matrix effect
var flyingChars = []rune{
	'ア', 'イ', 'ウ', 'エ', 'オ', 'カ', 'キ', 'ク', 'ケ', 'コ', // Katakana
	'░', '▒', '▓', '█', '▀', '▄', '▌', '▐', '■', '□', '▪', '▫',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'G', 'O', 'C', 'L', 'I', 'T', 'A', 'Z',
}

func NewEnhancedIntro(width, height int) *EnhancedIntroModel {
	m := &EnhancedIntroModel{
		width:     width,
		height:    height,
		startTime: time.Now(),
		particles: make([]EnhancedParticle, 0),
		eyeX:      float64(width), // start off-screen
		eyeY:      10,
		eyeScale:  0.0,
	}

	// Initialize glitch map
	m.glitchMap = make([][]rune, 8)
	for i := range m.glitchMap {
		m.glitchMap[i] = make([]rune, 50)
		for j := range m.glitchMap[i] {
			m.glitchMap[i][j] = flyingChars[rand.Intn(len(flyingChars))]
		}
	}

	// Initialize character fly-in for logo
	m.initCharFlyIn()

	return m
}

func (m *EnhancedIntroModel) initCharFlyIn() {
	// Calculate delays based on position for wave effect
	for y, line := range mainLogo {
		for x, ch := range line {
			if ch != ' ' {
				// Delay based on distance from center-left (wave sweeps right)
				delay := float64(x) / float64(len(line)) * 0.6
				delay += float64(y) / float64(len(mainLogo)) * 0.2
				delay += rand.Float64() * 0.1 // slight randomization

				// Random starting position (off-screen)
				startX := float64(x) + (rand.Float64()-0.5)*100
				startY := float64(y) + (rand.Float64()-0.5)*50 - 20 // mostly from above

				// Color based on position
				colorIdx := int(float64(x) / float64(len(line)) * float64(len(rainbowColors)))
				if colorIdx >= len(rainbowColors) {
					colorIdx = len(rainbowColors) - 1
				}

				m.chars = append(m.chars, CharReveal{
					char:     ch,
					x:        x,
					y:        y,
					currX:    startX,
					currY:    startY,
					delay:    delay,
					revealed: false,
					color:    rainbowColors[colorIdx],
				})
			}
		}
	}
	m.charsReady = true
}

func (m *EnhancedIntroModel) Init() tea.Cmd {
	return tea.Batch(tickEnhanced(), tea.EnterAltScreen)
}

func tickEnhanced() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { // 60 FPS
		return tickMsg(t)
	})
}

func (m *EnhancedIntroModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		ms := elapsed.Milliseconds()

		// Update phase
		switch {
		case ms < pDarkEnd:
			m.phase = 0
		case ms < pMatrixEnd:
			m.phase = 1
		case ms < pFlyInEnd:
			m.phase = 2
		case ms < pShimmerEnd:
			m.phase = 3
		default:
			m.phase = 4
		}

		// Update rainbow offset for animated gradient
		m.rainbowOffset += 0.02
		if m.rainbowOffset >= 1.0 {
			m.rainbowOffset = 0
		}

		// Update wave phase
		m.wavePhase += 0.15

		// Update glitch map
		if m.phase <= 1 {
			m.updateGlitch()
		}

		// Update character fly-in
		if m.phase >= 2 {
			m.updateCharFlyIn(float64(ms-pMatrixEnd) / float64(pFlyInEnd-pMatrixEnd))
		}

		// Update eye animation
		if m.phase >= 2 {
			m.updateEye(float64(ms-pMatrixEnd) / float64(pFlyInEnd-pMatrixEnd))
		}

		// Spawn particles during shimmer phase
		if m.phase == 3 && m.frame%2 == 0 {
			m.spawnParticle()
		}

		// Update particles
		m.updateParticles()

		// Update eye frame
		if m.frame%12 == 0 {
			m.eyeFrame = rand.Intn(3) // 0, 1, 2 (center, left, right)
		}
		if m.frame%60 == 0 {
			m.eyeFrame = 3 // blink
		}
		if m.phase == 3 {
			m.eyeFrame = 4 // excited
		}

		// Check done
		if ms >= pTotalEnd {
			m.done = true
		}

		return m, tickEnhanced()
	}

	return m, nil
}

func (m *EnhancedIntroModel) updateGlitch() {
	intensity := 0.4
	if m.phase == 0 {
		intensity = 0.1 // slower at start
	}

	for i := range m.glitchMap {
		for j := range m.glitchMap[i] {
			if rand.Float64() < intensity {
				m.glitchMap[i][j] = flyingChars[rand.Intn(len(flyingChars))]
			}
		}
	}
}

func (m *EnhancedIntroModel) updateCharFlyIn(progress float64) {
	if progress > 1 {
		progress = 1
	}

	for i := range m.chars {
		ch := &m.chars[i]

		// Check if this character should start animating
		effectiveProgress := progress - ch.delay
		if effectiveProgress < 0 {
			continue
		}
		if effectiveProgress > 1 {
			effectiveProgress = 1
		}

		// Ease-out cubic for smooth landing
		eased := 1 - math.Pow(1-effectiveProgress, 3)

		// Interpolate position
		ch.currX = ch.currX + (float64(ch.x)-ch.currX)*eased
		ch.currY = ch.currY + (float64(ch.y)-ch.currY)*eased

		if effectiveProgress >= 0.8 {
			ch.revealed = true
		}
	}
}

func (m *EnhancedIntroModel) updateEye(progress float64) {
	if progress > 1 {
		progress = 1
	}

	// Ease-out for smooth fly-in
	eased := 1 - math.Pow(1-progress, 3)

	// Target position (right of logo)
	targetX := float64(m.width/2 + 25)
	startX := float64(m.width + 20)

	m.eyeX = startX + (targetX-startX)*eased
	m.eyeScale = eased

	// Add slight bounce at end
	if progress > 0.9 {
		bounce := math.Sin((progress - 0.9) * 10 * math.Pi) * 0.05
		m.eyeX += bounce * 5
	}
}

func (m *EnhancedIntroModel) spawnParticle() {
	m.particles = append(m.particles, EnhancedParticle{
		x:       m.eyeX + 5,
		y:       m.eyeY + 2,
		vx:      (rand.Float64() - 0.5) * 4,
		vy:      -rand.Float64()*3 - 1,
		ax:      (rand.Float64() - 0.5) * 0.1,
		ay:      0.15, // gravity
		char:    enhancedSparkles[rand.Intn(len(enhancedSparkles))],
		color:   neonPalette[rand.Intn(len(neonPalette))],
		life:    0,
		maxLife: 20 + rand.Intn(15),
		size:    0.5 + rand.Float64()*0.5,
	})
}

func (m *EnhancedIntroModel) updateParticles() {
	alive := make([]EnhancedParticle, 0)
	for _, p := range m.particles {
		p.vx += p.ax
		p.vy += p.ay
		p.x += p.vx
		p.y += p.vy
		p.life++
		p.rotation += 0.1

		if p.life < p.maxLife && p.y < float64(m.height) && p.y > 0 {
			alive = append(alive, p)
		}
	}
	m.particles = alive
}

func (m *EnhancedIntroModel) View() string {
	var b strings.Builder

	padX := (m.width - 50) / 2
	if padX < 0 {
		padX = 0
	}
	pad := strings.Repeat(" ", padX)

	b.WriteString("\n\n")

	switch m.phase {
	case 0:
		// Dark phase - just a hint of matrix
		b.WriteString(m.renderDark(pad))
	case 1:
		// Matrix phase - full glitch
		b.WriteString(m.renderMatrix(pad))
	default:
		// Main animation with character fly-in
		b.WriteString(m.renderFlyIn(pad))
	}

	return b.String()
}

func (m *EnhancedIntroModel) renderDark(pad string) string {
	var b strings.Builder

	// Mostly dark with occasional flickers
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#1a1a2e"))

	for i := 0; i < 8; i++ {
		b.WriteString(pad)
		for j := 0; j < 50; j++ {
			if rand.Float64() < 0.05 {
				// Rare flicker
				color := neonPalette[rand.Intn(len(neonPalette))]
				b.WriteString(lipgloss.NewStyle().Foreground(color).Render(string(flyingChars[rand.Intn(len(flyingChars))])))
			} else {
				b.WriteString(dimStyle.Render(" "))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m *EnhancedIntroModel) renderMatrix(pad string) string {
	var b strings.Builder

	for i, row := range m.glitchMap {
		b.WriteString(pad)
		for j, ch := range row {
			// Color based on position + time for wave effect
			offset := m.rainbowOffset + float64(j)/50 + float64(i)/20
			colorIdx := int(offset * float64(len(rainbowColors))) % len(rainbowColors)
			color := rainbowColors[colorIdx]

			// Vary intensity
			if rand.Float64() < 0.3 {
				color = lipgloss.Color("#334155") // dim
			}

			b.WriteString(lipgloss.NewStyle().Foreground(color).Render(string(ch)))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m *EnhancedIntroModel) renderFlyIn(pad string) string {
	var b strings.Builder

	// Build 2D grid for logo
	logoHeight := len(mainLogo)
	logoWidth := 0
	for _, line := range mainLogo {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}

	// Create empty grid
	grid := make([][]struct {
		ch    rune
		color lipgloss.Color
	}, logoHeight)
	for i := range grid {
		grid[i] = make([]struct {
			ch    rune
			color lipgloss.Color
		}, logoWidth)
		for j := range grid[i] {
			grid[i][j].ch = ' '
		}
	}

	// Place characters that have arrived
	for _, ch := range m.chars {
		if ch.revealed {
			x := int(math.Round(ch.currX))
			y := int(math.Round(ch.currY))
			if y >= 0 && y < logoHeight && x >= 0 && x < logoWidth {
				// Apply wave shimmer effect
				shimmerOffset := math.Sin(m.wavePhase + float64(x)*0.3)
				colorIdx := int((float64(x)/float64(logoWidth) + m.rainbowOffset + shimmerOffset*0.1) * float64(len(rainbowColors)))
				colorIdx = colorIdx % len(rainbowColors)
				if colorIdx < 0 {
					colorIdx += len(rainbowColors)
				}

				grid[y][x].ch = ch.char
				grid[y][x].color = rainbowColors[colorIdx]
			}
		}
	}

	// "Welcome to" title
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Italic(true)
	b.WriteString(pad + titleStyle.Render("  Welcome to") + "\n\n")

	// Render logo grid
	for _, row := range grid {
		b.WriteString(pad)
		for _, cell := range row {
			if cell.ch == ' ' {
				b.WriteString(" ")
			} else {
				style := lipgloss.NewStyle().Foreground(cell.color).Bold(true)
				b.WriteString(style.Render(string(cell.ch)))
			}
		}
		b.WriteString("\n")
	}

	// Tagline
	if m.phase >= 3 {
		tagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")).Italic(true)
		b.WriteString(pad + "                " + tagStyle.Render("The Dream CLI") + "\n")
	}

	b.WriteString("\n")

	// Render eye
	if m.phase >= 2 && m.eyeScale > 0.3 {
		b.WriteString(m.renderEye())
	}

	// Render particles
	if len(m.particles) > 0 {
		b.WriteString(m.renderParticles())
	}

	// Version
	if m.phase >= 4 {
		verStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B"))
		b.WriteString("\n" + pad + verStyle.Render("  CLI v0.2.0") + "\n")
	}

	return b.String()
}

func (m *EnhancedIntroModel) renderEye() string {
	var b strings.Builder

	eye := eyeFrames[m.eyeFrame]
	eyeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#22D3EE"))
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))
	pupilStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Bold(true)

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

	return b.String()
}

func (m *EnhancedIntroModel) renderParticles() string {
	var b strings.Builder

	for _, p := range m.particles {
		opacity := 1.0 - float64(p.life)/float64(p.maxLife)
		if opacity > 0.3 {
			style := lipgloss.NewStyle().Foreground(p.color)
			// Simple representation - just add sparkles to output
			b.WriteString(style.Render(p.char))
		}
	}
	if len(m.particles) > 0 {
		b.WriteString("\n")
	}

	return b.String()
}

func (m *EnhancedIntroModel) Done() bool {
	return m.done
}
