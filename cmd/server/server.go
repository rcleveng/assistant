// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rcleveng/assistant/cmd/server/chat"
	"github.com/rcleveng/assistant/cmd/server/docs"

	"os"
)

func main() {
	cloudRunExecution := os.Getenv("CLOUD_RUN_EXECUTION")
	log.Print("starting server: " + cloudRunExecution)
	router := mux.NewRouter()

	chatHandler := chat.NewChatHandler()
	defer chatHandler.Close()

	router.HandleFunc("/docs", docs.DocsHandler).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/authorizeFile", docs.AuthFileHandler).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/chat", chatHandler.HandleRequest).Methods(http.MethodPost, http.MethodGet)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
