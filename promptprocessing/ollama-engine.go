package promptprocessing

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// OllamaEngine implements the LLMEngineType interface using the Ollama model.
type OllamaEngine struct {
	model       string
	mu          sync.Mutex
	activeTasks map[string]context.CancelFunc
}

func NewOllamaEngine(model string) *OllamaEngine {
	return &OllamaEngine{
		model:       model,
		activeTasks: make(map[string]context.CancelFunc),
	}
}

func (o *OllamaEngine) GenerateTokens(ctx context.Context, prompt string) (<-chan string, error) {
	o.mu.Lock()
	ctx, cancel := context.WithCancel(ctx)
	o.activeTasks[prompt] = cancel
	o.mu.Unlock()

	tokenChan := make(chan string, 100)

	go func() {
		defer close(tokenChan)
		defer func() {
			o.mu.Lock()
			delete(o.activeTasks, prompt)
			o.mu.Unlock()
		}()

		llm, err := ollama.New(ollama.WithModel(o.model))
		if err != nil {
			log.Printf("Failed to create Ollama LLM: %v", err)
			return
		}

		_, err = llm.Call(ctx, prompt,
			llms.WithTemperature(0.8),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				select {
				case <-ctx.Done():
					log.Println("Context canceled, stopping token generation")
					return ctx.Err()
				case tokenChan <- string(chunk):
				}
				return nil
			}),
		)

		if err != nil && ctx.Err() == nil {
			log.Printf("Error generating tokens: %v", err)
		}
	}()

	return tokenChan, nil
}

func (o *OllamaEngine) StopGeneration(ctx context.Context, prompt string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	cancel, exists := o.activeTasks[prompt]
	if !exists {
		return fmt.Errorf("prompt %q not found or already completed", prompt)
	}

	cancel()
	delete(o.activeTasks, prompt)

	return nil
}
