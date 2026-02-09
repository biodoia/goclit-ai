// Hephaestus - The Autonomy Agent
// God of the forge - creates and builds autonomously
package agents

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Hephaestus is the autonomy agent - builds and creates independently
// Named after the Greek god of the forge and craftsmanship
type Hephaestus struct {
	mu           sync.RWMutex
	name         string
	status       AgentStatus
	provider     LLMProvider
	tools        []Tool
	memory       Memory
	workQueue    chan Task
	results      chan Result
	workers      int
	onArtifact   func(Artifact)
}

type Task struct {
	ID          string
	Description string
	Priority    int
	Dependencies []string
	CreatedAt   time.Time
}

type Result struct {
	TaskID    string
	Success   bool
	Output    any
	Error     error
	Duration  time.Duration
}

type Artifact struct {
	Type     string // "file", "code", "config", "doc"
	Path     string
	Content  []byte
	Metadata map[string]string
}

// NewHephaestus creates the autonomy agent
func NewHephaestus(provider LLMProvider, opts ...HephaestusOption) *Hephaestus {
	h := &Hephaestus{
		name:      "Hephaestus",
		status:    StatusIdle,
		provider:  provider,
		tools:     make([]Tool, 0),
		workQueue: make(chan Task, 100),
		results:   make(chan Result, 100),
		workers:   3, // Parallel workers
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

type HephaestusOption func(*Hephaestus)

func WithWorkers(n int) HephaestusOption {
	return func(h *Hephaestus) { h.workers = n }
}

func WithHephaestusTools(tools ...Tool) HephaestusOption {
	return func(h *Hephaestus) { h.tools = append(h.tools, tools...) }
}

func WithHephaestusMemory(m Memory) HephaestusOption {
	return func(h *Hephaestus) { h.memory = m }
}

func WithArtifactCallback(fn func(Artifact)) HephaestusOption {
	return func(h *Hephaestus) { h.onArtifact = fn }
}

// Start begins the autonomous work loop
func (h *Hephaestus) Start(ctx context.Context) error {
	h.mu.Lock()
	h.status = StatusRunning
	h.mu.Unlock()

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < h.workers; i++ {
		wg.Add(1)
		go h.worker(ctx, &wg, i)
	}

	// Wait for context cancellation
	<-ctx.Done()

	h.mu.Lock()
	h.status = StatusPaused
	h.mu.Unlock()

	close(h.workQueue)
	wg.Wait()

	return ctx.Err()
}

// worker processes tasks from the queue
func (h *Hephaestus) worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-h.workQueue:
			if !ok {
				return
			}
			result := h.processTask(ctx, task)
			h.results <- result
		}
	}
}

// processTask handles a single task
func (h *Hephaestus) processTask(ctx context.Context, task Task) Result {
	start := time.Now()

	prompt := fmt.Sprintf(`You are Hephaestus, the autonomy agent.
You are the god of the forge - you CREATE and BUILD things independently.

Your task: %s

Rules:
1. Analyze the task thoroughly
2. Plan your approach
3. Execute step by step
4. Create high-quality artifacts
5. Verify your work
6. Report what you created

Output artifacts using:
ARTIFACT_START
type: [file|code|config|doc]
path: [path/to/artifact]
---
[content]
ARTIFACT_END

Begin:`, task.Description)

	response, err := h.provider.GenerateWithTools(ctx, prompt, h.tools)
	if err != nil {
		return Result{
			TaskID:   task.ID,
			Success:  false,
			Error:    err,
			Duration: time.Since(start),
		}
	}

	// Parse and emit artifacts
	artifacts := parseArtifacts(response)
	for _, artifact := range artifacts {
		if h.onArtifact != nil {
			h.onArtifact(artifact)
		}
	}

	return Result{
		TaskID:   task.ID,
		Success:  true,
		Output:   response,
		Duration: time.Since(start),
	}
}

// Submit adds a task to the work queue
func (h *Hephaestus) Submit(task Task) error {
	h.mu.RLock()
	status := h.status
	h.mu.RUnlock()

	if status != StatusRunning {
		return fmt.Errorf("agent not running")
	}

	select {
	case h.workQueue <- task:
		return nil
	default:
		return fmt.Errorf("work queue full")
	}
}

// Results returns the results channel
func (h *Hephaestus) Results() <-chan Result {
	return h.results
}

// Status returns current status
func (h *Hephaestus) Status() AgentStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status
}

// parseArtifacts extracts artifacts from response
func parseArtifacts(response string) []Artifact {
	// Simple parser - look for ARTIFACT_START/END blocks
	artifacts := make([]Artifact, 0)
	// TODO: Implement full parser
	return artifacts
}
