package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/juancwu/Gorangu/fullstack/api"
	"github.com/juancwu/Gorangu/fullstack/www"
)

func main() {
    router := chi.NewRouter()
    router.Use(middleware.Logger)

    apiRouter := api.New()
    wwwRouter := www.New()

    router.Mount("/", wwwRouter)
    router.Mount("/api", apiRouter)

    fmt.Printf("Starting server to listen on http://localhost:%d\n", 3000)
    if err := http.ListenAndServe(":3000", router); err != nil {
        fmt.Println("Error starting server: %w\n", err)
    }
}
