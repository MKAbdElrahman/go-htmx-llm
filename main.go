package main

import (
	"fmt"
	"net/http"

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
		// Generate a random UUID for each request using google uuid
		promptID := uuid.New().String()

		// Trigger an event to notify the client
		w.Header().Set("HX-Trigger", fmt.Sprintf(`{"PromptSubmitted": {"id": "%s"}}`, promptID))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Prompt submitted with ID: %s", promptID)))
	})

	r.Get("/stream", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the query parameter value
		id := r.URL.Query().Get("id")

		// Simulate streaming data
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: Streaming response for ID: " + id + "\n\n"))
	})

	fmt.Println("Server is running on http://localhost:8081")
	http.ListenAndServe(":8081", r)
}
