// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"

	"os"
	"text/template"
)

func main() {
	log.Print("starting server...")
	router := mux.NewRouter()

	router.HandleFunc("/", handler).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/authorizeFile", authFileHandler).Methods(http.MethodPost, http.MethodGet)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func mainCard(w http.ResponseWriter, r *http.Request) string {
	s := `
	{
		"sections": [
		  {
			"widgets": [
			  {
				"buttonList": {
				  "buttons": [
					{
					  "text": "Button 1",
					  "onClick": {
						"action": {
						  "function": "TODO",
						  "parameters": []
						}
					  }
					}
				  ]
				},
				"horizontalAlignment": "END"
			  }
			]
		  }
		]
	  }
`
	return s
}

type RenderActionOneCard struct {
	Card string
}

func authFileHandler(w http.ResponseWriter, r *http.Request) {
	if reqText, err := httputil.DumpRequest(r, true); err == nil {
		log.Default().Println(string(reqText))
	}
	text := `
	{
		renderActions: {
		  hostAppAction: {
			editorAction: {
			  requestFileScopeForActiveDocument: {},
			 },
		  },
		},
	  }	`
	fmt.Fprintln(w, text)
}

func handler(w http.ResponseWriter, r *http.Request) {

	uri := r.RequestURI
	log.Default().Println("URI: " + uri)
	if reqText, err := httputil.DumpRequest(r, true); err == nil {
		log.Default().Println(string(reqText))
	}

	var body map[string]interface{}
	json.NewDecoder(r.Body).Decode(&body)

	if d, ok := body["foo"]; ok {
		docs := d.(map[string]interface{})
		if len(docs) == 0 {

		}
	}

	text := `
	{
        action: {
            navigations: [{
                pushCard: {{ .Card}}
            }]
        }
    }
	`
	t, err := template.New("main").Parse(text)

	if err != nil {
		panic(err)
	}

	mainCard := mainCard(w, r)
	ra := RenderActionOneCard{mainCard}
	err = t.Execute(w, ra)
	if err != nil {
		panic(err)
	}

}
