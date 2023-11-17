// Sample run-helloworld is a minimal Cloud Run service.
package main

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
	"github.com/gorilla/mux"
	"github.com/rcleveng/assistant/apps"
	"github.com/rcleveng/assistant/cards"
	"github.com/tidwall/gjson"

	"os"
)

func main() {
	cloudRunExecution := os.Getenv("CLOUD_RUN_EXECUTION")
	log.Print("starting server: " + cloudRunExecution)
	router := mux.NewRouter()

	router.HandleFunc("/", handler).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/authorizeFile", authFileHandler).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/chat", chatHandler).Methods(http.MethodPost, http.MethodGet)

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

func mainCard(w http.ResponseWriter, r *http.Request, event *apps.WorkspaceAppEvent) *cards.Card {
	uri := fmt.Sprintf("%s/authorizeFile", getURI(r))

	sampleButton := &cards.Button{
		Text: "Sample Button",
		OnClick: &cards.OnClick{
			Action: &cards.Action{
				Function: uri,
			},
		},
	}
	card := &cards.Card{
		// CardActions: []*cards.CardAction{},
		// Header:      &cards.CardHeader{},
		Name: "Sample Card",
		Sections: []*cards.Section{{
			Header: "",
			Widgets: []*cards.WidgetMarkup{{
				ButtonList:          &cards.ButtonList{Buttons: []*cards.Button{sampleButton}},
				HorizontalAlignment: "CENTER",
				// Image:      &cards.Image{},
				// TextParagraph: &cards.TextParagraph{
				// 	Text: "Hello World",
				// },
			}},
		}},
	}

	if event == nil || event.CommonEventObject == nil {
		return card
	}
	switch event.CommonEventObject.HostApp {
	case "DOCS":
		if event.Docs.Id == nil {
			// Need to auth
			button := &cards.Button{
				Text: "Authorize File",
				OnClick: &cards.OnClick{
					Action: &cards.Action{
						Function: uri,
					},
				},
			}
			card.FixedFooter = &cards.CardFixedFooter{
				PrimaryButton: button,
			}
		}
	}

	return card
}

func authFileHandler(w http.ResponseWriter, r *http.Request) {
	response := cards.SubmitFormResponse{
		RenderAction: &cards.RenderActions{
			HostAppAction: &cards.HostAppAction{
				EditorAction: &cards.EditorClientAction{
					RequestFileScopeForActiveDocument: cards.RequestFileScopeForActiveDocument{},
				},
			},
		},
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(response)
	log.Default().Println("Responrese: " + b.String())
	fmt.Fprintln(w, b.String())
}

func getURI(r *http.Request) string {
	proto := "http"
	if protos, ok := r.Header["X-Forwarded-Proto"]; ok {
		proto = protos[0]
	}
	return fmt.Sprintf("%s://%s", proto, r.Host)
}

func handler(w http.ResponseWriter, r *http.Request) {

	uri := getURI(r)
	log.Default().Println("URI: " + uri)
	if reqText, err := httputil.DumpRequest(r, true); err == nil {
		log.Default().Println(string(reqText))
	}

	event := apps.WorkspaceAppEvent{}
	json.NewDecoder(r.Body).Decode(&event)

	mc := mainCard(w, r, &event)
	action := cards.RenderActions{
		Action: &cards.RenderActionsAction{
			Navigations: &[]cards.Navigation{
				{PushCard: mc}}},
	}

	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetIndent("", "  ")
	enc.Encode(action)
	log.Default().Println("Response: " + b.String())

	json.NewEncoder(w).Encode(action)
}

// CHAT

func chatHandler(w http.ResponseWriter, r *http.Request) {
	uri := getURI(r)
	log.Default().Println("URI: " + uri)
	if reqText, err := httputil.DumpRequest(r, true); err == nil {
		log.Default().Println(string(reqText))
	}

	authHeader := r.Header["Authorization"]
	if len(authHeader) > 0 {
		jwtURL := "https://www.googleapis.com/service_accounts/v1/jwk/"
		chatIssuer := "chat@system.gserviceaccount.com"
		tokenString := strings.Split(authHeader[0], " ")[1]
		audience := "1007744422436"
		context := context.Background()
		keySet := oidc.NewRemoteKeySet(context, jwtURL+chatIssuer)
		config := &oidc.Config{
			SkipClientIDCheck: true,
			ClientID:          audience,
		}
		verifier := oidc.NewVerifier(chatIssuer, keySet, config)
		payload, err := verifier.Verify(context, tokenString)
		if err != nil {
			panic(err)
		}
		var claims struct {
			Aud string `json:"aud"`
			Iss string `json:"iss"`
		}
		if err := payload.Claims(&claims); err != nil {
			panic(err)
		}
		fmt.Printf("\n\nAud: %s; %v\n\n", claims.Aud, claims.Iss)
	} else {
		fmt.Println("No Auth Header!")
		// TODO return here.
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "Error getting name!")
		return
	}
	json := string(bytes)
	name := gjson.Get(json, "message.sender.displayName")
	text := gjson.Get(json, "message.text")
	fmt.Fprintf(w, `
{
	text: "Hello %s, you said; %s"
}`, name, text)

}
