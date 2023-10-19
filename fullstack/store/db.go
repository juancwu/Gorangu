package store

import (
	"database/sql"
    "fmt"

	_ "github.com/libsql/libsql-client-go/libsql"

	"github.com/juancwu/Gorangu/fullstack/env"
)

func New() (*sql.DB, error) {
    db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", env.DB_URL, env.DB_AUTH_TOKEN))
    if err != nil {
        return nil, fmt.Errorf("Error opening a connection to database: %w", err)
    }

    return db, nil
}
