// Package agents implements specialized AI agents
// Pattern from oh-my-opencode: agents that work until the task is done
package agents

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Sisyphus is the discipline agent - works until the task is done
// "It just works until the task is done. It is a discipline agent."
type Sisyphus struct {
	mu          sync.RWMutex
	name        string
	task        string
	status      AgentStatus
	iterations  int
	maxRetries  int
	provider    LLMProvider
	tools       []Tool
	memory      Memory
	onProgress  func(Progress)
	startTime   time.Time
}

type AgentStatus string

const (
	StatusIdle     AgentStatus = "idle"
	StatusRunning  AgentStatus = "running"
	StatusPaused   AgentStatus = "paused"
	StatusComplete AgentStatus = "complete"
	StatusFailed   AgentStatus = "failed"
)

type Progress struct {
	Iteration   int
	TotalSteps  int
	CurrentStep string
	Percentage  float64
	Message     string
}

type LLMProvider interface {
	Generate(ctx context.Context, prompt string) (string, error)
	GenerateWithTools(ctx context.Context, prompt string, tools []Tool) (string, error)
}

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, args map[string]any) (any, error)
}

type Memory interface {
	Store(key string, value any) error
	Recall(key string) (any, error)
	Context() string
}

// NewSisyphus creates a new Sisyphus agent
func NewSisyphus(provider LLMProvider, opts ...SisyphusOption) *Sisyphus {
	s := &Sisyphus{
		name:       "Sisyphus",
		status:     StatusIdle,
		maxRetries: 100, // Never give up easily
		provider:   provider,
		tools:      make([]Tool, 0),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type SisyphusOption func(*Sisyphus)

func WithMaxRetries(n int) SisyphusOption {
	return func(s *Sisyphus) { s.maxRetries = n }
}

func WithTools(tools ...Tool) SisyphusOption {
	return func(s *Sisyphus) { s.tools = append(s.tools, tools...) }
}

func WithMemory(m Memory) SisyphusOption {
	return func(s *Sisyphus) { s.memory = m }
}

func WithProgressCallback(fn func(Progress)) SisyphusOption {
	return func(s *Sisyphus) { s.onProgress = fn }
}

// Work starts the Sisyphus agent on a task - runs until complete
func (s *Sisyphus) Work(ctx context.Context, task string) error {
	s.mu.Lock()
	s.task = task
	s.status = StatusRunning
	s.iterations = 0
	s.startTime = time.Now()
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		if s.status == StatusRunning {
			s.status = StatusComplete
		}
		s.mu.Unlock()
	}()

	systemPrompt := `You are Sisyphus, the discipline agent.
Your job is to work on a task until it is COMPLETELY done.
Never give up. Never stop halfway. Push through obstacles.

Rules:
1. Break the task into steps
2. Execute each step methodically
3. Verify completion before moving on
4. If something fails, try a different approach
5. Report progress clearly
6. Only stop when the task is 100% complete

Current task: ` + task

	for s.iterations < s.maxRetries {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			s.status = StatusPaused
			s.mu.Unlock()
			return ctx.Err()
		default:
		}

		s.iterations++

		// Build prompt with memory context
		prompt := systemPrompt
		if s.memory != nil {
			prompt += "\n\nContext from memory:\n" + s.memory.Context()
		}

		// Generate next action
		response, err := s.provider.GenerateWithTools(ctx, prompt, s.tools)
		if err != nil {
			continue // Sisyphus doesn't give up on errors
		}

		// Report progress
		if s.onProgress != nil {
			s.onProgress(Progress{
				Iteration:   s.iterations,
				CurrentStep: response[:min(100, len(response))],
				Message:     fmt.Sprintf("Iteration %d", s.iterations),
			})
		}

		// Check if task is complete
		if isTaskComplete(response) {
			return nil
		}

		// Store in memory
		if s.memory != nil {
			s.memory.Store(fmt.Sprintf("iteration_%d", s.iterations), response)
		}
	}

	s.mu.Lock()
	s.status = StatusFailed
	s.mu.Unlock()
	return fmt.Errorf("max retries (%d) exceeded", s.maxRetries)
}

// Status returns the current agent status
func (s *Sisyphus) Status() AgentStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// Iterations returns the number of iterations completed
func (s *Sisyphus) Iterations() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.iterations
}

// Pause pauses the agent
func (s *Sisyphus) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.status == StatusRunning {
		s.status = StatusPaused
	}
}

// Resume resumes a paused agent
func (s *Sisyphus) Resume(ctx context.Context) error {
	s.mu.Lock()
	if s.status != StatusPaused {
		s.mu.Unlock()
		return fmt.Errorf("agent not paused")
	}
	s.status = StatusRunning
	s.mu.Unlock()

	return s.Work(ctx, s.task)
}

func isTaskComplete(response string) bool {
	// Check for completion markers
	markers := []string{
		"TASK_COMPLETE",
		"task is complete",
		"successfully completed",
		"all done",
		"finished",
	}
	for _, marker := range markers {
		if contains(response, marker) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsLower(s, substr))
}

func containsLower(s, substr string) bool {
	// Simple contains check
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
