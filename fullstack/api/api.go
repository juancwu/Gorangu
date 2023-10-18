package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
    router chi.Router
}

const (
    STATUS_OK = 200
    STATUS_CREATED = 201
    STATUS_BAD_REQUEST = 400
    STATUS_UNAUTHORIZED = 401
    STATUS_NOT_FOUND = 404
    STATUS_SERVER_ERROR = 500
)

func New() *Handler {
    h := &Handler{}

    h.router = chi.NewRouter()
    h.router.Use(middleware.Logger)

    h.router.Get("/", h.getLandingPage)

    return h
}

func (h* Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.router.ServeHTTP(w, r)
}
