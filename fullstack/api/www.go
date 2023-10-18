package api

import (
    "net/http"
    _ "embed"
)

//go:embed html/index.html
var LANDING_PAGE []byte

func (h* Handler) getLandingPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(STATUS_OK)
    w.Write(LANDING_PAGE)
}
