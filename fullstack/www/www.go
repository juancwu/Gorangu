package www 

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"

    "github.com/juancwu/Gorangu/fullstack/constants"
)

//go:embed html/index.html
var LANDING_PAGE []byte

func New() chi.Router {
    router := chi.NewRouter()

    router.Get("/", getLandingPage)

    return router
}

func getLandingPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(constants.STATUS_OK)
    w.Write(LANDING_PAGE)
}
