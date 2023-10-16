package oauth2

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
)

// go:embed login.html
var LOGIN_HTML string

const SERVER_URL string = "http://localhost:3000"

func Auth() error {
    ch := make(chan string, 1)
    server, err := createCallbackServer(ch)
    if err != nil {
        return fmt.Errorf("Could not create callback server: %w", err)
    }

    port, err := runServer(server)
    if err != nil {
        return fmt.Errorf("Could not bind callback server to port: %w", err)
    }

    authUrl, err := url.Parse(SERVER_URL)
    if err != nil {
        return fmt.Errorf("Error parsing auth url: %w", err)
    }

    authUrl.RawQuery = url.Values{
        "port": {strconv.Itoa(port)},
        "redirect": {"true"},
    }.Encode()

    url := authUrl.String()

    fmt.Println("Visit this URL on this device to log in:")
    fmt.Println(url)
    fmt.Println("Waiting for authentication...")

    token := <- ch
    username := <- ch

    server.Shutdown(context.Background())

    fmt.Printf("Success! Logged in as %s\n", username)
    fmt.Printf("Token: %s\n", token)

    return nil
}

func createCallbackServer(ch chan string) (*http.Server, error) {
    tmpl, err := template.New("login.html").Parse(LOGIN_HTML)
    if err != nil {
        return nil, fmt.Errorf("Could not parse login.html template: %w", err)
    }

    handler := http.NewServeMux()
    handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        ch <- query.Get("token")
        ch <- query.Get("username")

        w.WriteHeader(200)
        tmpl.Execute(w, map[string]string{
            "dynamic": "this is dynamic content",
        })
    })

    return &http.Server{Handler: handler}, nil
}

func runServer(server *http.Server) (int, error) {
    // :0 means dynamically assign an available port
    listener, err := net.Listen("tcp", ":0")
    if err != nil {
        return 0, fmt.Errorf("Could not allocate port for http server: %w", err)
    }

    go func() {
        server.Serve(listener)
    }()

    return listener.Addr().(*net.TCPAddr).Port, nil
}
