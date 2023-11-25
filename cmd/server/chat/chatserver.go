// Sample run-helloworld is a minimal Cloud Run service.
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

	pb "google.golang.org/api/chat/v1"
)

const chatAppProject = "1007744422436"
const chatIssuer = "chat@system.gserviceaccount.com"
const jwtURL = "https://www.googleapis.com/service_accounts/v1/jwk/"

type ChatHandler struct {
	verifier *oidc.IDTokenVerifier
	llm      server.LlmClient
}

func NewChatHandler(ctx context.Context) *ChatHandler {
	config := &oidc.Config{
		SkipClientIDCheck: true,
		ClientID:          chatAppProject,
	}
	ks := oidc.NewRemoteKeySet(ctx, jwtURL+chatIssuer)
	verifier := oidc.NewVerifier(chatIssuer, ks, config)

	llm, err := server.NewPalmLLMClient(ctx)
	if err != nil {
		panic(err)
	}
	return &ChatHandler{verifier, llm}
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

func (handler *ChatHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
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

	text, err := handler.llm.Call(ctx, originalText)
	if err != nil {
		log.Default().Println("Error: " + err.Error())
		text = fmt.Sprintf(`Error getting LLM, so: Hello '%s', you said '%s'`, name, originalText)
	}

	button := pb.GoogleAppsCardV1Button{
		Text: "Sample Button",
		OnClick: &pb.GoogleAppsCardV1OnClick{
			OpenLink: &pb.GoogleAppsCardV1OpenLink{
				Url: uri + "/chat/action/12344",
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
						Text: "Ipsem Lorem",
					},
				}, {
					ButtonList: &pb.GoogleAppsCardV1ButtonList{
						Buttons: []*pb.GoogleAppsCardV1Button{&button},
					},
				}},
			}},
		},
		CardId: "EchoCard1",
	}

	resp := pb.Message{
		Text:    text,
		CardsV2: []*pb.CardWithId{card},
	}

	server.EncodeAndLogResponse(&resp, w)
}

func (handler *ChatHandler) Close() {
	if handler.llm != nil {
		handler.llm.Close()
	}
}
