// Layout - Declarative pane layout using bubblelayout
package tui

import (
	bl "github.com/winder/bubblelayout"
)

// LayoutIDs for each component
type LayoutIDs struct {
	Header  bl.ID
	Agents  bl.ID
	Chat    bl.ID
	Input   bl.ID
}

// NewLayout creates the declarative layout
func NewLayout() (bl.BubbleLayout, LayoutIDs) {
	layout := bl.New()

	ids := LayoutIDs{}

	// Simple 2-column layout
	// Agents pane - left side, fixed width
	ids.Agents = layout.Add("width 25")

	// Chat pane - right side, grows to fill
	ids.Chat = layout.Add("grow")

	return layout, ids
}

// SimpleLayout for 2-pane split
func NewSimpleLayout() (bl.BubbleLayout, LayoutIDs) {
	layout := bl.New()

	ids := LayoutIDs{}

	// Left pane (agents) - 20% width
	ids.Agents = layout.Add("width 20%")

	// Right pane (chat) - grows
	ids.Chat = layout.Add("grow")

	return layout, ids
}
