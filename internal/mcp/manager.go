// Package mcp implements Model Context Protocol support
// Pattern from oh-my-opencode: curated MCPs
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Manager handles MCP server connections
type Manager struct {
	mu      sync.RWMutex
	servers map[string]*Server
	tools   map[string]Tool
}

// Server represents an MCP server connection
type Server struct {
	Name        string
	URL         string
	Transport   string // "stdio" | "http" | "ws"
	Status      string
	Tools       []Tool
	Resources   []Resource
	Prompts     []Prompt
}

// Tool is an MCP tool definition
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
	ServerName  string          `json:"serverName,omitempty"`
}

// Resource is an MCP resource
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// Prompt is an MCP prompt template
type Prompt struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Arguments   []Argument `json:"arguments,omitempty"`
}

// Argument is a prompt argument
type Argument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// NewManager creates a new MCP manager
func NewManager() *Manager {
	return &Manager{
		servers: make(map[string]*Server),
		tools:   make(map[string]Tool),
	}
}

// RegisterServer adds an MCP server
func (m *Manager) RegisterServer(server *Server) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.servers[server.Name] = server

	// Index tools by name
	for _, tool := range server.Tools {
		tool.ServerName = server.Name
		m.tools[tool.Name] = tool
	}

	return nil
}

// UnregisterServer removes an MCP server
func (m *Manager) UnregisterServer(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if server, ok := m.servers[name]; ok {
		// Remove tools
		for _, tool := range server.Tools {
			delete(m.tools, tool.Name)
		}
		delete(m.servers, name)
	}
}

// ListServers returns all registered servers
func (m *Manager) ListServers() []*Server {
	m.mu.RLock()
	defer m.mu.RUnlock()

	servers := make([]*Server, 0, len(m.servers))
	for _, s := range m.servers {
		servers = append(servers, s)
	}
	return servers
}

// ListTools returns all available tools
func (m *Manager) ListTools() []Tool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tools := make([]Tool, 0, len(m.tools))
	for _, t := range m.tools {
		tools = append(tools, t)
	}
	return tools
}

// CallTool invokes an MCP tool
func (m *Manager) CallTool(ctx context.Context, name string, args map[string]any) (any, error) {
	m.mu.RLock()
	tool, ok := m.tools[name]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	m.mu.RLock()
	server := m.servers[tool.ServerName]
	m.mu.RUnlock()

	if server == nil {
		return nil, fmt.Errorf("server not found for tool: %s", name)
	}

	// TODO: Implement actual MCP protocol call
	return nil, fmt.Errorf("MCP protocol call not implemented yet")
}

// CuratedMCPs returns pre-configured MCP servers
// Pattern from oh-my-opencode
func CuratedMCPs() []*Server {
	return []*Server{
		{
			Name:      "filesystem",
			Transport: "stdio",
			Tools: []Tool{
				{Name: "read_file", Description: "Read a file's contents"},
				{Name: "write_file", Description: "Write content to a file"},
				{Name: "list_directory", Description: "List directory contents"},
				{Name: "create_directory", Description: "Create a new directory"},
				{Name: "delete_file", Description: "Delete a file"},
				{Name: "move_file", Description: "Move/rename a file"},
			},
		},
		{
			Name:      "shell",
			Transport: "stdio",
			Tools: []Tool{
				{Name: "execute", Description: "Execute a shell command"},
				{Name: "background", Description: "Run command in background"},
			},
		},
		{
			Name:      "git",
			Transport: "stdio",
			Tools: []Tool{
				{Name: "status", Description: "Get git status"},
				{Name: "diff", Description: "Get git diff"},
				{Name: "commit", Description: "Create a commit"},
				{Name: "push", Description: "Push changes"},
				{Name: "pull", Description: "Pull changes"},
				{Name: "log", Description: "Get commit history"},
			},
		},
		{
			Name:      "search",
			Transport: "stdio",
			Tools: []Tool{
				{Name: "grep", Description: "Search for pattern in files"},
				{Name: "find", Description: "Find files by name"},
				{Name: "ripgrep", Description: "Fast regex search"},
			},
		},
		{
			Name:      "browser",
			Transport: "http",
			Tools: []Tool{
				{Name: "navigate", Description: "Navigate to URL"},
				{Name: "screenshot", Description: "Take screenshot"},
				{Name: "click", Description: "Click element"},
				{Name: "type", Description: "Type text"},
				{Name: "extract", Description: "Extract page content"},
			},
		},
		{
			Name:      "database",
			Transport: "stdio",
			Tools: []Tool{
				{Name: "query", Description: "Execute SQL query"},
				{Name: "schema", Description: "Get database schema"},
			},
		},
	}
}

// SetupCuratedMCPs registers all curated MCP servers
func (m *Manager) SetupCuratedMCPs() {
	for _, server := range CuratedMCPs() {
		m.RegisterServer(server)
	}
}
