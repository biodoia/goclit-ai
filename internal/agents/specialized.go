// Specialized agents - domain-specific AI agents
// Pattern from oh-my-opencode: oracle, librarian, frontend engineer
package agents

import (
	"context"
	"fmt"
)

// Oracle - The knowledge and decision agent
// Answers questions, provides insights, makes recommendations
type Oracle struct {
	BaseAgent
	knowledgeBase []string
}

func NewOracle(provider LLMProvider) *Oracle {
	return &Oracle{
		BaseAgent: BaseAgent{
			name:     "Oracle",
			role:     "Knowledge & Decision Agent",
			provider: provider,
			systemPrompt: `You are Oracle, the knowledge agent.
You answer questions with deep insight and wisdom.
You make recommendations based on analysis.
You connect disparate pieces of information.

Your strengths:
- Deep knowledge of programming patterns
- Understanding of best practices
- Ability to see the big picture
- Making informed recommendations

When asked a question:
1. Consider all angles
2. Provide context
3. Give a clear answer
4. Suggest related insights`,
		},
	}
}

// Ask queries the Oracle
func (o *Oracle) Ask(ctx context.Context, question string) (string, error) {
	prompt := o.systemPrompt + "\n\nQuestion: " + question
	return o.provider.Generate(ctx, prompt)
}

// Recommend provides recommendations
func (o *Oracle) Recommend(ctx context.Context, situation string) (string, error) {
	prompt := o.systemPrompt + "\n\nSituation requiring recommendation:\n" + situation + "\n\nProvide your recommendation:"
	return o.provider.Generate(ctx, prompt)
}

// Librarian - The documentation and search agent
// Finds information, organizes knowledge, maintains docs
type Librarian struct {
	BaseAgent
	indices map[string][]string
}

func NewLibrarian(provider LLMProvider) *Librarian {
	return &Librarian{
		BaseAgent: BaseAgent{
			name:     "Librarian",
			role:     "Documentation & Search Agent",
			provider: provider,
			systemPrompt: `You are Librarian, the documentation agent.
You find information, organize knowledge, and maintain documentation.

Your strengths:
- Finding relevant code and documentation
- Organizing information logically
- Creating clear documentation
- Summarizing complex topics

When searching:
1. Identify key terms
2. Look in relevant places
3. Return ranked results
4. Suggest related resources`,
		},
		indices: make(map[string][]string),
	}
}

// Search finds information
func (l *Librarian) Search(ctx context.Context, query string) (string, error) {
	prompt := l.systemPrompt + "\n\nSearch query: " + query + "\n\nFind and return relevant information:"
	return l.provider.Generate(ctx, prompt)
}

// Document creates documentation
func (l *Librarian) Document(ctx context.Context, code string) (string, error) {
	prompt := l.systemPrompt + "\n\nCode to document:\n```\n" + code + "\n```\n\nCreate clear documentation:"
	return l.provider.Generate(ctx, prompt)
}

// Summarize condenses information
func (l *Librarian) Summarize(ctx context.Context, content string) (string, error) {
	prompt := l.systemPrompt + "\n\nContent to summarize:\n" + content + "\n\nProvide a clear summary:"
	return l.provider.Generate(ctx, prompt)
}

// FrontendEngineer - The UI/UX specialist agent
// Builds interfaces, handles styling, creates components
type FrontendEngineer struct {
	BaseAgent
	framework string // react, vue, svelte, etc.
}

func NewFrontendEngineer(provider LLMProvider, framework string) *FrontendEngineer {
	return &FrontendEngineer{
		BaseAgent: BaseAgent{
			name:     "Frontend Engineer",
			role:     "UI/UX Specialist Agent",
			provider: provider,
			systemPrompt: fmt.Sprintf(`You are Frontend Engineer, the UI/UX specialist.
You build beautiful, functional user interfaces.
Your framework of choice: %s

Your strengths:
- Creating responsive layouts
- Building reusable components
- Handling state management
- Writing clean CSS/styling
- Ensuring accessibility

When building UI:
1. Understand the requirements
2. Plan the component structure
3. Write clean, maintainable code
4. Consider edge cases
5. Ensure good UX`, framework),
		},
		framework: framework,
	}
}

// BuildComponent creates a UI component
func (f *FrontendEngineer) BuildComponent(ctx context.Context, spec string) (string, error) {
	prompt := f.systemPrompt + "\n\nComponent specification:\n" + spec + "\n\nBuild the component:"
	return f.provider.Generate(ctx, prompt)
}

// StyleComponent creates styling
func (f *FrontendEngineer) StyleComponent(ctx context.Context, component string) (string, error) {
	prompt := f.systemPrompt + "\n\nComponent to style:\n" + component + "\n\nCreate beautiful styling:"
	return f.provider.Generate(ctx, prompt)
}

// BackendEngineer - The server-side specialist
type BackendEngineer struct {
	BaseAgent
	language string // go, python, rust, etc.
}

func NewBackendEngineer(provider LLMProvider, language string) *BackendEngineer {
	return &BackendEngineer{
		BaseAgent: BaseAgent{
			name:     "Backend Engineer",
			role:     "Server-side Specialist Agent",
			provider: provider,
			systemPrompt: fmt.Sprintf(`You are Backend Engineer, the server-side specialist.
You build robust, scalable backend systems.
Your language of choice: %s

Your strengths:
- Designing APIs
- Database optimization
- Security best practices
- Performance tuning
- Error handling

When building backend:
1. Design the API contract
2. Plan data models
3. Implement with security in mind
4. Handle errors gracefully
5. Document endpoints`, language),
		},
		language: language,
	}
}

// BuildAPI creates an API endpoint
func (b *BackendEngineer) BuildAPI(ctx context.Context, spec string) (string, error) {
	prompt := b.systemPrompt + "\n\nAPI specification:\n" + spec + "\n\nImplement the API:"
	return b.provider.Generate(ctx, prompt)
}

// DevOpsEngineer - The infrastructure specialist
type DevOpsEngineer struct {
	BaseAgent
}

func NewDevOpsEngineer(provider LLMProvider) *DevOpsEngineer {
	return &DevOpsEngineer{
		BaseAgent: BaseAgent{
			name:     "DevOps Engineer",
			role:     "Infrastructure Specialist Agent",
			provider: provider,
			systemPrompt: `You are DevOps Engineer, the infrastructure specialist.
You build and maintain deployment pipelines and infrastructure.

Your strengths:
- CI/CD pipelines
- Container orchestration
- Infrastructure as code
- Monitoring and logging
- Security hardening

When building infrastructure:
1. Plan the architecture
2. Write infrastructure as code
3. Set up CI/CD
4. Configure monitoring
5. Document operations`,
		},
	}
}

// BuildPipeline creates a CI/CD pipeline
func (d *DevOpsEngineer) BuildPipeline(ctx context.Context, spec string) (string, error) {
	prompt := d.systemPrompt + "\n\nPipeline specification:\n" + spec + "\n\nCreate the pipeline:"
	return d.provider.Generate(ctx, prompt)
}

// BaseAgent provides common functionality
type BaseAgent struct {
	name         string
	role         string
	provider     LLMProvider
	systemPrompt string
	tools        []Tool
	memory       Memory
}

func (b *BaseAgent) Name() string { return b.name }
func (b *BaseAgent) Role() string { return b.role }
