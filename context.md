## File: `./chat/chat_repository.go`
- **Type**: File
- **Extension**: `.go`
### Content:
```go
package chat

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Chat represents a chat in the repository.
type Chat struct {
	id        string
	name      string
	prompts   []Prompt
	createdAt time.Time
	updatedAt time.Time
}

// Prompt represents a prompt in a chat.
type Prompt struct {
	id        string
	text      string
	responses []Response
	createdAt time.Time
	updatedAt time.Time
}

func (p Prompt) Id() string {
	return p.id
}

// Response represents a response to a prompt.
type Response struct {
	id        string
	text      string
	createdAt time.Time
}

var (
	ErrChatNotFound   = errors.New("chat not found")
	ErrPromptNotFound = errors.New("prompt not found")
)

// ChatRepository manages the storage and retrieval of chats, prompts, and responses.
type ChatRepository struct {
	chats map[string]*Chat
	mu    sync.Mutex
}

// NewChatRepository creates a new ChatRepository.
func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		chats: make(map[string]*Chat),
	}
}

// AddChat adds a new chat to the repository and returns its ID.
func (r *ChatRepository) AddChat(name string) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	chat := &Chat{
		id:        uuid.New().String(),
		name:      name,
		prompts:   make([]Prompt, 0),
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	r.chats[chat.id] = chat
	return chat.id
}

// GetChat retrieves a chat by its ID.
func (r *ChatRepository) GetChat(chatId string) (*Chat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	chat, exists := r.chats[chatId]
	if !exists {
		return nil, ErrChatNotFound
	}

	return chat, nil
}

// RenameChat updates the name of an existing chat.
func (r *ChatRepository) RenameChat(chatId, newName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	chat, exists := r.chats[chatId]
	if !exists {
		return ErrChatNotFound
	}

	chat.name = newName
	chat.updatedAt = time.Now()
	return nil
}

// DeleteChat removes a chat from the repository.
func (r *ChatRepository) DeleteChat(chatId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.chats[chatId]
	if !exists {
		return ErrChatNotFound
	}

	delete(r.chats, chatId)
	return nil
}

// SubmitPrompt submits a prompt to a chat.
func (r *ChatRepository) SubmitPrompt(chatId, promptText string) (*Prompt, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	chat, exists := r.chats[chatId]
	if !exists {
		return nil, ErrChatNotFound
	}

	prompt := Prompt{
		id:        uuid.New().String(),
		text:      promptText,
		responses: make([]Response, 0),
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	chat.prompts = append(chat.prompts, prompt)
	chat.updatedAt = time.Now()

	return &prompt, nil
}

// AddResponseToPrompt adds a response to a specific prompt in a chat.
func (r *ChatRepository) AddResponseToPrompt(chatId, promptId, responseText string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	chat, exists := r.chats[chatId]
	if !exists {
		return ErrChatNotFound
	}

	for i, prompt := range chat.prompts {
		if prompt.id == promptId {
			response := Response{
				id:        uuid.New().String(),
				text:      responseText,
				createdAt: time.Now(),
			}
			prompt.responses = append(prompt.responses, response)
			prompt.updatedAt = time.Now()
			chat.prompts[i] = prompt
			chat.updatedAt = time.Now()
			return nil
		}
	}

	return ErrPromptNotFound
}

```

## File: `./chat/chat_service.go`
- **Type**: File
- **Extension**: `.go`
### Content:
```go
package chat

import (
	"demo/pubsub"
	"log"
	"strings"
	"sync"
)

// ChatService orchestrates operations on chats, prompts, and responses.
type ChatService struct {
	repo   *ChatRepository
	pubSub *pubsub.PubSub
	mu     sync.Mutex
}

// NewChatService creates a new ChatService with the given repository and PubSub system.
func NewChatService(repo *ChatRepository, pubSub *pubsub.PubSub) *ChatService {
	return &ChatService{
		repo:   repo,
		pubSub: pubSub,
	}
}

// CreateChat creates a new chat and publishes a "ChatCreated" event.
func (s *ChatService) CreateChat(name string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	chatID := s.repo.AddChat(name)
	s.pubSub.Publish("ChatCreated", map[string]interface{}{
		"chatId": chatID,
		"name":   name,
	})

	return chatID
}

// RenameChat renames an existing chat and publishes an event.
func (s *ChatService) RenameChat(chatID, newName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.repo.RenameChat(chatID, newName)
	if err != nil {
		return err
	}

	// Publish a "ChatRenamed" event.
	s.pubSub.Publish("ChatRenamed", map[string]interface{}{
		"chatId":  chatID,
		"newName": newName,
	})

	return nil
}

// DeleteChat deletes a chat and publishes an event.
func (s *ChatService) DeleteChat(chatID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.repo.DeleteChat(chatID)
	if err != nil {
		return err
	}

	// Publish a "ChatDeleted" event.
	s.pubSub.Publish("ChatDeleted", map[string]interface{}{
		"chatId": chatID,
	})

	return nil
}

// SubmitPrompt submits a prompt, stores it in the repository, and publishes a "PromptSubmitted" event.
func (s *ChatService) SubmitPrompt(chatID, promptText string) (*Prompt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prompt, err := s.repo.SubmitPrompt(chatID, promptText)
	if err != nil {
		return nil, err
	}

	s.pubSub.Publish("PromptSubmitted", map[string]interface{}{
		"chatId":     chatID,
		"promptId":   prompt.id,
		"promptText": promptText,
	})

	return prompt, nil
}

// HandleTokensGenerated processes TokensGenerated events and updates the prompt with the response.
func (s *ChatService) HandleTokensGenerated(chatId, promptId, responseText string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Append the token to the prompt's response.
	err := s.repo.AddResponseToPrompt(chatId, promptId, responseText)
	if err != nil {
		log.Printf("Error adding response to prompt: %v\n", err)
		return err
	}
	log.Printf("Token added to ChatID=%s, PromptID=%s: %s\n", "***", promptId, responseText)
	return nil
}

// GetRecentAggregatedTokens retrieves the recent aggregated tokens as a sentence.
func (s *ChatService) GetRecentAggregatedTokens(chatId, promptId string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, err := s.repo.GetChat(chatId)
	if err != nil {
		return "", err
	}

	for _, prompt := range chat.prompts {
		if prompt.id == promptId {
			var responseText strings.Builder
			for _, response := range prompt.responses {
				responseText.WriteString(response.text)
			}
			return responseText.String(), nil
		}
	}

	return "", ErrPromptNotFound
}

func (s *ChatService) ListenForTokensGenerated() <-chan string {
	tokenCh := make(chan string, 100)
	go func() {
		s.pubSub.Subscribe("TokensGenerated", func(payload interface{}) {
			data, ok := payload.(map[string]interface{})
			if !ok {
				log.Println("Invalid payload for TokensGenerated event")
				return
			}

			// Extract event data.
			chatID, ok := data["chatId"].(string)
			if !ok {
				log.Println("Invalid chatId in TokensGenerated event")
				return
			}

			promptID, ok := data["promptId"].(string)
			if !ok {
				log.Println("Invalid promptId in TokensGenerated event")
				return
			}

			responseText, ok := data["responseText"].(string)
			if !ok {
				log.Println("Invalid responseText in TokensGenerated event")
				return
			}

			go func() {
				// Send the response text to the token channel.
				tokenCh <- responseText
			}()

			// Handle the TokensGenerated event and send token to the channel.
			err := s.HandleTokensGenerated(chatID, promptID, responseText)
			if err != nil {
				log.Printf("Failed to handle TokensGenerated event: %v\n", err)
				return
			}

		})
	}()

	return tokenCh
}

```

## File: `./promptprocessing/ollama-engine.go`
- **Type**: File
- **Extension**: `.go`
### Content:
```go
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

```

## File: `./promptprocessing/prompt_processing_service.go`
- **Type**: File
- **Extension**: `.go`
### Content:
```go
package promptprocessing

import (
	"context"
	"demo/pubsub"
	"log"
)

// LLMEngineType defines the interface for any LLM engine.
type LLMEngineType interface {
	// Starts generating tokens and returns a channel for streaming responses.
	GenerateTokens(ctx context.Context, prompt string) (<-chan string, error)
	// Attempts to stop a request mid-processing.
	StopGeneration(ctx context.Context, prompt string) error
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

```

## File: `./cmd/main.go`
- **Type**: File
- **Extension**: `.go`
### Content:
```go
package main

import (
	"demo/chat"
	"demo/promptprocessing"
	"demo/pubsub"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	ps := pubsub.NewPubSub()

	chatRepository := chat.NewChatRepository()
	chatService := chat.NewChatService(chatRepository, ps)
	tokensCh := chatService.ListenForTokensGenerated()

	// Create an Ollama LLM engine.
	ollamaEngine := promptprocessing.NewOllamaEngine("llama3.1:8b")
	promptprocessingService := promptprocessing.NewPromptProcessingService(ps, ollamaEngine)
	promptprocessingService.Start()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Endpoint to handle prompt submission with UUID generation
	r.Post("/prompt", func(w http.ResponseWriter, r *http.Request) {
		// Extract the prompt submitted
		txt := r.FormValue("prompt")
		if txt == "" {
			http.Error(w, "Message is required", http.StatusBadRequest)
			return
		}

		chatId := chatService.CreateChat("TestChat")
		p, err := chatService.SubmitPrompt(chatId, txt)
		if err != nil {
			http.Error(w, "Failed to submit prompt", http.StatusInternalServerError)
			return
		}

		// Trigger an event to notify the client
		w.Header().Set("HX-Trigger", fmt.Sprintf(`{"PromptSubmitted": {"id": "%s"}}`, p.Id()))
		w.WriteHeader(http.StatusOK)
	})
	r.Post("/stop", func(w http.ResponseWriter, r *http.Request) {
		prompt := r.FormValue("prompt")
		if prompt == "" {
			http.Error(w, "prompt is required", http.StatusBadRequest)
			return
		}

		err := ollamaEngine.StopGeneration(r.Context(), prompt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to stop: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Generation stopped successfully"))
	})

	r.Get("/stream", func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Ensure the response writer supports flushing
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Get the request context
		ctx := r.Context()

		// Send initial message to confirm connection
		fmt.Fprintf(w, "event: connected\ndata: Connection established\n\n")
		flusher.Flush()

		for {
			select {
			case token, ok := <-tokensCh:
				if !ok {
					// Channel closed, signal the end of the stream
					fmt.Fprintf(w, "event: close\ndata: Stream completed\n\n")
					flusher.Flush()
					return
				}

				// Send the generated token as an SSE message
				fmt.Fprintf(w, "event: update\ndata: %s\n\n", token)
				flusher.Flush()
			case <-ctx.Done():
				// Client disconnected
				log.Println("Client disconnected")
				return
			}
		}
	})

	fmt.Println("Server is running on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}

```

## File: `./index.html`
- **Type**: File
- **Extension**: `.html`
### Content:
```html
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chat Application</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
    <script src="https://unpkg.com/htmx-ext-sse@2.2.2/sse.js"></script>
</head>

<body class="bg-[#1a1a1a] text-[#e5e5e5] p-6 flex flex-col h-screen">

    <!-- Display streaming response -->
    <div
        id="main-content"
        class="flex-grow p-8 mt-16 overflow-y-auto scrollbar-thin scrollbar-thumb-[#4C9C94] scrollbar-track-[#1a1a1a] border border-[#3a3a3c] rounded-lg mb-4"
    >
        <div 
            id="stream-response" 
            hx-ext="sse" 
            sse-connect="/stream" 
            sse-swap="update" 
            hx-swap="beforeend">
            <!-- Responses will be appended here -->
        </div>
    </div>

    <!-- Prompt area -->
    <div class="p-6 flex flex-col">
        <div class="text-[#e5e5e5] flex-grow flex flex-col">
            <form
                hx-post="/prompt" 
                hx-swap="none" 
                hx-on:htmx:before-request="document.getElementById('stream-response').innerHTML = '';"
                class="flex items-center space-x-2 bg-[#1a1a1a] rounded-lg border border-[#3a3a3c] p-2 hover:border-[#4C9C94] transition-colors duration-200"
            >
                <!-- Send Icon (Interactive Animations) -->
                <button
                    type="submit"
                    class="p-1 text-[#4C9C94] hover:text-[#007acc] transition-colors duration-200 flex items-center justify-center group"
                >
                    <svg
                        class="w-4 h-4 hover:w-5 hover:h-5 transition-all duration-200 animate-bounce group-hover:animate-pulse group-active:animate-ping"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M13 5l7 7-7 7M5 5l7 7-7 7"
                        ></path>
                    </svg>
                </button>
                <!-- Input Field -->
                <input
                    type="text"
                    id="prompt-input"
                    name="prompt"
                    placeholder="Type your prompt..."
                    class="w-full p-2 bg-transparent text-[#e5e5e5] focus:outline-none placeholder-[#a1a1aa]"
                    required
                />
                <input type="hidden" id="prompt-index" name="prompt-index" value="-1"/>
            </form>
        </div>
    </div>
</body>

</html>
```

## File: `./prompt.md`
- **Type**: File
- **Extension**: `.md`
### Content:
```markdown
## Requirements

add a stop method to the engine type 
```

