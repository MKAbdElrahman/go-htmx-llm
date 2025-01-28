package main

import (
	"demo/chat"
	promptprocessing "demo/prompt-processing"
	"demo/pubsub"
	"log"
	"time"
)

func main() {
	// Create a new in-memory PubSub instance.
	ps := pubsub.NewPubSub()

	chatRepository := chat.NewChatRepository()
	chatService := chat.NewChatService(chatRepository, ps)
	chatService.Start()

	// Create an Ollama LLM engine.
	ollamaEngine := promptprocessing.NewOllamaEngine("llama3.1:8b")
	promptprocessingService := promptprocessing.NewPromptProcessingService(ps, ollamaEngine)
	promptprocessingService.Start()

	chatID := chatService.CreateChat("General")
	prompt, err := chatService.SubmitPrompt(chatID, "What is the weather today?")
	if err != nil {
		log.Fatalf("Error submitting prompt: %v", err)
	}

	// Simulate real-time updates by periodically fetching the recent aggregated tokens.
	go func() {
		for {
			time.Sleep(2 * time.Second)
			response, err := chatService.GetRecentAggregatedTokens(chatID, prompt.Id())
			if err != nil {
				log.Printf("Error fetching recent tokens: %v", err)
			} else {
				log.Printf("Recent response: %s\n", response)
			}
		}
	}()

	select {}
}
