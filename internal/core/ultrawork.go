// Package core implements the ultrawork pattern from oh-my-opencode
// "ultrawork" - the magic word that triggers total automation
package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/biodoia/goclit-ai/internal/agents"
)

// UltraWork is the magic command - total automation mode
// When invoked, it orchestrates all agents to complete a complex task
type UltraWork struct {
	mu          sync.RWMutex
	sisyphus    *agents.Sisyphus
	hephaestus  *agents.Hephaestus
	oracle      *agents.Oracle
	librarian   *agents.Librarian
	frontend    *agents.FrontendEngineer
	backend     *agents.BackendEngineer
	devops      *agents.DevOpsEngineer
	provider    agents.LLMProvider
	memory      agents.Memory
	status      string
	startTime   time.Time
	taskLog     []TaskLogEntry
}

type TaskLogEntry struct {
	Time      time.Time
	Agent     string
	Action    string
	Result    string
	Duration  time.Duration
}

// NewUltraWork creates the ultrawork orchestrator
func NewUltraWork(provider agents.LLMProvider, memory agents.Memory) *UltraWork {
	return &UltraWork{
		sisyphus:   agents.NewSisyphus(provider, agents.WithMemory(memory)),
		hephaestus: agents.NewHephaestus(provider, agents.WithHephaestusMemory(memory)),
		oracle:     agents.NewOracle(provider),
		librarian:  agents.NewLibrarian(provider),
		frontend:   agents.NewFrontendEngineer(provider, "react"),
		backend:    agents.NewBackendEngineer(provider, "go"),
		devops:     agents.NewDevOpsEngineer(provider),
		provider:   provider,
		memory:     memory,
		status:     "idle",
		taskLog:    make([]TaskLogEntry, 0),
	}
}

// Execute runs ultrawork on a task - full automation
func (u *UltraWork) Execute(ctx context.Context, task string) error {
	u.mu.Lock()
	u.status = "running"
	u.startTime = time.Now()
	u.mu.Unlock()

	defer func() {
		u.mu.Lock()
		u.status = "complete"
		u.mu.Unlock()
	}()

	// Phase 1: Oracle analyzes the task
	u.log("Oracle", "analyzing task", "")
	analysis, err := u.oracle.Ask(ctx, fmt.Sprintf(`Analyze this task and break it down into steps:

Task: %s

Provide:
1. Task understanding
2. Required steps (numbered)
3. Which specialized agents to use
4. Potential challenges
5. Success criteria`, task))
	if err != nil {
		return fmt.Errorf("oracle analysis failed: %w", err)
	}
	u.log("Oracle", "analysis complete", analysis[:min(200, len(analysis))])

	// Phase 2: Librarian gathers context
	u.log("Librarian", "gathering context", "")
	context, err := u.librarian.Search(ctx, "Find relevant code, docs, and examples for: "+task)
	if err != nil {
		// Non-fatal - continue without context
		context = "No additional context found"
	}
	u.log("Librarian", "context gathered", context[:min(200, len(context))])

	// Phase 3: Start Hephaestus for background work
	hephaestusCtx, cancelHephaestus := context.WithCancel(ctx)
	defer cancelHephaestus()

	go func() {
		u.hephaestus.Start(hephaestusCtx)
	}()

	// Phase 4: Sisyphus executes the main task
	u.log("Sisyphus", "starting main execution", "")

	// Build comprehensive prompt
	mainPrompt := fmt.Sprintf(`ULTRAWORK MODE ACTIVATED

Original Task: %s

Oracle's Analysis:
%s

Context from Librarian:
%s

You have access to specialized agents:
- Frontend Engineer (React)
- Backend Engineer (Go)
- DevOps Engineer
- Hephaestus (background builder)

Execute this task completely. Do not stop until done.
Delegate to specialists when appropriate.
Use DELEGATE:[agent] to assign subtasks.

Begin:`, task, analysis, context)

	// Configure Sisyphus with progress reporting
	progressCh := make(chan agents.Progress, 100)
	go func() {
		for p := range progressCh {
			u.log("Sisyphus", fmt.Sprintf("iteration %d", p.Iteration), p.Message)
		}
	}()

	sisyphus := agents.NewSisyphus(u.provider,
		agents.WithMemory(u.memory),
		agents.WithMaxRetries(50),
		agents.WithProgressCallback(func(p agents.Progress) {
			progressCh <- p
		}),
	)

	err = sisyphus.Work(ctx, mainPrompt)
	close(progressCh)

	if err != nil {
		return fmt.Errorf("sisyphus execution failed: %w", err)
	}

	u.log("UltraWork", "task complete", fmt.Sprintf("Duration: %v", time.Since(u.startTime)))
	return nil
}

// Status returns current status
func (u *UltraWork) Status() string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.status
}

// TaskLog returns the task log
func (u *UltraWork) TaskLog() []TaskLogEntry {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.taskLog
}

// Duration returns elapsed time
func (u *UltraWork) Duration() time.Duration {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return time.Since(u.startTime)
}

func (u *UltraWork) log(agent, action, result string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	entry := TaskLogEntry{
		Time:   time.Now(),
		Agent:  agent,
		Action: action,
		Result: result,
	}

	if len(u.taskLog) > 0 {
		entry.Duration = time.Since(u.taskLog[len(u.taskLog)-1].Time)
	}

	u.taskLog = append(u.taskLog, entry)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// UltraWorkCommand is the CLI entry point for "goclit ultrawork"
func UltraWorkCommand(ctx context.Context, task string, provider agents.LLMProvider, memory agents.Memory) error {
	fmt.Println("🚀 ULTRAWORK MODE ACTIVATED")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Task: %s\n\n", task)

	uw := NewUltraWork(provider, memory)

	// Progress display
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				status := uw.Status()
				duration := uw.Duration()
				fmt.Printf("⏳ Status: %s | Duration: %v\n", status, duration.Round(time.Second))
			}
		}
	}()

	err := uw.Execute(ctx, task)

	// Print final log
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 TASK LOG:")
	for _, entry := range uw.TaskLog() {
		fmt.Printf("  [%s] %s: %s\n", entry.Agent, entry.Action, entry.Result)
	}

	if err != nil {
		fmt.Printf("\n❌ ULTRAWORK FAILED: %v\n", err)
		return err
	}

	fmt.Println("\n✅ ULTRAWORK COMPLETE")
	return nil
}
