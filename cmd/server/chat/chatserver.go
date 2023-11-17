// Sample run-helloworld is a minimal Cloud Run service.
package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/tidwall/gjson"

	pb "google.golang.org/api/chat/v1"
)

const chatAppProject = "1007744422436"

func getURI(r *http.Request) string {
	proto := "http"
	host := r.Host
	if protos, ok := r.Header["X-Forwarded-Proto"]; ok {
		proto = protos[0]
	}
	if hosts, ok := r.Header["X-Forwarded-Host"]; ok {
		host = hosts[0]
	}
	return fmt.Sprintf("%s://%s", proto, host)
}

func validateChatToken(tokenString string, chatAppProject string) error {
	context := context.Background()

	jwtURL := "https://www.googleapis.com/service_accounts/v1/jwk/"
	chatIssuer := "chat@system.gserviceaccount.com"
	keySet := oidc.NewRemoteKeySet(context, jwtURL+chatIssuer)
	config := &oidc.Config{
		SkipClientIDCheck: true,
		ClientID:          chatAppProject,
	}
	verifier := oidc.NewVerifier(chatIssuer, keySet, config)
	payload, err := verifier.Verify(context, tokenString)
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

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	uri := getURI(r)
	log.Default().Println("URI: " + uri)
	if reqText, err := httputil.DumpRequest(r, true); err == nil {
		log.Default().Println(string(reqText))
	}

	if authHeader, ok := r.Header["Authorization"]; ok {
		token := strings.Split(authHeader[0], " ")[1]
		if err := validateChatToken(token, chatAppProject); err != nil {
			fmt.Printf("Error validating Token: %v\n", err)
			fmt.Fprintf(w, "Error validating Token: %v\n", err)
			return
		}
	} else {
		// TODO return error to client here.
		fmt.Println("No Auth Header!")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "Error getting name!")
		return
	}
	requestJson := string(bodyBytes)
	name := gjson.Get(requestJson, "message.sender.displayName")
	originalText := gjson.Get(requestJson, "message.text")

	text := fmt.Sprintf(`Hello %s, you said %s`, name, originalText)

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

	encodeAndLogResponse(&resp, w)
}

func encodeAndLogResponse(resp json.Marshaler, w http.ResponseWriter) error {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetIndent("", "  ")
	enc.Encode(resp)
	log.Default().Println("Response: " + b.String())
	w.Write(b.Bytes())
	return nil
}
