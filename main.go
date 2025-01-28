package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid" // Import the google uuid package
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Endpoint to handle prompt submission with UUID generation
	r.Post("/prompt", func(w http.ResponseWriter, r *http.Request) {

		// Extract the prompt submitted

		p := r.FormValue("prompt")
		if p == "" {
			http.Error(w, "Message is required", http.StatusBadRequest)
			return
		}

		pID := uuid.New().String()

		fmt.Printf("New User Prompt Submitted Assigned ID: %s\n", pID)
		// Generate a random UUID for each request using google uuid

		

		// Trigger an event to notify the client
		w.Header().Set("HX-Trigger", fmt.Sprintf(`{"PromptSubmitted": {"id": "%s"}}`, pID))
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/stream-component", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the query parameter value
		id := r.URL.Query().Get("id")
		StreamComponent(id).Render(r.Context(), w)
	})

	r.Get("/stream", func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Token generation configuration
		totalTokens := 1000
		delayPerToken := 100 * time.Millisecond
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		// Use the `llm` function with the context
		for token := range llm(ctx, delayPerToken, totalTokens) {
			// Send the generated token as an SSE message
			fmt.Fprintf(w, "event: update\ndata: %s\n\n", token)
			flusher.Flush()
		}

		// Signal the end of the stream
		fmt.Fprintf(w, "event: close\ndata: Stream completed\n\n")
		flusher.Flush()
	})

	fmt.Println("Server is running on http://localhost:8081")
	http.ListenAndServe(":8081", r)
}
