package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/juancwu/Gorangu/fullstack/api"
)

func main() {
    h := api.New()

    router := chi.NewRouter()
    router.Use(middleware.Logger)

    router.Mount("/", h)

    fmt.Printf("Starting server to listen on port %d\n", 3000)
    if err := http.ListenAndServe(":3000", router); err != nil {
        fmt.Println("Error starting server: %w\n", err)
    }
}
