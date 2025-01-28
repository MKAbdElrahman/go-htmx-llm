package promptprocessing

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// OllamaEngine implements the LLMEngineType interface using the Ollama model.
type OllamaEngine struct {
	model string
}

// NewOllamaEngine creates a new OllamaEngine with the specified model.
func NewOllamaEngine(model string) *OllamaEngine {
	return &OllamaEngine{
		model: model,
	}
}

func (o *OllamaEngine) GenerateTokens(ctx context.Context, prompt string) (<-chan string, error) {
	llm, err := ollama.New(ollama.WithModel(o.model))
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama LLM: %w", err)
	}

	// Use a buffered channel to avoid dropping tokens
	tokenChan := make(chan string, 100) // Buffer size of 100

	go func() {
		defer close(tokenChan)

		_, err := llm.Call(ctx, prompt,
			llms.WithTemperature(0.8),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				// Send the token to the buffered channel
				tokenChan <- string(chunk)
				return nil
			}),
		)
		if err != nil {
			log.Printf("Error generating tokens: %v", err)
		}
	}()

	return tokenChan, nil
}
