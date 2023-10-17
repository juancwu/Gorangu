package api

import "net/http"

func (h* Handler) getLandingPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(STATUS_OK)
    w.Write([]byte("Sup"))
}
