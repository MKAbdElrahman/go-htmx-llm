package main

import (
	"demo/chat"
	"demo/cmd/components"
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

	// called after POST /prompt
	r.Get("/stream-component", func(w http.ResponseWriter, r *http.Request) {

		components.StreamListner().Render(r.Context(), w)
	})

	r.Get("/prompt-component", func(w http.ResponseWriter, r *http.Request) {
		components.Prompt().Render(r.Context(), w)
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
