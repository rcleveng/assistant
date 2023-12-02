package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rcleveng/assistant/server"
	"github.com/rcleveng/assistant/server/db"
	"github.com/rcleveng/assistant/server/env"
	"github.com/rcleveng/assistant/server/llm"

	pb "google.golang.org/api/chat/v1"
)

const chatAppProject = "1007744422436"
const chatIssuer = "chat@system.gserviceaccount.com"
const jwtURL = "https://www.googleapis.com/service_accounts/v1/jwk/"

type ChatHandler struct {
	verifier *oidc.IDTokenVerifier
	llm      llm.LlmClient
	db       db.EmbeddingsDB
}

func NewChatHandler(ctx context.Context, environment *env.Environment) (*ChatHandler, error) {
	config := &oidc.Config{
		SkipClientIDCheck: true,
		ClientID:          chatAppProject,
	}
	ks := oidc.NewRemoteKeySet(ctx, jwtURL+chatIssuer)
	verifier := oidc.NewVerifier(chatIssuer, ks, config)

	llm, err := llm.NewPalmLLMClient(ctx, environment)
	if err != nil {
		return nil, err
	}

	edb, err := db.NewPostgresDatabase(environment)
	if err != nil {
		return nil, err
	}

	return &ChatHandler{
		verifier: verifier,
		llm:      llm,
		db:       edb,
	}, err
}

// Validate the Chat Token

func (handler *ChatHandler) validateChatToken(context context.Context, tokenString string, chatAppProject string) error {

	payload, err := handler.verifier.Verify(context, tokenString)
	if err != nil {
		return err
	}
	var claims struct {
		Aud string `json:"aud"`
		Iss string `json:"iss"`
	}
	if err := payload.Claims(&claims); err != nil {
		return err
	}

	fmt.Printf("\n\nAud: %s; Iss: %s\n\n", claims.Aud, claims.Iss)
	if claims.Aud != chatAppProject {
		return fmt.Errorf("audience was not the correct chat project, got '%s'", claims.Aud)
	}
	if claims.Iss != chatIssuer {
		return fmt.Errorf("issuer was not '%s'; got: '%s'", chatIssuer, claims.Iss)
	}
	return nil
}

// CHAT

func CreateResponseCard(cardId, sessionId, text, uri string) (*pb.Message, error) {
	button := pb.GoogleAppsCardV1Button{
		Text: "Learn More",
		OnClick: &pb.GoogleAppsCardV1OnClick{
			OpenLink: &pb.GoogleAppsCardV1OpenLink{
				Url: uri + "/chat/action/12344",
			},
		},
	}
	// <a href="https://www.flaticon.com/free-icons/thumbs-down" title="thumbs down icons">Thumbs down icons created by Freepik - Flaticon</a>

	thumbsUpButton := pb.GoogleAppsCardV1Button{
		Icon: &pb.GoogleAppsCardV1Icon{
			AltText:   "Thumbs Up",
			IconUrl:   "https://storage.cloud.google.com/robsite-assistant-public/thumb-up.png",
			ImageType: "SQUARE",
		},
		//Text: "Thumbs Up",
		OnClick: &pb.GoogleAppsCardV1OnClick{
			OpenLink: &pb.GoogleAppsCardV1OpenLink{
				Url: fmt.Sprintf("%s/feedback/up/%s", uri, sessionId),
			},
		},
	}
	thumbsDownButton := pb.GoogleAppsCardV1Button{
		//Text: "Thumbs Down",
		Icon: &pb.GoogleAppsCardV1Icon{
			AltText:   "Thumbs Down",
			IconUrl:   "https://storage.cloud.google.com/robsite-assistant-public/thumb-down.png",
			ImageType: "SQUARE",
		},
		OnClick: &pb.GoogleAppsCardV1OnClick{
			OpenLink: &pb.GoogleAppsCardV1OpenLink{
				Url: fmt.Sprintf("%s/feedback/down/%s", uri, sessionId),
			},
		},
	}
	card := &pb.CardWithId{
		Card: &pb.GoogleAppsCardV1Card{
			Header: &pb.GoogleAppsCardV1CardHeader{
				ImageType: "CIRCLE",
				ImageUrl:  "https://developers.google.com/chat/images/chat-product-icon.png",
				Subtitle:  "Created with the Robsite Assistant",
				Title:     "Robsite Assistant Reply",
			},
			Sections: []*pb.GoogleAppsCardV1Section{{
				Widgets: []*pb.GoogleAppsCardV1Widget{{
					TextParagraph: &pb.GoogleAppsCardV1TextParagraph{
						Text: text,
					},
				}, {
					ButtonList: &pb.GoogleAppsCardV1ButtonList{
						Buttons: []*pb.GoogleAppsCardV1Button{
							&button, &thumbsUpButton, &thumbsDownButton},
					},
				}},
			}},
		},
		CardId: cardId,
	}

	resp := &pb.Message{
		CardsV2: []*pb.CardWithId{card},
	}
	return resp, nil
}

func (handler *ChatHandler) DebugCard(w http.ResponseWriter, r *http.Request) {
	uri := server.GetPublicEndpoint(r)
	log.Default().Println("URI: " + uri)
	if reqText, err := httputil.DumpRequest(r, true); err == nil {
		log.Default().Println(string(reqText))
	}
	r.URL.Query()

	sessionId := r.URL.Query().Get("id")
	text := r.URL.Query().Get("text")
	if text == "" {
		text = "This is a test"
	}
	resp, err := CreateResponseCard("TestCard", sessionId, text, uri)
	if err != nil {
		resp = &pb.Message{
			Text: "Error creating Card",
		}
	}
	server.EncodeAndLogResponse(resp, w)
}

func (handler *ChatHandler) runChain(ctx context.Context, cmd, rest, name string) (string, error) {
	switch cmd {
	case "ANSWER":
		return rest, nil
	case "NEEDMORE":
		// TODO - ask this as a follow up, then add the next answer to the context and
		// retry the original question
		return rest, nil
	case "CALENDAR":
		return fmt.Sprintf("I would use the calendar to look up '%s'", rest), nil
	case "REMEMBER":
		emb, err := handler.llm.EmbedText(ctx, rest)
		if err != nil {
			return rest, fmt.Errorf("error trying to remember: ")
		}
		_, err = handler.db.Add(0, rest, emb)
		if err != nil {
			return rest, fmt.Errorf("error trying to remember: ")
		}
		return fmt.Sprintf("I will remember that '%s'", rest), nil
	default:
		return cmd + " " + rest, nil
	}

}

func (handler *ChatHandler) handleChat(w http.ResponseWriter, r *http.Request, name, sessionId, text string) (string, error) {
	ctx := r.Context()

	emb, err := handler.llm.EmbedText(ctx, text)
	if err != nil {
		return "", err
	}

	context, err := handler.db.Find(emb, 2)
	if err != nil {
		context = []string{}
	}

	prompt, err := llm.ChatPrompt(text, context)
	if err != nil {
		fmt.Println("error generating chat prompt", err.Error())
		prompt = text
	}
	responseText, err := handler.llm.GenerateText(ctx, prompt)
	if err != nil {
		return "", err
	}

	// TODO - process response for remembering, looking up calendar and starting chain, etc
	cmd, rest, found := strings.Cut(responseText, " ")
	if found && cmd[len(cmd)-1:] == ":" {
		// we have a command
		cmd = cmd[:len(cmd)-1]
		responseText, err = handler.runChain(ctx, cmd, rest, name)
		if err != nil {
			fmt.Println("running chain failed ", err.Error())
			return "", err
		}
	}
	return responseText, nil
}

type BasicChat struct {
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
}

func (handler *ChatHandler) HandleChatBasic(w http.ResponseWriter, r *http.Request) {
	uri := server.GetPublicEndpoint(r)
	log.Default().Println("HandleChatBasic URI: " + uri)

	req := &BasicChat{}
	json.NewDecoder(r.Body).Decode(&req)
	fmt.Printf("Decoded Message: %#v", req)

	sessionId := "0"
	text, err := handler.handleChat(w, r, req.Name, sessionId, req.Text)
	if err != nil {
		log.Default().Println("Error: " + err.Error())
		server.EncodeAndLogResponse(&pb.Message{
			Text: "Error: " + err.Error(),
		}, w)
		return
	}

	server.EncodeAndLogResponse(&pb.Message{
		Text: text,
	}, w)

}

func (handler *ChatHandler) HandleChatApp(w http.ResponseWriter, r *http.Request) {
	uri := server.GetPublicEndpoint(r)
	log.Default().Println("URI: " + uri)
	if reqText, err := httputil.DumpRequest(r, true); err == nil {
		log.Default().Println(string(reqText))
	}
	ctx := r.Context()

	if authHeader, ok := r.Header["Authorization"]; ok {
		token := strings.Split(authHeader[0], " ")[1]
		if err := handler.validateChatToken(ctx, token, chatAppProject); err != nil {
			fmt.Printf("Error validating Token: %v\n", err)
			fmt.Fprintf(w, "Error validating Token: %v\n", err)
			return
		}
	} else {
		// TODO return error to client here.
		fmt.Println("No Auth Header!")
		return
	}

	req := pb.DeprecatedEvent{}
	json.NewDecoder(r.Body).Decode(&req)
	fmt.Printf("Decoded Message: %#v", req)

	if req.Message == nil || req.Message.Sender == nil {
		fmt.Println("No Message or Sender")
		return
	}

	name := req.Message.Sender.DisplayName
	_, originalText, _ := strings.Cut(req.Message.Text, " ")

	// todo - make a real session id
	sessionId := "0"

	text, err := handler.handleChat(w, r, name, sessionId, originalText)
	if err != nil {
		log.Default().Println("Error: " + err.Error())
		server.EncodeAndLogResponse(&pb.Message{
			Text: "Error creating Card",
		}, w)
		return
	}

	resp, err := CreateResponseCard("ChatResponseCard", sessionId, text, uri)
	if err != nil {
		log.Default().Println("Error: " + err.Error())
		server.EncodeAndLogResponse(&pb.Message{
			Text: "Error creating Card",
		}, w)
		return
	}

	server.EncodeAndLogResponse(resp, w)
}

func (handler *ChatHandler) Close() {
	if handler.llm != nil {
		handler.llm.Close()
	}
}
