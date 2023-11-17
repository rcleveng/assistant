package server

import (
	"fmt"
	"net/http"
)

// Gets the webserver's public endpoint (https://foo.com), even if the
// local server is behind a proxy.
func GetPublicEndpoint(r *http.Request) string {
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
