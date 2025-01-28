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
