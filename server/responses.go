package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Encodes the resource as JSON, logs to locally and writes it to the
// response w.
func EncodeAndLogResponse(resp json.Marshaler, w http.ResponseWriter) error {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetIndent("", "  ")
	enc.Encode(resp)
	log.Default().Println("Response: " + b.String())
	w.Write(b.Bytes())
	return nil
}
