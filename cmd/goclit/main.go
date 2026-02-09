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
)

const version = "0.2.0"

const banner = `
   ██████╗  ██████╗  ██████╗██╗     ██╗████████╗
  ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝
  ██║  ███╗██║   ██║██║     ██║     ██║   ██║   
  ██║   ██║██║   ██║██║     ██║     ██║   ██║   
  ╚██████╔╝╚██████╔╝╚██████╗███████╗██║   ██║   
   ╚═════╝  ╚═════╝  ╚═════╝╚══════╝╚═╝   ╚═╝   
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  The Dream CLI - Synthesis of 65 coding agents
  v%s
`

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
		printBanner()
		fmt.Println(help)
		return
	}

	cmd := os.Args[1]

	switch cmd {
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
	// Check if animation is enabled (default: yes for interactive terminals)
	if isInteractive() && !noAnimation() {
		// Animated banner with spring physics
		playAnimatedBanner()
	} else {
		// Static fallback
		fmt.Printf(banner, version)
	}
}

func isInteractive() bool {
	// TODO: Check if stdout is a TTY
	return true
}

func noAnimation() bool {
	// Check NO_ANIMATION or ACCESSIBILITY env vars
	return os.Getenv("NO_ANIMATION") == "1" || os.Getenv("GOCLIT_NO_ANIMATION") == "1"
}

func playAnimatedBanner() {
	// Quick animation (1.5 seconds)
	fmt.Print("\033[2J\033[H") // Clear screen
	
	frames := []string{
		// Frame 1 - Logo slides in
		`
   ██████╗  ██████╗  ██████╗██╗     ██╗████████╗
  ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝
`,
		// Frame 2 - More logo
		`
   ██████╗  ██████╗  ██████╗██╗     ██╗████████╗
  ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝
  ██║  ███╗██║   ██║██║     ██║     ██║   ██║   
  ██║   ██║██║   ██║██║     ██║     ██║   ██║   
`,
		// Frame 3 - Full logo
		`
   ██████╗  ██████╗  ██████╗██╗     ██╗████████╗
  ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝
  ██║  ███╗██║   ██║██║     ██║     ██║   ██║   
  ██║   ██║██║   ██║██║     ██║     ██║   ██║   
  ╚██████╔╝╚██████╔╝╚██████╗███████╗██║   ██║   
   ╚═════╝  ╚═════╝  ╚═════╝╚══════╝╚═╝   ╚═╝   
`,
		// Frame 4 - With tagline
		`
   ██████╗  ██████╗  ██████╗██╗     ██╗████████╗
  ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝
  ██║  ███╗██║   ██║██║     ██║     ██║   ██║   
  ██║   ██║██║   ██║██║     ██║     ██║   ██║   
  ╚██████╔╝╚██████╔╝╚██████╗███████╗██║   ██║   
   ╚═════╝  ╚═════╝  ╚═════╝╚══════╝╚═╝   ╚═╝   
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`,
		// Frame 5 - With robot
		`
   ██████╗  ██████╗  ██████╗██╗     ██╗████████╗
  ██╔════╝ ██╔═══██╗██╔════╝██║     ██║╚══██╔══╝
  ██║  ███╗██║   ██║██║     ██║     ██║   ██║   
  ██║   ██║██║   ██║██║     ██║     ██║   ██║   
  ╚██████╔╝╚██████╔╝╚██████╗███████╗██║   ██║   
   ╚═════╝  ╚═════╝  ╚═════╝╚══════╝╚═╝   ╚═╝   
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
        ✨ Agents are listening... ✨  🤖
`,
	}

	for _, frame := range frames {
		fmt.Print("\033[H") // Cursor home
		fmt.Print("\033[36m") // Cyan
		fmt.Print(frame)
		fmt.Print("\033[0m") // Reset
		time.Sleep(150 * time.Millisecond)
	}
	
	fmt.Printf("\n  v%s - The Dream CLI\n\n", version)
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
