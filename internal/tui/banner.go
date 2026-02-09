// Package tui implements terminal UI with animated ASCII banners
// Inspired by: GitHub Copilot CLI + Monkey Island "monkeys are listening"
// Using: Charmbracelet Harmonica for spring physics animations
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

// Banner frames - Monkey Island style "Agents are listening..."
var goclitBannerFrames = []string{
	// Frame 1 - Robot appears
	`
                    ▄▄▄▄▄▄▄
                   █ ◉   ◉ █
                   █   ▼   █
                   █ ╰───╯ █
                    ▀▀▀▀▀▀▀
`,
	// Frame 2 - Eyes move
	`
                    ▄▄▄▄▄▄▄
                   █  ◉ ◉  █
                   █   ▼   █
                   █ ╰───╯ █
                    ▀▀▀▀▀▀▀
`,
	// Frame 3 - Antenna blinks
	`
                       ★
                    ▄▄▄▄▄▄▄
                   █ ◉   ◉ █
                   █   ▼   █
                   █ ╰───╯ █
                    ▀▀▀▀▀▀▀
`,
}

// Full ASCII banner with spring animation positions
const goclitLogo = `
   ██████╗  ██████╗  ██████╗██╗     ██╗████████╗
  ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝
  ██║  ███╗██║   ██║██║     ██║     ██║   ██║   
  ██║   ██║██║   ██║██║     ██║     ██║   ██║   
  ╚██████╔╝╚██████╔╝╚██████╗███████╗██║   ██║   
   ╚═════╝  ╚═════╝  ╚═════╝╚══════╝╚═╝   ╚═╝   
`

const tagline = `          ✨ Agents are listening... ✨`

// ColorRole defines semantic color roles (from GitHub Copilot CLI article)
type ColorRole string

const (
	RoleLogo       ColorRole = "logo"
	RoleTagline    ColorRole = "tagline"
	RoleRobot      ColorRole = "robot"
	RoleEyes       ColorRole = "eyes"
	RoleAntenna    ColorRole = "antenna"
	RoleBorder     ColorRole = "border"
	RoleBackground ColorRole = "background"
)

// Theme maps color roles to actual colors (graceful degradation)
type Theme struct {
	Colors map[ColorRole]lipgloss.Color
}

// DefaultTheme for most terminals
var DefaultTheme = Theme{
	Colors: map[ColorRole]lipgloss.Color{
		RoleLogo:       lipgloss.Color("39"),  // Cyan
		RoleTagline:    lipgloss.Color("213"), // Pink
		RoleRobot:      lipgloss.Color("252"), // Light gray
		RoleEyes:       lipgloss.Color("226"), // Yellow
		RoleAntenna:    lipgloss.Color("226"), // Yellow
		RoleBorder:     lipgloss.Color("240"), // Dark gray
		RoleBackground: lipgloss.Color("0"),   // Black
	},
}

// DarkTheme for dark terminals
var DarkTheme = Theme{
	Colors: map[ColorRole]lipgloss.Color{
		RoleLogo:       lipgloss.Color("51"),  // Bright cyan
		RoleTagline:    lipgloss.Color("219"), // Bright pink
		RoleRobot:      lipgloss.Color("255"), // White
		RoleEyes:       lipgloss.Color("220"), // Gold
		RoleAntenna:    lipgloss.Color("220"), // Gold
		RoleBorder:     lipgloss.Color("245"), // Gray
		RoleBackground: lipgloss.Color("0"),   // Black
	},
}

// AccessibleTheme for high contrast
var AccessibleTheme = Theme{
	Colors: map[ColorRole]lipgloss.Color{
		RoleLogo:       lipgloss.Color("15"),  // White
		RoleTagline:    lipgloss.Color("15"),  // White
		RoleRobot:      lipgloss.Color("15"),  // White
		RoleEyes:       lipgloss.Color("11"),  // Bright yellow
		RoleAntenna:    lipgloss.Color("11"),  // Bright yellow
		RoleBorder:     lipgloss.Color("15"),  // White
		RoleBackground: lipgloss.Color("0"),   // Black
	},
}

// BannerAnimation manages the animated banner
type BannerAnimation struct {
	spring      harmonica.Spring
	xPos        float64
	xVel        float64
	yPos        float64
	yVel        float64
	targetX     float64
	targetY     float64
	frame       int
	theme       Theme
	showRobot   bool
	blinkState  bool
	startTime   time.Time
}

// NewBannerAnimation creates an animated banner with spring physics
func NewBannerAnimation(theme Theme) *BannerAnimation {
	return &BannerAnimation{
		// Spring: 60fps, angular freq 6.0, damping 0.5
		spring:    harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5),
		xPos:      -50, // Start off-screen left
		targetX:   0,   // Slide to center
		yPos:      -10, // Start above
		targetY:   0,   // Slide down
		theme:     theme,
		showRobot: true,
		startTime: time.Now(),
	}
}

// Update advances the animation by one frame
func (b *BannerAnimation) Update() {
	// Spring physics for smooth motion
	b.xPos, b.xVel = b.spring.Update(b.xPos, b.xVel, b.targetX)
	b.yPos, b.yVel = b.spring.Update(b.yPos, b.yVel, b.targetY)

	// Blink antenna every 500ms
	elapsed := time.Since(b.startTime)
	b.blinkState = (elapsed.Milliseconds()/500)%2 == 0

	// Cycle robot frames every 300ms
	b.frame = int(elapsed.Milliseconds()/300) % len(goclitBannerFrames)
}

// Render returns the current frame as a string
func (b *BannerAnimation) Render() string {
	var sb strings.Builder

	// Logo style with color role
	logoStyle := lipgloss.NewStyle().
		Foreground(b.theme.Colors[RoleLogo]).
		Bold(true)

	// Tagline style
	taglineStyle := lipgloss.NewStyle().
		Foreground(b.theme.Colors[RoleTagline]).
		Italic(true)

	// Robot style
	robotStyle := lipgloss.NewStyle().
		Foreground(b.theme.Colors[RoleRobot])

	// Apply horizontal offset (spring animation)
	indent := strings.Repeat(" ", int(b.xPos)+10)

	// Render logo with spring offset
	logoLines := strings.Split(goclitLogo, "\n")
	for _, line := range logoLines {
		sb.WriteString(indent)
		sb.WriteString(logoStyle.Render(line))
		sb.WriteString("\n")
	}

	// Render tagline
	sb.WriteString(indent)
	sb.WriteString(taglineStyle.Render(tagline))
	sb.WriteString("\n\n")

	// Render robot frame
	if b.showRobot {
		robotFrame := goclitBannerFrames[b.frame]
		robotLines := strings.Split(robotFrame, "\n")
		for _, line := range robotLines {
			sb.WriteString(indent)
			sb.WriteString(robotStyle.Render(line))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// IsComplete returns true when animation has settled
func (b *BannerAnimation) IsComplete() bool {
	// Animation complete when velocity is near zero
	return abs(b.xVel) < 0.01 && abs(b.yVel) < 0.01
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// PlayBanner plays the animated banner (blocking)
func PlayBanner(theme Theme, duration time.Duration) {
	banner := NewBannerAnimation(theme)
	ticker := time.NewTicker(time.Second / 60) // 60 FPS
	defer ticker.Stop()

	timeout := time.After(duration)

	// Clear screen
	fmt.Print("\033[2J\033[H")

	for {
		select {
		case <-timeout:
			return
		case <-ticker.C:
			banner.Update()

			// Clear and redraw
			fmt.Print("\033[H") // Cursor home
			fmt.Print(banner.Render())

			// Exit early if settled and past minimum time
			if banner.IsComplete() && time.Since(banner.startTime) > time.Second {
				time.Sleep(500 * time.Millisecond) // Pause to appreciate
				return
			}
		}
	}
}

// PlayBannerAsync plays banner in background (non-blocking)
func PlayBannerAsync(theme Theme, duration time.Duration, done chan<- struct{}) {
	go func() {
		PlayBanner(theme, duration)
		done <- struct{}{}
	}()
}

// QuickBanner shows a static banner (for accessibility/no-animation mode)
func QuickBanner(theme Theme) string {
	logoStyle := lipgloss.NewStyle().
		Foreground(theme.Colors[RoleLogo]).
		Bold(true)

	taglineStyle := lipgloss.NewStyle().
		Foreground(theme.Colors[RoleTagline]).
		Italic(true)

	return logoStyle.Render(goclitLogo) + "\n" + taglineStyle.Render(tagline) + "\n"
}
