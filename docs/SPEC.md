# 🚀 DREAM CLI - Il Super Coding Agent

> Specifica del CLI definitivo che unisce le migliori features di 65 tool

---

## 📐 Architettura

```
┌─────────────────────────────────────────────────────────────────┐
│                         DREAM CLI                                │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   TUI/GUI   │  │   Desktop   │  │      IDE Extension      │  │
│  │  (OpenCode) │  │    App      │  │   (Continue-style)      │  │
│  └──────┬──────┘  └──────┬──────┘  └───────────┬─────────────┘  │
│         └────────────────┼─────────────────────┘                │
│                          ▼                                       │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    CORE (Rust/Go)                          │  │
│  │  ┌─────────┐ ┌──────────┐ ┌─────────┐ ┌────────────────┐  │  │
│  │  │ RepoMap │ │ MCP Core │ │ Memory  │ │ Multi-Agent    │  │  │
│  │  │ (Aider) │ │ (Gemini) │ │(Cipher) │ │ (ClaudeSquad)  │  │  │
│  │  └─────────┘ └──────────┘ └─────────┘ └────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
│                          ▼                                       │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    PROVIDERS                               │  │
│  │  Claude | OpenAI | Gemini | Groq | Local | 50+ more       │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🧬 Core Features (da implementare)

### 1. 🗺️ RepoMap Intelligence (from Aider)
```go
type RepoMap struct {
    Files      []FileInfo
    Symbols    []Symbol      // functions, classes, types
    Imports    []ImportGraph // dependency graph
    Tests      []TestFile
    Docs       []DocFile
}

// Capisce il progetto PRIMA di agire
func (r *RepoMap) Analyze(path string) error
func (r *RepoMap) Summarize() string
func (r *RepoMap) FindRelevant(query string) []FileInfo
```

### 2. 🔌 MCP Native (from Gemini CLI)
```go
type MCPManager struct {
    Servers []MCPServer  // server MCP connessi
    Tools   []MCPTool    // tool disponibili
}

// Supporta sia server che client MCP
func (m *MCPManager) RegisterServer(config MCPConfig) error
func (m *MCPManager) CallTool(name string, args map[string]any) (any, error)
func (m *MCPManager) ListResources() []Resource
```

### 3. 🧠 Memory Layer (from Cipher)
```go
type Memory struct {
    Episodic  []Episode    // cosa è successo
    Semantic  []Concept    // cosa sa
    Working   []Context    // contesto attuale
}

// Memoria persistente cross-session
func (m *Memory) Remember(event Event) error
func (m *Memory) Recall(query string) []Memory
func (m *Memory) Forget(filter Filter) error
```

### 4. 🤖 Multi-Agent Orchestrator (from Claude Squad)
```go
type AgentOrchestrator struct {
    Agents    []Agent
    Sessions  []Session
    Scheduler Scheduler
}

// Gestisce più agenti in parallelo
func (o *AgentOrchestrator) Spawn(config AgentConfig) *Agent
func (o *AgentOrchestrator) Delegate(task Task, agent *Agent) error
func (o *AgentOrchestrator) Coordinate(tasks []Task) error
```

### 5. 📊 Real-time HUD (from Claude HUD)
```go
type HUD struct {
    ContextUsage  float64
    ActiveTools   []Tool
    RunningAgents []Agent
    TodoProgress  Progress
    TokenCount    int
}

// Visibilità in tempo reale
func (h *HUD) Render() string
func (h *HUD) Update(event Event)
```

---

## 🎯 Killer Features

### Auto-everything
```bash
# Zero config - capisce il progetto da solo
dream init  # analizza, configura, ready

# Auto-commit intelligente
dream commit  # genera messaggio, committa, push opzionale

# Auto-test
dream test  # genera test mancanti, esegue, report
```

### Voice & Browser
```bash
# Voice-to-code (Aider)
dream voice  # parla, codifica

# Browser control (Cline/GPTMe)
dream browse "apri docs e cerca X"
```

### Multi-Agent
```bash
# Spawna team di agenti
dream squad spawn --agents 3 --task "refactor module X"

# Monitor
dream squad status
dream squad logs agent-1
```

---

## 📁 Struttura Progetto

```
dream-cli/
├── cmd/
│   └── dream/
│       ├── main.go
│       ├── init.go
│       ├── chat.go
│       ├── squad.go
│       └── voice.go
├── internal/
│   ├── core/           # Core orchestration
│   ├── repomap/        # Project understanding (Aider)
│   ├── mcp/            # MCP server/client (Gemini)
│   ├── memory/         # Persistent memory (Cipher)
│   ├── agents/         # Multi-agent (Claude Squad)
│   ├── hud/            # Real-time display (Claude HUD)
│   ├── providers/      # AI providers (50+)
│   ├── tools/          # Built-in tools (20+)
│   └── tui/            # Terminal UI (OpenCode)
├── pkg/
│   ├── llm/            # Unified LLM interface
│   └── git/            # Git automation
└── configs/
    └── default.toml
```

---

## 🔧 CLI Commands

```bash
# Core
dream                    # Interactive TUI
dream chat "message"     # One-shot
dream init               # Initialize project

# RepoMap
dream map                # Show project map
dream map --symbols      # List all symbols
dream find "query"       # Find relevant files

# Memory
dream remember "fact"    # Store in memory
dream recall "query"     # Search memory
dream context            # Show current context

# Multi-Agent
dream squad spawn        # Create agent team
dream squad task "..."   # Assign task
dream squad status       # Monitor agents
dream squad merge        # Combine results

# Automation
dream commit             # Smart commit
dream test               # Generate & run tests
dream docs               # Generate documentation

# Voice
dream voice              # Voice input mode
dream listen             # Continuous listening

# Browser
dream browse "url"       # Control browser
dream scrape "selector"  # Extract data

# MCP
dream mcp list           # List MCP servers
dream mcp add "server"   # Add MCP server
dream mcp call "tool"    # Call MCP tool

# Config
dream config             # Edit config
dream providers          # List AI providers
dream provider set X     # Set default provider
```

---

## 🏗️ Implementation Priority

### Phase 1: Core (Week 1-2)
1. [ ] Go project scaffold
2. [ ] Basic TUI (bubbletea)
3. [ ] Single provider (Claude/OpenAI)
4. [ ] File read/write tools

### Phase 2: Intelligence (Week 3-4)
5. [ ] RepoMap implementation
6. [ ] Multi-file editing
7. [ ] Git integration
8. [ ] Auto-commit

### Phase 3: MCP & Memory (Week 5-6)
9. [ ] MCP client
10. [ ] MCP server
11. [ ] Cipher memory integration
12. [ ] Context management

### Phase 4: Multi-Agent (Week 7-8)
13. [ ] Agent spawning
14. [ ] Task delegation
15. [ ] Parallel execution
16. [ ] Result merging

### Phase 5: Advanced (Week 9-10)
17. [ ] Voice input
18. [ ] Browser control
19. [ ] Test generation
20. [ ] Desktop app

---

## 🎨 Nome Candidati

1. **DreamCode** - Il coding agent dei sogni
2. **OmniCode** - Onnicomprensivo
3. **Shogun** - Il comandante degli agenti
4. **Nexus** - Il punto di connessione
5. **Synthesis** - Sintesi di tutti i tool
6. **GoBro** - Il nome attuale del progetto! 🎯

---

*Spec v1.0 - Generated from 65 CLI tools analysis*
*Ready to implement in GoBro*
