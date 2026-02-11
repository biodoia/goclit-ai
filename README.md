# ⚡ Goclit-AI

> The Dream CLI - Synthesis of coding agents

## Cos'è

Goclit-AI è la fusione di 65 coding agent patterns in un unico tool:

- **Multi-model** - Routing intelligente tra provider
- **MCP Integration** - Model Context Protocol
- **Specialized agents** - Hephaestus, Sisyphus, etc.
- **Animated TUI** - Robot mascot + modern UI

## Installazione

```bash
go install github.com/biodoia/goclit-ai/cmd/goclit@latest
```

## Usage

### TUI Mode
```bash
# Avvia TUI interattiva
goclit
```

### One-Shot
```bash
# Query diretta
goclit "Fix the bug in main.go"
```

### Demo
```bash
# Demo introduttiva
goclit --intro
```

## Agenti Specializzati

| Agent | Focus |
|-------|-------|
| **Hephaestus** | Code forging, creation |
| **Sisyphus** | Repetitive tasks, persistence |
| **Architect** | System design |
| **Coder** | Implementation |
| **Reviewer** | Code review |
| **Debugger** | Bug fixing |

## Smart Routing

```
Task → Analyze → Select Best Model → Execute

Low complexity → Fast model (GLM-4-32B, Gemini Flash)
High complexity → Deep model (Claude, GPT-4)
Code generation → CodeGeeX, Codex
```

## Model Registry

| Provider | Models |
|----------|--------|
| Z.AI | GLM-4-32B, GLM-Z1-32B, CodeGeeX-4 |
| OpenRouter | Free models rotation |
| Anthropic | Claude Opus/Sonnet |
| OpenAI | GPT-4, Codex |
| Google | Gemini Pro/Flash |

## MCP Support

```bash
# Con MCP servers
goclit --mcp-config ~/.goclit/mcp.json
```

## Ultrawork Mode

```bash
# Modalità intensiva (no pause)
goclit --ultrawork "Refactor entire codebase"
```

## Configurazione

```yaml
# ~/.goclit/config.yaml
default_model: glm-4-32b
smart_routing: true
mcp_enabled: true

providers:
  zai:
    api_key: $ZAI_API_KEY
  openrouter:
    api_key: $OPENROUTER_API_KEY
```

## TUI Features

```
      ★      
   ▄▄▄▄▄▄▄   
   █ ◉ ◉ █   ← Animated robot mascot
   █  ▼  █   
   █ ╰─╯ █   
   ▀▀▀▀▀▀▀   
```

- Multicolor animations
- Fly-in effects
- Spring physics
- Syntax highlighting
- File tree browser

## Architettura

```
goclit-ai/
├── cmd/
│   ├── goclit/      # Main CLI
│   └── introdemo/   # Demo animation
├── internal/
│   ├── tui/         # Bubble Tea UI
│   ├── agents/      # Specialized agents
│   ├── providers/   # LLM providers
│   ├── mcp/         # MCP manager
│   └── core/        # Ultrawork engine
└── assets/
    └── animations/  # TUI animations
```

---

*Part of the Autoschei ecosystem*
