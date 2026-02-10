// Demo for the enhanced intro animation
// Run: go run ./cmd/introdemo
package main

import (
	"fmt"
	"os"

	"github.com/biodoia/goclit-ai/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create enhanced intro model with reasonable size
	intro := tui.NewEnhancedIntro(80, 24)

	p := tea.NewProgram(intro, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
