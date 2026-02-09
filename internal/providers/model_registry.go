// Model registry - all known models and their availability
package providers

// ModelRegistry contains all known models across providers
var ModelRegistry = map[string]ModelInfo{
	// ═══════════════════════════════════════════════════════════
	// CLAUDE (Anthropic) - also on OpenRouter, AWS Bedrock
	// ═══════════════════════════════════════════════════════════
	"claude-opus-4": {
		Name:      "Claude Opus 4",
		Type:      ModelLLM,
		Context:   200000,
		Providers: []string{"anthropic", "openrouter", "aws-bedrock"},
	},
	"claude-sonnet-4": {
		Name:      "Claude Sonnet 4",
		Type:      ModelLLM,
		Context:   200000,
		Providers: []string{"anthropic", "openrouter", "aws-bedrock", "google-vertex"},
	},
	"claude-haiku-3.5": {
		Name:      "Claude Haiku 3.5",
		Type:      ModelLLM,
		Context:   200000,
		Providers: []string{"anthropic", "openrouter", "aws-bedrock"},
	},

	// ═══════════════════════════════════════════════════════════
	// GPT (OpenAI) - also on Azure, OpenRouter
	// ═══════════════════════════════════════════════════════════
	"gpt-4o": {
		Name:      "GPT-4o",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"openai", "openrouter", "azure"},
	},
	"gpt-4o-mini": {
		Name:      "GPT-4o Mini",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"openai", "openrouter", "azure"},
	},
	"gpt-4-turbo": {
		Name:      "GPT-4 Turbo",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"openai", "openrouter", "azure"},
	},
	"o1": {
		Name:      "o1 (Reasoning)",
		Type:      ModelLLM,
		Context:   200000,
		Providers: []string{"openai"}, // RARE - OpenAI only
	},
	"o1-mini": {
		Name:      "o1-mini",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"openai"}, // RARE
	},
	"o3-mini": {
		Name:      "o3-mini",
		Type:      ModelLLM,
		Context:   200000,
		Providers: []string{"openai"}, // RARE
	},

	// ═══════════════════════════════════════════════════════════
	// GEMINI (Google) - also on OpenRouter
	// ═══════════════════════════════════════════════════════════
	"gemini-2.5-pro": {
		Name:      "Gemini 2.5 Pro",
		Type:      ModelLLM,
		Context:   2000000, // 2M!
		Providers: []string{"google", "openrouter"},
	},
	"gemini-2.5-flash": {
		Name:      "Gemini 2.5 Flash",
		Type:      ModelLLM,
		Context:   1000000,
		Providers: []string{"google", "openrouter"},
	},
	"gemini-2.0-flash-thinking": {
		Name:      "Gemini 2.0 Flash Thinking",
		Type:      ModelLLM,
		Context:   1000000,
		Providers: []string{"google"}, // RARE - Google only
	},

	// ═══════════════════════════════════════════════════════════
	// LLAMA (Meta) - on Groq, Together, Fireworks, Cerebras, etc.
	// ═══════════════════════════════════════════════════════════
	"llama-3.3-70b": {
		Name:      "Llama 3.3 70B",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"groq", "together", "fireworks", "cerebras", "sambanova", "openrouter"},
	},
	"llama-3.1-405b": {
		Name:      "Llama 3.1 405B",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"together", "fireworks", "openrouter"},
	},
	"llama-3.1-70b": {
		Name:      "Llama 3.1 70B",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"groq", "together", "fireworks", "cerebras", "openrouter"},
	},
	"llama-3.1-8b": {
		Name:      "Llama 3.1 8B",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"groq", "together", "fireworks", "cerebras", "sambanova", "openrouter", "local"},
	},

	// ═══════════════════════════════════════════════════════════
	// MISTRAL - on Mistral, Together, OpenRouter
	// ═══════════════════════════════════════════════════════════
	"mistral-large": {
		Name:      "Mistral Large",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"mistral", "openrouter"},
	},
	"mistral-medium": {
		Name:      "Mistral Medium",
		Type:      ModelLLM,
		Context:   32000,
		Providers: []string{"mistral", "openrouter"},
	},
	"mixtral-8x22b": {
		Name:      "Mixtral 8x22B",
		Type:      ModelLLM,
		Context:   65000,
		Providers: []string{"mistral", "together", "groq", "openrouter"},
	},
	"codestral": {
		Name:      "Codestral",
		Type:      ModelLLM,
		Context:   32000,
		Providers: []string{"mistral"}, // RARE - Mistral only
	},

	// ═══════════════════════════════════════════════════════════
	// DEEPSEEK - DeepSeek API, OpenRouter
	// ═══════════════════════════════════════════════════════════
	"deepseek-chat": {
		Name:      "DeepSeek Chat",
		Type:      ModelLLM,
		Context:   64000,
		Providers: []string{"deepseek", "openrouter"},
	},
	"deepseek-coder": {
		Name:      "DeepSeek Coder",
		Type:      ModelLLM,
		Context:   64000,
		Providers: []string{"deepseek", "openrouter"},
	},
	"deepseek-r1": {
		Name:      "DeepSeek R1 (Reasoning)",
		Type:      ModelLLM,
		Context:   64000,
		Providers: []string{"deepseek", "together", "openrouter"},
	},

	// ═══════════════════════════════════════════════════════════
	// QWEN - Alibaba, Together, OpenRouter
	// ═══════════════════════════════════════════════════════════
	"qwen-2.5-72b": {
		Name:      "Qwen 2.5 72B",
		Type:      ModelLLM,
		Context:   128000,
		Providers: []string{"together", "fireworks", "openrouter"},
	},
	"qwen-coder-32b": {
		Name:      "Qwen Coder 32B",
		Type:      ModelLLM,
		Context:   32000,
		Providers: []string{"together", "openrouter"},
	},

	// ═══════════════════════════════════════════════════════════
	// STT (Speech-to-Text)
	// ═══════════════════════════════════════════════════════════
	"whisper-large-v3": {
		Name:      "Whisper Large v3",
		Type:      ModelSTT,
		Providers: []string{"groq", "openai"}, // Groq has FREE tier!
	},
	"whisper-large-v3-turbo": {
		Name:      "Whisper Large v3 Turbo",
		Type:      ModelSTT,
		Providers: []string{"groq"}, // RARE - Groq only, free!
	},

	// ═══════════════════════════════════════════════════════════
	// TTS (Text-to-Speech)
	// ═══════════════════════════════════════════════════════════
	"tts-1": {
		Name:      "OpenAI TTS-1",
		Type:      ModelTTS,
		Providers: []string{"openai"},
	},
	"tts-1-hd": {
		Name:      "OpenAI TTS-1 HD",
		Type:      ModelTTS,
		Providers: []string{"openai"}, // RARE
	},
	"elevenlabs": {
		Name:      "ElevenLabs",
		Type:      ModelTTS,
		Providers: []string{"elevenlabs"}, // RARE
	},

	// ═══════════════════════════════════════════════════════════
	// EMBEDDINGS
	// ═══════════════════════════════════════════════════════════
	"text-embedding-3-large": {
		Name:      "OpenAI Embedding Large",
		Type:      ModelEmbedding,
		Providers: []string{"openai"},
	},
	"text-embedding-3-small": {
		Name:      "OpenAI Embedding Small",
		Type:      ModelEmbedding,
		Providers: []string{"openai"},
	},
	"voyage-3": {
		Name:      "Voyage 3",
		Type:      ModelEmbedding,
		Providers: []string{"voyage"}, // RARE - best for code
	},

	// ═══════════════════════════════════════════════════════════
	// IMAGE GENERATION
	// ═══════════════════════════════════════════════════════════
	"dall-e-3": {
		Name:      "DALL-E 3",
		Type:      ModelImage,
		Providers: []string{"openai"}, // RARE
	},
	"imagen-3": {
		Name:      "Imagen 3",
		Type:      ModelImage,
		Providers: []string{"google"}, // RARE
	},
	"flux-1-pro": {
		Name:      "Flux 1 Pro",
		Type:      ModelImage,
		Providers: []string{"together", "replicate"}, // RARE
	},
}

type ModelInfo struct {
	Name      string
	Type      ModelType
	Context   int
	Providers []string
}

// RareModels returns models only available on one provider
func RareModels() map[string]string {
	rare := make(map[string]string)
	for id, info := range ModelRegistry {
		if len(info.Providers) == 1 {
			rare[id] = info.Providers[0]
		}
	}
	return rare
}

// CommonModels returns models available on multiple providers
func CommonModels() map[string][]string {
	common := make(map[string][]string)
	for id, info := range ModelRegistry {
		if len(info.Providers) > 1 {
			common[id] = info.Providers
		}
	}
	return common
}

// ProviderRareModels returns rare models for each provider
func ProviderRareModels() map[string][]string {
	result := make(map[string][]string)
	for id, info := range ModelRegistry {
		if len(info.Providers) == 1 {
			provider := info.Providers[0]
			result[provider] = append(result[provider], id)
		}
	}
	return result
}

// GetBestProvider returns the best provider for a model based on strategy
func GetBestProvider(modelID string, strategy string) string {
	info, ok := ModelRegistry[modelID]
	if !ok {
		return ""
	}

	if len(info.Providers) == 1 {
		return info.Providers[0]
	}

	// Strategy-based selection
	switch strategy {
	case "free":
		// Prefer free tiers: Groq > Google > OpenRouter
		preferOrder := []string{"groq", "google", "openrouter", "cerebras", "sambanova"}
		for _, pref := range preferOrder {
			for _, p := range info.Providers {
				if p == pref {
					return p
				}
			}
		}
	case "fast":
		// Prefer low latency: Groq > Cerebras > Fireworks
		preferOrder := []string{"groq", "cerebras", "fireworks", "together"}
		for _, pref := range preferOrder {
			for _, p := range info.Providers {
				if p == pref {
					return p
				}
			}
		}
	case "cheap":
		// Prefer low cost: DeepSeek > Together > OpenRouter
		preferOrder := []string{"deepseek", "together", "openrouter", "groq"}
		for _, pref := range preferOrder {
			for _, p := range info.Providers {
				if p == pref {
					return p
				}
			}
		}
	}

	// Default: first provider
	return info.Providers[0]
}
