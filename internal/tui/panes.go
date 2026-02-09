// Panes - Split pane layout with bubbletea
package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PaneID identifies a pane
type PaneID int

const (
	PaneAgents PaneID = iota
	PaneChat
	PaneInput
)

// Pane is a single pane with viewport
type Pane struct {
	ID       PaneID
	Title    string
	viewport viewport.Model
	focused  bool
	width    int
	height   int
}

// PaneLayout manages multiple panes
type PaneLayout struct {
	panes       []*Pane
	focusedPane PaneID
	width       int
	height      int

	// Layout config
	leftWidth  float64 // 0.0 to 1.0
	showHeader bool
	showFooter bool
}

// NewPaneLayout creates a new pane layout
func NewPaneLayout() *PaneLayout {
	agentsPane := &Pane{
		ID:       PaneAgents,
		Title:    "AGENTS",
		viewport: viewport.New(0, 0),
	}

	chatPane := &Pane{
		ID:       PaneChat,
		Title:    "CHAT",
		viewport: viewport.New(0, 0),
	}

	return &PaneLayout{
		panes:       []*Pane{agentsPane, chatPane},
		focusedPane: PaneChat,
		leftWidth:   0.25, // 25% for agents
		showHeader:  true,
		showFooter:  true,
	}
}

// SetSize updates layout dimensions
func (pl *PaneLayout) SetSize(width, height int) {
	pl.width = width
	pl.height = height

	// Calculate pane dimensions
	headerHeight := 0
	footerHeight := 0
	if pl.showHeader {
		headerHeight = 1
	}
	if pl.showFooter {
		footerHeight = 1
	}

	contentHeight := height - headerHeight - footerHeight - 2 // borders

	leftWidth := int(float64(width) * pl.leftWidth)
	rightWidth := width - leftWidth - 3 // separator

	// Update pane sizes
	for _, p := range pl.panes {
		switch p.ID {
		case PaneAgents:
			p.width = leftWidth
			p.height = contentHeight
			p.viewport.Width = leftWidth - 4
			p.viewport.Height = contentHeight - 2
		case PaneChat:
			p.width = rightWidth
			p.height = contentHeight
			p.viewport.Width = rightWidth - 4
			p.viewport.Height = contentHeight - 2
		}
	}
}

// FocusNext moves focus to the next pane
func (pl *PaneLayout) FocusNext() {
	pl.focusedPane = (pl.focusedPane + 1) % PaneID(len(pl.panes))
	pl.updateFocus()
}

// FocusPrev moves focus to the previous pane
func (pl *PaneLayout) FocusPrev() {
	pl.focusedPane = (pl.focusedPane - 1 + PaneID(len(pl.panes))) % PaneID(len(pl.panes))
	pl.updateFocus()
}

func (pl *PaneLayout) updateFocus() {
	for _, p := range pl.panes {
		p.focused = p.ID == pl.focusedPane
	}
}

// GetPane returns a pane by ID
func (pl *PaneLayout) GetPane(id PaneID) *Pane {
	for _, p := range pl.panes {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// SetContent sets content for a pane
func (pl *PaneLayout) SetContent(id PaneID, content string) {
	if p := pl.GetPane(id); p != nil {
		p.viewport.SetContent(content)
	}
}

// Update handles messages for the focused pane
func (pl *PaneLayout) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			pl.FocusNext()
			return nil
		case "shift+tab":
			pl.FocusPrev()
			return nil
		}
	}

	// Update focused pane's viewport
	if p := pl.GetPane(pl.focusedPane); p != nil {
		var cmd tea.Cmd
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// View renders the pane layout
func (pl *PaneLayout) View() string {
	if pl.width == 0 || pl.height == 0 {
		return ""
	}

	var sections []string

	// Header
	if pl.showHeader {
		sections = append(sections, pl.renderHeader())
	}

	// Panes
	sections = append(sections, pl.renderPanes())

	// Footer
	if pl.showFooter {
		sections = append(sections, pl.renderFooter())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (pl *PaneLayout) renderHeader() string {
	logo := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true).
		Render("⚡ GOCLIT")

	ver := lipgloss.NewStyle().
		Foreground(Gray500).
		Render(" v0.2.0")

	return lipgloss.NewStyle().
		Width(pl.width).
		Background(BgHighlight).
		Padding(0, 1).
		Render(logo + ver)
}

func (pl *PaneLayout) renderPanes() string {
	var paneViews []string

	for _, p := range pl.panes {
		paneViews = append(paneViews, pl.renderPane(p))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, paneViews...)
}

func (pl *PaneLayout) renderPane(p *Pane) string {
	// Border style based on focus
	borderColor := Gray700
	if p.focused {
		borderColor = Cyan
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(White).
		Padding(0, 1)

	if p.focused {
		titleStyle = titleStyle.Foreground(Cyan)
	}

	title := titleStyle.Render(p.Title)

	// Content
	contentStyle := lipgloss.NewStyle().
		Width(p.width - 2).
		Height(p.height - 3).
		Padding(0, 1)

	content := contentStyle.Render(p.viewport.View())

	// Border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(p.width).
		Height(p.height)

	inner := lipgloss.JoinVertical(lipgloss.Left, title, content)
	return borderStyle.Render(inner)
}

func (pl *PaneLayout) renderFooter() string {
	left := lipgloss.NewStyle().
		Foreground(Gray500).
		Render("Tab: switch pane • ?: help • q: quit")

	focusedName := ""
	if p := pl.GetPane(pl.focusedPane); p != nil {
		focusedName = p.Title
	}

	right := lipgloss.NewStyle().
		Foreground(Cyan).
		Render("Focus: " + focusedName)

	gap := pl.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 0 {
		gap = 1
	}

	return lipgloss.NewStyle().
		Width(pl.width).
		Background(BgHighlight).
		Padding(0, 1).
		Render(left + strings.Repeat(" ", gap) + right)
}

// AgentItem for the agents pane
type AgentItem struct {
	Icon   string
	Name   string
	Role   string
	Color  lipgloss.Color
	Active bool
}

// DefaultAgents returns the default agent list
func DefaultAgents() []AgentItem {
	return []AgentItem{
		{"⚙️", "Sisyphus", "Discipline", Purple, false},
		{"🔨", "Hephaestus", "Autonomy", Blue, false},
		{"🔮", "Oracle", "Knowledge", Cyan, false},
		{"📚", "Librarian", "Search", Green, false},
		{"🎨", "Frontend", "UI/UX", Pink, false},
		{"⚡", "Backend", "Server", Orange, false},
		{"🔧", "DevOps", "Infra", Yellow, false},
	}
}

// RenderAgentList renders the agent list for the agents pane
func RenderAgentList(agents []AgentItem, selected int) string {
	var lines []string

	for i, a := range agents {
		isSelected := i == selected

		// Cursor
		cursor := "  "
		if isSelected {
			cursor = lipgloss.NewStyle().Foreground(Cyan).Render("▸ ")
		}

		// Icon with color
		iconStyle := lipgloss.NewStyle()
		if isSelected {
			iconStyle = iconStyle.Background(a.Color).Foreground(Black)
		}
		icon := iconStyle.Render(a.Icon)

		// Name
		nameStyle := lipgloss.NewStyle().Foreground(Gray300)
		if isSelected {
			nameStyle = nameStyle.Foreground(White).Bold(true)
		}
		name := nameStyle.Render(" " + a.Name)

		// Role (dimmed)
		roleStyle := lipgloss.NewStyle().Foreground(Gray500).Italic(true)
		role := roleStyle.Render(" - " + a.Role)

		lines = append(lines, cursor+icon+name+role)
	}

	return strings.Join(lines, "\n")
}
