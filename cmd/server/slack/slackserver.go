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

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const defaultChatAppProject = "1007744422436"

// const chatIssuer = "chat@system.gserviceaccount.com"
// const jwtURL = "https://www.googleapis.com/service_accounts/v1/jwk/"

type SlackHandler struct {
	llm           llm.LlmClient
	db            db.EmbeddingsDB
	api           *slack.Client
	clientId      string
	clientSecret  string
	signingSecret string
	projectID     string
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

	api := slack.New(environment.SlackBotOAuthToken)

	handler := &SlackHandler{
		llm:           llm,
		db:            edb,
		projectID:     projectID,
		api:           api,
		clientId:      environment.SlackClientID,
		clientSecret:  environment.SlackClientSecret,
		signingSecret: environment.SlackSigningSecret,
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

	ctx := r.Context()
	uri := server.GetPublicEndpoint(r)
	slog.InfoContext(ctx, "SlackHandler:actionEndpoint called", "URI", uri)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "SlackHandler:actionEndpoint", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, handler.signingSecret)
	if err != nil {
		slog.ErrorContext(ctx, "SlackHandler:actionEndpoint", "error", "StatusBadRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		slog.ErrorContext(ctx, "SlackHandler:actionEndpoint", "error", "StatusInternalServerError")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO - this is failing, fix soon
	// if err := sv.Ensure(); err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// Debug log full request
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, body, "", "\t"); err != nil {
		slog.ErrorContext(ctx, "JSON parse error", "error", err)
		return
	}
	slog.InfoContext(ctx, "SlackHandler:actionEndpoint", "request", prettyJSON.String())

	// TODO - verify token here?
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.ErrorContext(ctx, "SlackHandler:actionEndpoint error parsing event", "error", "StatusInternalServerError")
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			slog.ErrorContext(ctx, "SlackHandler:actionEndpoint error Unmarshal event", "error", "StatusInternalServerError")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		slog.InfoContext(ctx, "Responded with Challenge")
		w.Write([]byte(r.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		slog.InfoContext(ctx, "Received CallbackEvent")
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			channel, ts, err := handler.api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err == nil {
				slog.InfoContext(ctx, "posted message", "channel", channel, "timestamp", ts)
			} else {
				slog.ErrorContext(ctx, "error posting message to channel", "err", err)
			}
		default:
			slog.InfoContext(ctx, "unhandled event type", "ev", ev)
		}
	}
}

func (handler *SlackHandler) Close() {
	if handler.llm != nil {
		handler.llm.Close()
	}
}
