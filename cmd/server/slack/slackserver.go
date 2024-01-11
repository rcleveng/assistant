package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rcleveng/assistant/server"
	"github.com/rcleveng/assistant/server/db"
	"github.com/rcleveng/assistant/server/env"
	"github.com/rcleveng/assistant/server/llm"
)

const defaultChatAppProject = "1007744422436"

// const chatIssuer = "chat@system.gserviceaccount.com"
// const jwtURL = "https://www.googleapis.com/service_accounts/v1/jwk/"

type SlackHandler struct {
	llm       llm.LlmClient
	db        db.EmbeddingsDB
	projectID string
}

func NewSlackHandler(ctx context.Context, environment *env.Environment, router *mux.Router) (*SlackHandler, error) {
	// TODO: query the metadata server if we're on cloud run.
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = defaultChatAppProject
	}

	slog.Info("NewSlackHandler: Using Cloud", "projectID", projectID)

	llm, err := llm.NewPalmLLMClient(ctx, environment)
	if err != nil {
		return nil, err
	}

	edb, err := db.NewPostgresDatabase(environment)
	if err != nil {
		return nil, err
	}

	handler := &SlackHandler{
		// verifier:  verifier,
		llm:       llm,
		db:        edb,
		projectID: projectID,
	}

	// TODO - remove GET from these
	router.HandleFunc("/commands/help", handler.slashHelp).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/action-endpoint", handler.actionEndpoint).Methods(http.MethodPost, http.MethodGet)

	// This is jsut here for testing
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Slack handler is active.\n")) }).Methods(http.MethodGet)

	return handler, err
}

func (handler *SlackHandler) slashHelp(w http.ResponseWriter, r *http.Request) {
	// TODO publish help
	uri := server.GetPublicEndpoint(r)
	slog.Info("SlackHandler:slashHelp URI: " + uri)

	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "SlackHandler:slashHelp", "error", err)
		return
	}
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, body, "", "\t"); err != nil {
		slog.ErrorContext(ctx, "JSON parse error", "error", err)
		return
	}
	slog.InfoContext(ctx, "SlackHandler:slashHelp", "request", prettyJSON.String())
}

func (handler *SlackHandler) actionEndpoint(w http.ResponseWriter, r *http.Request) {
	// TODO publish help
	uri := server.GetPublicEndpoint(r)
	slog.Info("SlackHandler:slashHelp URI: " + uri)

	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "SlackHandler:actionEndpoint", "error", err)
		return
	}
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, body, "", "\t"); err != nil {
		slog.ErrorContext(ctx, "JSON parse error", "error", err)
		return
	}
	slog.InfoContext(ctx, "SlackHandler:actionEndpoint", "request", prettyJSON.String())
}

func (handler *SlackHandler) Close() {
	if handler.llm != nil {
		handler.llm.Close()
	}
}
