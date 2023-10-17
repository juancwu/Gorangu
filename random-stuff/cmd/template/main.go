package main

import (
    "fmt"
    "text/template"
    _ "embed"
    "net/http"
)

//go:embed index.html
var HTML string

func main() {
    tmpl, err := template.New("index.html").Parse(HTML)
    if err != nil {
        fmt.Println("Could not parse index.html: %w", err)
    }

    handler := http.NewServeMux()
    handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        tmpl.Execute(w, map[string]string{
            "dynamic": "this should be dynamic content",
        })
    })

    server := http.Server{Handler: handler, Addr: ":3000"}
    server.ListenAndServe()
}
