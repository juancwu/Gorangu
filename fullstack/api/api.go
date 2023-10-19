package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

    "github.com/juancwu/Gorangu/fullstack/constants"
)

func New() chi.Router {
    router := chi.NewRouter()

    router.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(constants.STATUS_OK)
        w.Write([]byte("hello from api router"))
    })

    return router
}
