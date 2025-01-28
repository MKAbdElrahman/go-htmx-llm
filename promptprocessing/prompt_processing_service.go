package promptprocessing

import (
	"context"
	"demo/pubsub"
	"log"
)

// LLMEngineType defines the interface for any LLM engine.
type LLMEngineType interface {
	GenerateTokens(ctx context.Context, prompt string) (<-chan string, error)
}

// PromptProcessingService handles processing prompts.
type PromptProcessingService struct {
	pubSub    *pubsub.PubSub
	llmEngine LLMEngineType
}

// NewPromptProcessingService creates a new PromptProcessingService with the given LLM engine.
func NewPromptProcessingService(pubSub *pubsub.PubSub, llmEngine LLMEngineType) *PromptProcessingService {
	return &PromptProcessingService{
		pubSub:    pubSub,
		llmEngine: llmEngine,
	}
}

func (s *PromptProcessingService) Start() {
	s.pubSub.Subscribe("PromptSubmitted", func(payload interface{}) {
		data, ok := payload.(map[string]interface{})
		if !ok {
			log.Println("Invalid payload for PromptSubmitted event")
			return
		}

		// Extract event data
		chatID := data["chatId"].(string)
		promptID := data["promptId"].(string)
		promptText := data["promptText"].(string)

		log.Printf("Processing prompt: ChatID=%s, PromptID=%s, Text=%s\n", chatID, promptID, promptText)

		// Generate tokens using the LLM engine
		ctx := context.Background()
		tokenChan, err := s.llmEngine.GenerateTokens(ctx, promptText)
		if err != nil {
			log.Printf("Error generating tokens: %v", err)
			return
		}

		// Publish TokensGenerated events for each token
		go func() {
			for token := range tokenChan {
				log.Printf("Generated token for ChatID=%s, PromptID=%s: %s\n", "***", promptID, token)
				s.pubSub.Publish("TokensGenerated", map[string]interface{}{
					"chatId":       chatID,
					"promptId":     promptID,
					"responseText": token,
				})
			}
		}()
	})
}
