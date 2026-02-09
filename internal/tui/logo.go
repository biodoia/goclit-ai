// Package tui - Logo animation (graphic, not text)
// Black screen → Logo appears center → Then title
package tui

import (
	"fmt"
	"strings"
	"time"
)

// GoclitLogo - Graphic robot logo using block characters
const GoclitLogo = `
            ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
          ▄█░░░░░░░░░░░░░░░█▄
         ██░░░░░░░░░░░░░░░░░██
        ██░░▓▓░░░░░░░░░▓▓░░░██
        ██░░▓▓░░░░░░░░░▓▓░░░██
        ██░░░░░░░▄▄▄░░░░░░░░██
        ██░░░░░░█████░░░░░░░██
        ██░░░░░░░░░░░░░░░░░░██
         ██░░░░░░░░░░░░░░░░██
          ▀█▄░░░░░░░░░░░░▄█▀
            ▀▀▀▀▀▀▀▀▀▀▀▀▀▀
`

// GoclitLogoSmall - Compact version
const GoclitLogoSmall = `
       ▄▄▄▄▄▄▄
      █ ◉   ◉ █
      █   ▼   █
      █ ╰───╯ █
       ▀▀▀▀▀▀▀
`

// GoclitLogoMinimal - Minimal robot face
const GoclitLogoMinimal = `
    ┌─────────┐
    │  ◉   ◉  │
    │    ▼    │
    │  ╰───╯  │
    └─────────┘
`

// AntennaBlink frames for animation
var AntennaFrames = []string{
	"       ·       ",
	"       ★       ",
	"       ✦       ",
	"       ✧       ",
	"       ★       ",
}

// PlayLogoAnimation shows logo animation (center screen, black bg)
func PlayLogoAnimation(width, height int) {
	// 1. Clear screen (black)
	fmt.Print("\033[2J")     // Clear
	fmt.Print("\033[?25l")   // Hide cursor
	fmt.Print("\033[40m")    // Black background
	
	// Fill with black
	for i := 0; i < height; i++ {
		fmt.Printf("\033[%d;0H%s", i+1, strings.Repeat(" ", width))
	}
	
	// 2. Calculate center position
	logoLines := strings.Split(GoclitLogoSmall, "\n")
	logoHeight := len(logoLines)
	logoWidth := 0
	for _, line := range logoLines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}
	
	startRow := (height - logoHeight) / 2
	startCol := (width - logoWidth) / 2
	
	// 3. Fade in logo (character by character or line by line)
	fmt.Print("\033[36m") // Cyan color
	
	// Appear effect - flash then solid
	for flash := 0; flash < 3; flash++ {
		// Flash on
		for i, line := range logoLines {
			if line == "" { continue }
			fmt.Printf("\033[%d;%dH%s", startRow+i, startCol, line)
		}
		time.Sleep(80 * time.Millisecond)
		
		// Flash off (if not last)
		if flash < 2 {
			for i := range logoLines {
				fmt.Printf("\033[%d;%dH%s", startRow+i, startCol, strings.Repeat(" ", logoWidth))
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
	
	// 4. Antenna blink animation
	antennaRow := startRow - 1
	antennaCol := startCol + (logoWidth / 2) - 7
	
	for _, frame := range AntennaFrames {
		fmt.Printf("\033[%d;%dH\033[33m%s\033[36m", antennaRow, antennaCol, frame) // Yellow antenna
		time.Sleep(100 * time.Millisecond)
	}
	
	// 5. Hold for a moment
	time.Sleep(300 * time.Millisecond)
	
	// 6. Reset
	fmt.Print("\033[0m")    // Reset colors
	fmt.Print("\033[?25h")  // Show cursor
	fmt.Print("\033[2J")    // Clear
	fmt.Print("\033[H")     // Home
}

// PlayFullIntro plays the complete intro sequence
// 1. Black screen
// 2. Logo appears (center)
// 3. Logo pulses
// 4. Title slides in
// 5. "Agents are listening..."
func PlayFullIntro(width, height int) {
	// Hide cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")
	
	// Phase 1: Black screen (500ms)
	fmt.Print("\033[2J\033[40m")
	for i := 0; i < height; i++ {
		fmt.Printf("\033[%d;0H%s", i+1, strings.Repeat(" ", width))
	}
	time.Sleep(500 * time.Millisecond)
	
	// Phase 2: Logo fades in (center)
	logoLines := strings.Split(GoclitLogoSmall, "\n")
	logoHeight := len(logoLines)
	logoWidth := 15
	startRow := (height - logoHeight) / 2
	startCol := (width - logoWidth) / 2
	
	// Pixel-by-pixel reveal (simplified: line by line)
	fmt.Print("\033[36m") // Cyan
	for i, line := range logoLines {
		if strings.TrimSpace(line) == "" { continue }
		fmt.Printf("\033[%d;%dH%s", startRow+i, startCol, line)
		time.Sleep(50 * time.Millisecond)
	}
	
	// Phase 3: Antenna blink
	antennaRow := startRow - 1
	antennaCol := startCol + 7
	for i := 0; i < 5; i++ {
		color := "\033[33m" // Yellow
		if i % 2 == 0 {
			fmt.Printf("\033[%d;%dH%s★\033[0m", antennaRow, antennaCol, color)
		} else {
			fmt.Printf("\033[%d;%dH%s·\033[0m", antennaRow, antennaCol, color)
		}
		time.Sleep(150 * time.Millisecond)
	}
	
	// Phase 4: Title appears below
	titleRow := startRow + logoHeight + 2
	title := "G O C L I T"
	titleCol := (width - len(title)) / 2
	
	fmt.Print("\033[1;36m") // Bold cyan
	for i, char := range title {
		fmt.Printf("\033[%d;%dH%c", titleRow, titleCol+i, char)
		time.Sleep(30 * time.Millisecond)
	}
	
	// Phase 5: Tagline
	tagline := "✨ Agents are listening... ✨"
	taglineRow := titleRow + 2
	taglineCol := (width - len(tagline)) / 2
	
	time.Sleep(200 * time.Millisecond)
	fmt.Print("\033[35m") // Magenta
	fmt.Printf("\033[%d;%dH%s", taglineRow, taglineCol, tagline)
	
	// Hold
	time.Sleep(800 * time.Millisecond)
	
	// Clear and continue
	fmt.Print("\033[0m\033[2J\033[H")
}

// GetTerminalSize attempts to get terminal dimensions
func GetTerminalSize() (width, height int) {
	// Default fallback
	return 80, 24
}
