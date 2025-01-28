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
