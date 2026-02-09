// Theme - Copilot-style colors and styling
package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Copilot gradient colors (purple → blue → cyan)
var (
	// Primary gradient
	Purple     = lipgloss.Color("#A855F7") // violet-500
	PurpleDark = lipgloss.Color("#7C3AED") // violet-600
	Blue       = lipgloss.Color("#3B82F6") // blue-500
	BlueDark   = lipgloss.Color("#2563EB") // blue-600
	Cyan       = lipgloss.Color("#06B6D4") // cyan-500
	CyanDark   = lipgloss.Color("#0891B2") // cyan-600

	// Accent colors
	Green      = lipgloss.Color("#22C55E") // green-500
	Yellow     = lipgloss.Color("#EAB308") // yellow-500
	Orange     = lipgloss.Color("#F97316") // orange-500
	Red        = lipgloss.Color("#EF4444") // red-500
	Pink       = lipgloss.Color("#EC4899") // pink-500

	// Neutrals
	White      = lipgloss.Color("#FAFAFA")
	Gray100    = lipgloss.Color("#F3F4F6")
	Gray300    = lipgloss.Color("#D1D5DB")
	Gray500    = lipgloss.Color("#6B7280")
	Gray700    = lipgloss.Color("#374151")
	Gray900    = lipgloss.Color("#111827")
	Black      = lipgloss.Color("#0A0A0A")

	// Background
	BgDark     = lipgloss.Color("#0D1117") // GitHub dark bg
	BgPanel    = lipgloss.Color("#161B22") // Panel bg
	BgHighlight = lipgloss.Color("#21262D") // Highlight bg
)

// Gradient returns interpolated color for smooth animation
func Gradient(progress float64) lipgloss.Color {
	// 0.0 = purple, 0.5 = blue, 1.0 = cyan
	if progress < 0.5 {
		// Purple → Blue
		return interpolateColor(Purple, Blue, progress*2)
	}
	// Blue → Cyan
	return interpolateColor(Blue, Cyan, (progress-0.5)*2)
}

func interpolateColor(c1, c2 lipgloss.Color, t float64) lipgloss.Color {
	// Simple hex interpolation
	r1, g1, b1 := hexToRGB(string(c1))
	r2, g2, b2 := hexToRGB(string(c2))

	r := int(float64(r1) + t*(float64(r2)-float64(r1)))
	g := int(float64(g1) + t*(float64(g2)-float64(g1)))
	b := int(float64(b1) + t*(float64(b2)-float64(b1)))

	return lipgloss.Color(rgbToHex(r, g, b))
}

func hexToRGB(hex string) (r, g, b int) {
	if len(hex) < 7 {
		return 128, 128, 128
	}
	hex = hex[1:] // Remove #
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return
}

func rgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// Styles
var (
	// Logo style with glow effect
	LogoStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Purple)

	// Title style
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(White).
		Background(BgDark).
		Padding(0, 1)

	// Subtitle
	SubtitleStyle = lipgloss.NewStyle().
		Foreground(Gray500).
		Italic(true)

	// Panel border
	PanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Gray700).
		Background(BgPanel).
		Padding(1, 2)

	// Active panel
	ActivePanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Cyan).
		Background(BgPanel).
		Padding(1, 2)

	// Status bar
	StatusStyle = lipgloss.NewStyle().
		Background(BgHighlight).
		Foreground(Gray300).
		Padding(0, 1)

	// Agent badge
	AgentBadge = func(color lipgloss.Color) lipgloss.Style {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(Black).
			Background(color).
			Padding(0, 1).
			MarginRight(1)
	}

	// Command input
	InputStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(Blue).
		Padding(0, 1)

	// Spinner
	SpinnerStyle = lipgloss.NewStyle().
		Foreground(Cyan)

	// Success/Error
	SuccessStyle = lipgloss.NewStyle().Foreground(Green)
	ErrorStyle   = lipgloss.NewStyle().Foreground(Red)
	WarnStyle    = lipgloss.NewStyle().Foreground(Yellow)
)
