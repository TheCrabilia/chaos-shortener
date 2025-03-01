package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func errorJSON(w http.ResponseWriter, msg string, code int) {
	h := w.Header()
	h.Del("Content-Length")

	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")

	error := struct {
		Msg string `json:"message"`
	}{
		Msg: msg,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&error)
}

func fullURL(r *http.Request) string {
	schema := "http"
	if r.URL.Scheme != "" {
		schema = r.URL.Scheme
	}

	return fmt.Sprintf("%s://%s", schema, r.Host)
}
