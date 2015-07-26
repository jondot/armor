package middleware

import (
	"fmt"
	"net/http"
)

// XXX expand to a ping that also takes a func for custom
// ping validation and custom response (not 'ok')
func Ping(product string, ver string, sha string, host string) http.HandlerFunc {
	const (
		pingver  = "x-ping-version"
		pinghost = "x-ping-host"
		pingfac  = "x-ping-facility"
		pingsha  = "x-ping-sha"
		pingok   = "ok"
	)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(pingver, ver)
		w.Header().Add(pinghost, host)
		w.Header().Add(pingfac, product)
		w.Header().Add(pingsha, sha)
		fmt.Fprint(w, pingok)
	}
}
