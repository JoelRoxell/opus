package handlers

import (
	"net/http"
)

// Status is the system's health check.
func Status(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
}
