// goclit - The Dream CLI
// Synthesis of 65 coding agents into one supreme tool
// Pattern from oh-my-opencode
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/biodoia/goclit-ai/internal/tui"
)

const version = "0.2.0"

// Robot logo for animated intro
var robotLogo = []string{
	"      ★      ",
	"   ▄▄▄▄▄▄▄   ",
	"   █ ◉ ◉ █   ",
	"   █  ▼  █   ",
	"   █ ╰─╯ █   ",
	"   ▀▀▀▀▀▀▀   ",
}

var robotLogoNoAntenna = []string{
	"             ",
	"   ▄▄▄▄▄▄▄   ",
	"   █ ◉ ◉ █   ",
	"   █  ▼  █   ",
	"   █ ╰─╯ █   ",
	"   ▀▀▀▀▀▀▀   ",
}

const help = `
USAGE:
  goclit [command] [options]

COMMANDS:
  ultrawork <task>    🚀 Total automation mode (from oh-my-opencode)
  chat <message>      💬 Chat with the AI
  agents              🤖 List available agents
  mcp                 🔌 MCP server management
  oracle <question>   🔮 Ask the Oracle
  librarian <query>   📚 Search with Librarian
  sisyphus <task>     ⚙️  Start Sisyphus (discipline agent)
  hephaestus <task>   🔨 Start Hephaestus (autonomy agent)
  
OPTIONS:
  --help, -h          Show this help
  --version, -v       Show version
  --provider <name>   Set AI provider (default: claude)

EXAMPLES:
  goclit ultrawork "build a REST API for user management"
  goclit chat "explain this code"
  goclit oracle "what's the best approach for X?"
  goclit sisyphus "refactor module Y"

MAGIC WORDS:
  ultrawork           Total automation - all agents work together
  
SPECIALIZED AGENTS:
  Oracle              Knowledge & decision making
  Librarian           Documentation & search
  Sisyphus            Discipline - works until done
  Hephaestus          Autonomy - builds independently
  Frontend Engineer   UI/UX specialist
  Backend Engineer    Server-side specialist
  DevOps Engineer     Infrastructure specialist
`

func main() {
	if len(os.Args) < 2 {
		// Interactive TUI mode (no args)
		if err := tui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "tui":
		// Explicit TUI mode
		if err := tui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "--help", "-h", "help":
		printBanner()
		fmt.Println(help)

	case "--version", "-v", "version":
		fmt.Printf("goclit v%s\n", version)
		fmt.Println("The Dream CLI - Synthesis of 65 coding agents")

	case "ultrawork":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclit ultrawork <task>")
			os.Exit(1)
		}
		task := strings.Join(os.Args[2:], " ")
		runUltrawork(task)

	case "chat":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclit chat <message>")
			os.Exit(1)
		}
		message := strings.Join(os.Args[2:], " ")
		runChat(message)

	case "agents":
		listAgents()

	case "mcp":
		runMCP()

	case "oracle":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclit oracle <question>")
			os.Exit(1)
		}
		question := strings.Join(os.Args[2:], " ")
		runOracle(question)

	case "librarian":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclit librarian <query>")
			os.Exit(1)
		}
		query := strings.Join(os.Args[2:], " ")
		runLibrarian(query)

	case "sisyphus":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclit sisyphus <task>")
			os.Exit(1)
		}
		task := strings.Join(os.Args[2:], " ")
		runSisyphus(task)

	case "hephaestus":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclit hephaestus <task>")
			os.Exit(1)
		}
		task := strings.Join(os.Args[2:], " ")
		runHephaestus(task)

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Run 'goclit --help' for usage")
		os.Exit(1)
	}
}

func printBanner() {
	if isInteractive() && !noAnimation() {
		playAnimatedBanner()
	} else {
		// Static fallback
		printStaticBanner()
	}
}

func printStaticBanner() {
	cyan := "\033[36m"
	reset := "\033[0m"
	for _, line := range robotLogo {
		fmt.Println(cyan + line + reset)
	}
	fmt.Println()
	fmt.Println("  G O C L I T")
	fmt.Printf("  v%s - The Dream CLI\n\n", version)
}

func isInteractive() bool {
	return true
}

func noAnimation() bool {
	return os.Getenv("NO_ANIMATION") == "1" || os.Getenv("GOCLIT_NO_ANIMATION") == "1"
}

func playAnimatedBanner() {
	cyan := "\033[36m"
	purple := "\033[35m"
	white := "\033[97m"
	dim := "\033[2m"
	reset := "\033[0m"
	clear := "\033[2J\033[H"
	home := "\033[H"

	// Phase 1: Black screen (300ms)
	fmt.Print(clear)
	time.Sleep(300 * time.Millisecond)

	// Phase 2: Logo with flicker (600ms)
	for i := 0; i < 8; i++ {
		fmt.Print(home)
		fmt.Println()
		fmt.Println()

		// Flicker effect - random visibility
		if i%3 == 0 || i > 4 {
			// Show logo
			color := purple
			if i > 5 {
				color = cyan
			}
			for _, line := range robotLogoNoAntenna {
				fmt.Println(color + "        " + line + reset)
			}
		} else {
			// Blank during flicker
			for range robotLogoNoAntenna {
				fmt.Println()
			}
		}
		time.Sleep(75 * time.Millisecond)
	}

	// Phase 3: Antenna flicker (400ms)
	for i := 0; i < 6; i++ {
		fmt.Print(home)
		fmt.Println()
		fmt.Println()

		var logo []string
		if i%2 == 0 {
			logo = robotLogo // With antenna
		} else {
			logo = robotLogoNoAntenna // Without antenna
		}

		for _, line := range logo {
			fmt.Println(cyan + "        " + line + reset)
		}
		time.Sleep(70 * time.Millisecond)
	}

	// Final logo stable
	fmt.Print(home)
	fmt.Println()
	fmt.Println()
	for _, line := range robotLogo {
		fmt.Println(cyan + "        " + line + reset)
	}
	fmt.Println()

	// Phase 4: GOCLIT letter by letter (500ms)
	title := "G O C L I T"
	fmt.Print("        ")
	for _, ch := range title {
		fmt.Print(white + string(ch) + reset)
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println()
	time.Sleep(100 * time.Millisecond)

	// Phase 5: Tagline
	fmt.Println(dim + "        The Dream CLI" + reset)
	fmt.Println(cyan + "        v" + version + reset)
	fmt.Println()

	// Phase 6: Agents listening with sparkle
	sparkles := []string{"✨", "⚡", "💫", "🌟"}
	for i := 0; i < 3; i++ {
		s := sparkles[i%len(sparkles)]
		fmt.Print("\r        " + s + " Agents are listening... " + s + "  ")
		time.Sleep(150 * time.Millisecond)
	}
	fmt.Println()
	fmt.Println()
}

func runUltrawork(task string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n⏸️  Pausing ultrawork...")
		cancel()
	}()

	fmt.Println("🚀 ULTRAWORK MODE ACTIVATED")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Task: %s\n\n", task)
	fmt.Println("⚠️  Provider not configured. Run: goclit config --provider <name>")
	fmt.Println("📋 Supported: claude, openai, gemini, groq, local")

	// TODO: Implement actual ultrawork
	_ = ctx
}

func runChat(message string) {
	fmt.Println("💬 Chat mode")
	fmt.Printf("You: %s\n\n", message)
	fmt.Println("⚠️  Provider not configured.")
}

func listAgents() {
	fmt.Println("🤖 Available Agents:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	agents := []struct {
		name string
		role string
		icon string
	}{
		{"Sisyphus", "Discipline agent - works until done", "⚙️"},
		{"Hephaestus", "Autonomy agent - builds independently", "🔨"},
		{"Oracle", "Knowledge & decision making", "🔮"},
		{"Librarian", "Documentation & search", "📚"},
		{"Frontend Engineer", "UI/UX specialist (React)", "🎨"},
		{"Backend Engineer", "Server-side specialist (Go)", "⚡"},
		{"DevOps Engineer", "Infrastructure specialist", "🔧"},
	}

	for _, a := range agents {
		fmt.Printf("  %s %s\n     └─ %s\n", a.icon, a.name, a.role)
	}
}

func runMCP() {
	fmt.Println("🔌 MCP Servers (Curated):")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	servers := []struct {
		name  string
		tools int
	}{
		{"filesystem", 6},
		{"shell", 2},
		{"git", 6},
		{"search", 3},
		{"browser", 5},
		{"database", 2},
	}

	for _, s := range servers {
		fmt.Printf("  📦 %s (%d tools)\n", s.name, s.tools)
	}
}

func runOracle(question string) {
	fmt.Println("🔮 Oracle")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Question: %s\n\n", question)
	fmt.Println("⚠️  Provider not configured.")
}

func runLibrarian(query string) {
	fmt.Println("📚 Librarian")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Query: %s\n\n", query)
	fmt.Println("⚠️  Provider not configured.")
}

func runSisyphus(task string) {
	fmt.Println("⚙️  Sisyphus - The Discipline Agent")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Task: %s\n\n", task)
	fmt.Println("\"It just works until the task is done.\"")
	fmt.Println("\n⚠️  Provider not configured.")
}

func runHephaestus(task string) {
	fmt.Println("🔨 Hephaestus - The Autonomy Agent")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Task: %s\n\n", task)
	fmt.Println("\"God of the forge - creates and builds independently.\"")
	fmt.Println("\n⚠️  Provider not configured.")
}
