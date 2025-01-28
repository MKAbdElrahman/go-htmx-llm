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

	log.Printf("Token added to ChatID=%s, PromptID=%s: %s\n", chatId, promptId, responseText)
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
	"strings"

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

		var buffer strings.Builder
		for {
			select {
			case token, ok := <-tokensCh:
				if !ok {
					// Channel closed, signal the end of the stream
					if buffer.Len() > 0 {
						fmt.Fprintf(w, "event: update\ndata: %s\n\n", buffer.String())
						buffer.Reset()
						flusher.Flush()
					}
					fmt.Fprintf(w, "event: close\ndata: Stream completed\n\n")
					flusher.Flush()
					return
				}

				log.Printf("read token %s from channel\n", token)
				// Append token to buffer
				buffer.WriteString(token + " ")

				// Check if the buffer contains a complete sentence (ends with a period)
				if strings.HasSuffix(strings.TrimSpace(buffer.String()), ".") {
					fmt.Fprintf(w, "event: update\ndata: %s\n\n", buffer.String())
					buffer.Reset()
					flusher.Flush()
				}
			case <-ctx.Done():
				// Client disconnected
				log.Println("Client disconnected")
				return
			}
		}
	})

	fmt.Println("Server is running on http://localhost:8081")
	http.ListenAndServe(":8081", r)
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
    
    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
    <script src="https://unpkg.com/htmx-ext-sse@2.2.2/sse.js"></script>
</head>

<body>
    <h1>Chat Application</h1>

    <!-- Form for submitting a prompt -->
    <form hx-post="/prompt" hx-swap="none" hx-on:htmx:before-request="document.getElementById('stream-response').innerHTML = '';">
        <input type="text" id="prompt" name="prompt" placeholder="Enter your prompt..." required>
        <button type="submit">Submit</button>
    </form>

    <!-- Display streaming response -->
    <div id="stream-response" hx-ext="sse" sse-connect="/stream" sse-swap="update" hx-swap="beforeend">
        <!-- Responses will be appended here -->
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

 - fix the stream endpoint

 the event stream look like thisin chrome  chat
connected	Connection established	
17:52:34.814
update	are	
17:52:38.699
update	today	
17:52:39.100
update	Is	
17:52:39.511
update	something	
17:52:39.893
update	can	
17:52:40.277
update	with	
17:52:40.650
update	would	
17:52:41.034
update	like	
17:52:41.439
update	chat	
17:52:41.815
update		
17:52:42.191
```

