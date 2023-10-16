package main

import (
    "fmt"
    "net/http"
    _ "embed"
    "os"

    "github.com/juancwu/Gorangu/internal/oauth2"
)

//go:embed index.html
var HOME []byte

func main() {
    handler := http.NewServeMux()
    handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        w.Write(HOME)
    })

    server := http.Server{Handler: handler, Addr: ":3000"}
    go func() {
        server.ListenAndServe()
    }()

    err := oauth2.Auth()
    if err != nil {
        fmt.Println("Error: %w", err)
        os.Exit(1)
    }
}
