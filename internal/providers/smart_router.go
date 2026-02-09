// Package providers implements intelligent multi-provider routing
// Auto-discovers accounts, pools quotas, smart prioritization
package providers

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// SmartRouter manages all providers with intelligent routing
type SmartRouter struct {
	mu            sync.RWMutex
	providers     map[string]*Provider
	models        map[string][]*ModelAvailability // model -> providers that have it
	quotaPool     map[string]*QuotaPool           // model -> pooled quota
	rareModels    map[string]string               // model -> exclusive provider
	usageHistory  []UsageRecord
	config        RouterConfig
}

// Provider represents an AI provider account
type Provider struct {
	Name         string
	Type         ProviderType
	APIKey       string
	BaseURL      string
	Models       []Model
	Quota        Quota
	Status       ProviderStatus
	LastUsed     time.Time
	ErrorCount   int
	Latency      time.Duration // avg response time
}

type ProviderType string

const (
	ProviderOpenAI      ProviderType = "openai"
	ProviderAnthropic   ProviderType = "anthropic"
	ProviderGoogle      ProviderType = "google"
	ProviderGroq        ProviderType = "groq"
	ProviderOpenRouter  ProviderType = "openrouter"
	ProviderTogether    ProviderType = "together"
	ProviderMistral     ProviderType = "mistral"
	ProviderDeepSeek    ProviderType = "deepseek"
	ProviderCerebras    ProviderType = "cerebras"
	ProviderSambaNova   ProviderType = "sambanova"
	ProviderFireworks   ProviderType = "fireworks"
	ProviderPerplexity  ProviderType = "perplexity"
	ProviderCohere      ProviderType = "cohere"
	ProviderLocal       ProviderType = "local"
)

type ProviderStatus string

const (
	StatusActive    ProviderStatus = "active"
	StatusRateLimited ProviderStatus = "rate_limited"
	StatusQuotaExhausted ProviderStatus = "quota_exhausted"
	StatusError     ProviderStatus = "error"
	StatusDisabled  ProviderStatus = "disabled"
)

// Model represents an AI model
type Model struct {
	ID          string
	Name        string
	Provider    string
	Type        ModelType // llm, stt, tts, embedding, image
	ContextSize int
	InputCost   float64  // per 1M tokens
	OutputCost  float64
	IsRare      bool     // only available on this provider
}

type ModelType string

const (
	ModelLLM       ModelType = "llm"
	ModelSTT       ModelType = "stt"
	ModelTTS       ModelType = "tts"
	ModelEmbedding ModelType = "embedding"
	ModelImage     ModelType = "image"
)

// Quota represents provider quota/limits
type Quota struct {
	RequestsPerMinute int
	RequestsPerDay    int
	TokensPerMinute   int
	TokensPerDay      int
	UsedRequests      int
	UsedTokens        int
	ResetTime         time.Time
}

// QuotaPool pools quota across providers for same model
type QuotaPool struct {
	Model            string
	TotalRequests    int
	TotalTokens      int
	AvailableRequests int
	AvailableTokens  int
	Providers        []string // providers contributing to pool
}

// ModelAvailability tracks where a model is available
type ModelAvailability struct {
	Provider    *Provider
	Model       Model
	IsExclusive bool // only this provider has it
	Priority    int  // lower = use first
}

// UsageRecord tracks usage for smart routing
type UsageRecord struct {
	Time      time.Time
	Provider  string
	Model     string
	Tokens    int
	Latency   time.Duration
	Success   bool
}

type RouterConfig struct {
	PreferFree          bool    // prefer free tiers
	PreferFast          bool    // prefer low latency
	PreferCheap         bool    // prefer low cost
	MaxLatencyMs        int     // max acceptable latency
	FallbackEnabled     bool    // try next provider on failure
	QuotaReservePercent float64 // reserve % of quota for rare models
}

// NewSmartRouter creates a new intelligent router
func NewSmartRouter(config RouterConfig) *SmartRouter {
	return &SmartRouter{
		providers:    make(map[string]*Provider),
		models:       make(map[string][]*ModelAvailability),
		quotaPool:    make(map[string]*QuotaPool),
		rareModels:   make(map[string]string),
		usageHistory: make([]UsageRecord, 0),
		config:       config,
	}
}

// AutoDiscover finds all configured provider accounts
func (r *SmartRouter) AutoDiscover() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	discoveries := []struct {
		envKey   string
		provider ProviderType
		discover func(string) (*Provider, error)
	}{
		{"OPENAI_API_KEY", ProviderOpenAI, discoverOpenAI},
		{"ANTHROPIC_API_KEY", ProviderAnthropic, discoverAnthropic},
		{"GOOGLE_API_KEY", ProviderGoogle, discoverGoogle},
		{"GROQ_API_KEY", ProviderGroq, discoverGroq},
		{"OPENROUTER_API_KEY", ProviderOpenRouter, discoverOpenRouter},
		{"TOGETHER_API_KEY", ProviderTogether, discoverTogether},
		{"MISTRAL_API_KEY", ProviderMistral, discoverMistral},
		{"DEEPSEEK_API_KEY", ProviderDeepSeek, discoverDeepSeek},
		{"CEREBRAS_API_KEY", ProviderCerebras, discoverCerebras},
		{"SAMBANOVA_API_KEY", ProviderSambaNova, discoverSambaNova},
		{"FIREWORKS_API_KEY", ProviderFireworks, discoverFireworks},
		{"PERPLEXITY_API_KEY", ProviderPerplexity, discoverPerplexity},
		{"COHERE_API_KEY", ProviderCohere, discoverCohere},
	}

	for _, d := range discoveries {
		key := getEnv(d.envKey)
		if key != "" {
			provider, err := d.discover(key)
			if err == nil && provider != nil {
				r.providers[provider.Name] = provider
				r.indexModels(provider)
			}
		}
	}

	// Discover local models (Ollama, etc)
	if local := discoverLocal(); local != nil {
		r.providers[local.Name] = local
		r.indexModels(local)
	}

	// Identify rare models (only one provider has them)
	r.identifyRareModels()
	
	// Build quota pools
	r.buildQuotaPools()

	return nil
}

// indexModels adds provider's models to the index
func (r *SmartRouter) indexModels(provider *Provider) {
	for _, model := range provider.Models {
		avail := &ModelAvailability{
			Provider: provider,
			Model:    model,
		}
		r.models[model.ID] = append(r.models[model.ID], avail)
	}
}

// identifyRareModels finds models only available on one provider
func (r *SmartRouter) identifyRareModels() {
	for modelID, availabilities := range r.models {
		if len(availabilities) == 1 {
			r.rareModels[modelID] = availabilities[0].Provider.Name
			availabilities[0].IsExclusive = true
			availabilities[0].Model.IsRare = true
		}
	}
}

// buildQuotaPools creates pooled quotas for shared models
func (r *SmartRouter) buildQuotaPools() {
	for modelID, availabilities := range r.models {
		if len(availabilities) > 1 {
			pool := &QuotaPool{
				Model:     modelID,
				Providers: make([]string, 0),
			}
			for _, avail := range availabilities {
				pool.TotalRequests += avail.Provider.Quota.RequestsPerDay
				pool.TotalTokens += avail.Provider.Quota.TokensPerDay
				pool.Providers = append(pool.Providers, avail.Provider.Name)
			}
			pool.AvailableRequests = pool.TotalRequests
			pool.AvailableTokens = pool.TotalTokens
			r.quotaPool[modelID] = pool
		}
	}
}

// Route selects the best provider for a model request
func (r *SmartRouter) Route(ctx context.Context, modelID string, tokenEstimate int) (*Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	availabilities, ok := r.models[modelID]
	if !ok || len(availabilities) == 0 {
		return nil, fmt.Errorf("model not found: %s", modelID)
	}

	// Calculate priority for each availability
	ranked := r.rankProviders(availabilities, tokenEstimate)
	
	// Return best available
	for _, avail := range ranked {
		if r.isProviderAvailable(avail.Provider, tokenEstimate) {
			return avail.Provider, nil
		}
	}

	return nil, fmt.Errorf("no available provider for model: %s", modelID)
}

// rankProviders sorts providers by smart priority
func (r *SmartRouter) rankProviders(availabilities []*ModelAvailability, tokenEstimate int) []*ModelAvailability {
	// Copy to avoid modifying original
	ranked := make([]*ModelAvailability, len(availabilities))
	copy(ranked, availabilities)

	sort.Slice(ranked, func(i, j int) bool {
		pi, pj := ranked[i], ranked[j]

		// Rule 1: Preserve providers with rare models
		// If provider has rare models, deprioritize for common model requests
		iHasRare := r.providerHasRareModels(pi.Provider.Name)
		jHasRare := r.providerHasRareModels(pj.Provider.Name)
		
		if iHasRare && !jHasRare && !pi.IsExclusive {
			return false // j first (doesn't have rare models)
		}
		if jHasRare && !iHasRare && !pj.IsExclusive {
			return true // i first
		}

		// Rule 2: If this IS the rare model, use its exclusive provider
		if pi.IsExclusive {
			return true
		}
		if pj.IsExclusive {
			return false
		}

		// Rule 3: Prefer providers with more remaining quota
		iQuotaRatio := r.quotaRemaining(pi.Provider)
		jQuotaRatio := r.quotaRemaining(pj.Provider)
		if iQuotaRatio != jQuotaRatio {
			return iQuotaRatio > jQuotaRatio
		}

		// Rule 4: Prefer lower latency
		if r.config.PreferFast {
			if pi.Provider.Latency != pj.Provider.Latency {
				return pi.Provider.Latency < pj.Provider.Latency
			}
		}

		// Rule 5: Prefer lower cost
		if r.config.PreferCheap {
			iCost := pi.Model.InputCost + pi.Model.OutputCost
			jCost := pj.Model.InputCost + pj.Model.OutputCost
			if iCost != jCost {
				return iCost < jCost
			}
		}

		// Rule 6: Prefer less recently used (load balancing)
		return pi.Provider.LastUsed.Before(pj.Provider.LastUsed)
	})

	return ranked
}

// providerHasRareModels checks if provider has any exclusive models
func (r *SmartRouter) providerHasRareModels(providerName string) bool {
	for _, exclusiveProvider := range r.rareModels {
		if exclusiveProvider == providerName {
			return true
		}
	}
	return false
}

// quotaRemaining returns the fraction of quota remaining
func (r *SmartRouter) quotaRemaining(p *Provider) float64 {
	if p.Quota.RequestsPerDay == 0 {
		return 1.0 // unlimited
	}
	return float64(p.Quota.RequestsPerDay-p.Quota.UsedRequests) / float64(p.Quota.RequestsPerDay)
}

// isProviderAvailable checks if provider can handle request
func (r *SmartRouter) isProviderAvailable(p *Provider, tokenEstimate int) bool {
	if p.Status != StatusActive {
		return false
	}
	if p.Quota.UsedRequests >= p.Quota.RequestsPerDay {
		return false
	}
	if p.Quota.TokensPerDay > 0 && p.Quota.UsedTokens+tokenEstimate > p.Quota.TokensPerDay {
		return false
	}
	// Reserve quota for rare models
	if r.providerHasRareModels(p.Name) {
		reserve := int(float64(p.Quota.RequestsPerDay) * r.config.QuotaReservePercent)
		if p.Quota.UsedRequests >= p.Quota.RequestsPerDay-reserve {
			return false
		}
	}
	return true
}

// RecordUsage tracks usage for smart routing
func (r *SmartRouter) RecordUsage(provider, model string, tokens int, latency time.Duration, success bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Update provider stats
	if p, ok := r.providers[provider]; ok {
		p.LastUsed = time.Now()
		p.Quota.UsedRequests++
		p.Quota.UsedTokens += tokens
		
		// Update latency (rolling average)
		if p.Latency == 0 {
			p.Latency = latency
		} else {
			p.Latency = (p.Latency + latency) / 2
		}
		
		if !success {
			p.ErrorCount++
			if p.ErrorCount > 3 {
				p.Status = StatusError
			}
		}
	}

	// Record for history
	r.usageHistory = append(r.usageHistory, UsageRecord{
		Time:     time.Now(),
		Provider: provider,
		Model:    model,
		Tokens:   tokens,
		Latency:  latency,
		Success:  success,
	})

	// Keep last 1000 records
	if len(r.usageHistory) > 1000 {
		r.usageHistory = r.usageHistory[100:]
	}
}

// GetStats returns router statistics
func (r *SmartRouter) GetStats() RouterStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := RouterStats{
		TotalProviders:   len(r.providers),
		TotalModels:      len(r.models),
		RareModels:       len(r.rareModels),
		ProviderStats:    make(map[string]ProviderStats),
	}

	for name, p := range r.providers {
		stats.ProviderStats[name] = ProviderStats{
			Status:        string(p.Status),
			QuotaUsed:     p.Quota.UsedRequests,
			QuotaTotal:    p.Quota.RequestsPerDay,
			AvgLatency:    p.Latency,
			ErrorCount:    p.ErrorCount,
			ModelCount:    len(p.Models),
			HasRareModels: r.providerHasRareModels(name),
		}
	}

	return stats
}

type RouterStats struct {
	TotalProviders int
	TotalModels    int
	RareModels     int
	ProviderStats  map[string]ProviderStats
}

type ProviderStats struct {
	Status        string
	QuotaUsed     int
	QuotaTotal    int
	AvgLatency    time.Duration
	ErrorCount    int
	ModelCount    int
	HasRareModels bool
}

// Helper functions (stubs - implement per provider)
func getEnv(key string) string {
	// TODO: Read from env
	return ""
}

func discoverOpenAI(key string) (*Provider, error)      { return nil, nil }
func discoverAnthropic(key string) (*Provider, error)   { return nil, nil }
func discoverGoogle(key string) (*Provider, error)      { return nil, nil }
func discoverGroq(key string) (*Provider, error)        { return nil, nil }
func discoverOpenRouter(key string) (*Provider, error)  { return nil, nil }
func discoverTogether(key string) (*Provider, error)    { return nil, nil }
func discoverMistral(key string) (*Provider, error)     { return nil, nil }
func discoverDeepSeek(key string) (*Provider, error)    { return nil, nil }
func discoverCerebras(key string) (*Provider, error)    { return nil, nil }
func discoverSambaNova(key string) (*Provider, error)   { return nil, nil }
func discoverFireworks(key string) (*Provider, error)   { return nil, nil }
func discoverPerplexity(key string) (*Provider, error)  { return nil, nil }
func discoverCohere(key string) (*Provider, error)      { return nil, nil }
func discoverLocal() *Provider                           { return nil }
