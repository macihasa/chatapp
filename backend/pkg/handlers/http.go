package handlers

import (
	"fmt"
	"net/http"
)

// HTTP HANDLERS
func Landing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is live!", r.RemoteAddr)
}
