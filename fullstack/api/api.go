package api

import (
	"fmt"
	"net/http"
    "log"

	"github.com/go-chi/chi/v5"

	"github.com/juancwu/Gorangu/fullstack/constants"
	"github.com/juancwu/Gorangu/fullstack/store"
)

func New() chi.Router {
    router := chi.NewRouter()

    router.Get("/", func(w http.ResponseWriter, r *http.Request) {
        db, err := store.New()
        defer db.Close()
        if err != nil {
            log.Fatal(err)
            w.WriteHeader(constants.STATUS_SERVER_ERROR)
            w.Write([]byte(fmt.Sprintf("%s", err)))
            return
        }

        rows, err := db.Query("CREATE TABLE IF NOT EXISTS test_table (id integer primary key autoincrement, name text)")
        if err != nil {
            log.Fatal(err)
            w.WriteHeader(constants.STATUS_BAD_REQUEST)
            w.Write([]byte(fmt.Sprintf("%s", err)))
            return
        }
        defer rows.Close()

        w.WriteHeader(constants.STATUS_OK)
        w.Write([]byte("hello from api router"))
    })

    router.Post("/test", func (w http.ResponseWriter, r *http.Request) {
        r.ParseForm()

        for key, value := range r.Form {
            fmt.Printf("%s = %s\n", key, value)
        }

        w.WriteHeader(constants.STATUS_OK)
        fmt.Fprintln(w, "Post success!")
    })

    return router
}
