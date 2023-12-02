// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rcleveng/assistant/cmd/server/chat"
	"github.com/rcleveng/assistant/cmd/server/docs"
	"github.com/rcleveng/assistant/server/env"

	"os"
)

// TODO - This may go away if we don't need Environment per-request since the
// handlers have it already.
func addRequestEnvironment(next http.Handler, environment *env.Environment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := env.NewContext(r.Context(), environment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	cloudRunExecution := os.Getenv("CLOUD_RUN_EXECUTION")
	log.Print("starting server: " + cloudRunExecution)
	router := mux.NewRouter()

	// TODO(rcleveng): Use correct environment
	environment, err := env.NewEnvironment(env.GOTEST)
	if err != nil {
		log.Fatal(err)
	}
	ctx := env.NewContext(context.Background(), environment)

	chatHandler, err := chat.NewChatHandler(ctx, environment)
	if err != nil {
		log.Fatal(err)
	}

	defer chatHandler.Close()
	feedbackHandler := chat.FeedbackHandler{}
	fs := http.FileServer(http.Dir("./static_files/"))
	//router.Handle("/debug/chat/", fs)
	router.Handle("/", fs)

	router.HandleFunc("/docs", docs.DocsHandler).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/authorizeFile", docs.AuthFileHandler).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/chat", chatHandler.HandleChatApp).Methods(http.MethodPost)
	router.HandleFunc("/chat/basic", chatHandler.HandleChatBasic).Methods(http.MethodPost)
	router.HandleFunc("/debug/card", chatHandler.DebugCard).Methods(http.MethodGet)
	router.PathPrefix("/debug/chat/").Handler(http.StripPrefix("/debug/chat/", fs))
	router.HandleFunc("/feedback/{type}/{id}", feedbackHandler.HandleFeedback).Methods(http.MethodPost)
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{'o', 'k', '\n'})
	})

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, addRequestEnvironment(router, environment)); err != nil {
		log.Fatal(err)
	}
}
