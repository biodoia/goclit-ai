# TUI Engineering Reference

## GitHub Copilot CLI ASCII Banner
Source: https://github.blog/engineering/from-pixels-to-characters-the-engineering-behind-github-copilot-clis-animated-ascii-banner/

### Key Learnings

1. **Color Roles (not RGB)**
   - Don't hardcode colors
   - Use semantic roles: "eyes", "shadow", "border"
   - Map to ANSI colors that degrade gracefully

2. **Accessibility First**
   - Screen readers: avoid rapid re-renders
   - Color-blind: meaning not only from color
   - Low-vision: high contrast support
   - Animations: opt-in, not automatic

3. **Terminal Inconsistencies**
   - ANSI codes behave differently per terminal
   - 3-bit, 4-bit, 8-bit, truecolor support varies
   - Test on: iTerm2, Terminal.app, Windows Terminal, Alacritty, etc.

4. **Ink Limitations**
   - Re-renders on every state change
   - No frame delta management
   - No flicker prevention
   - Need custom animation logic

5. **Custom Tooling Needed**
   - Frame-by-frame editor
   - Multi-color ANSI preview
   - Color role export
   - Accessibility testing

### For goclit TUI

```go
// Color roles for theming
type ColorRole string

const (
    RolePrimary    ColorRole = "primary"    // Main text
    RoleSecondary  ColorRole = "secondary"  // Dimmed text
    RoleAccent     ColorRole = "accent"     // Highlights
    RoleSuccess    ColorRole = "success"    // Green indicators
    RoleWarning    ColorRole = "warning"    // Yellow/orange
    RoleError      ColorRole = "error"      // Red indicators
    RoleBorder     ColorRole = "border"     // Box borders
    RoleBackground ColorRole = "background" // Background
)

// Theme maps roles to actual ANSI colors
type Theme map[ColorRole]lipgloss.Color
```

---
*Reference added 2026-02-09*
