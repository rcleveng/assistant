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

	chatHandler := chat.NewChatHandler(ctx, environment)
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
	if err := http.ListenAndServe(":"+port, addRequestEnvironment(router, environment)); err != nil {
		log.Fatal(err)
	}
}
