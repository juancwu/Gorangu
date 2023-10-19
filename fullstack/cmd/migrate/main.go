package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
	"github.com/spf13/cobra"
)

//go:embed sql/create_migrations_table.sql
var createMigrationsTableSQL string

//go:embed sql/check_migrations_table.sql
var checkMigrationsTableSQL string

var migrationsPath string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Use:   "go run ./cmd/migrate",
		Short: "A CLI tool for generating migration files",
	}

	var generateCmd = &cobra.Command{
		Use:   "gen",
		Short: "Genereate new migration files",
		Args:  cobra.ExactArgs(1),
		Run:   generate,
	}
	generateCmd.Flags().StringVarP(&migrationsPath, "path", "p", "./", "Path to store the migration files")

	var migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Perform UP migrations",
		Args:  cobra.NoArgs,
		Run:   migrateUp,
	}
	migrateUpCmd.Flags().StringVarP(&migrationsPath, "path", "p", "./", "Path where the migration files are located")

	var migrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Perform DOWN migrations",
		Args:  cobra.NoArgs,
		Run:   migrateDown,
	}
	migrateDownCmd.Flags().StringVarP(&migrationsPath, "path", "p", "./", "Path where the migration files are located")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(migrateUpCmd)
	rootCmd.AddCommand(migrateDownCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

func migrateUp(cmd *cobra.Command, args []string) {
	DB_URL := os.Getenv("DB_URL")
	DB_AUTH_TOKEN := os.Getenv("DB_AUTH_TOKEN")

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", DB_URL, DB_AUTH_TOKEN))
	if err != nil {
		log.Fatalf("Error opening a connection to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// create migration table new table
	var name string
	err = db.QueryRow(checkMigrationsTableSQL).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			// table does not exists, create table
			_, err = db.Exec(createMigrationsTableSQL)
			if err != nil {
				log.Fatalf("Failed to create migrations table: %v", err)
				os.Exit(1)
			}
		}
	}

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_up.sql") {
			var id int32
			err := db.QueryRow("SELECT id FROM migrations WHERE name = ?", file.Name()).Scan(&id)
			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("Failed to query migrations table: %v", err)
				os.Exit(1)
			}

			if err == sql.ErrNoRows {
				fmt.Printf("Applying migration: %s\n", file.Name())

                // apply migrations
                content, err := os.ReadFile(filepath.Join(migrationsPath, file.Name()))
                if err != nil {
                    log.Fatalf("Failed to apply migration from file %s: %v", file.Name(), err)
                    os.Exit(1)
                }

                _, err = db.Exec(string(content))
                if err != nil {
                    log.Fatalf("Failed to apply migration from file %s: %v", file.Name(), err)
                    os.Exit(1)
                }

				_, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", file.Name())
				if err != nil {
					log.Fatalf("Failed to insert into migrations table: %v", err)
					os.Exit(1)
				}
			}
		}
	}

    fmt.Println("Finished applying migrations!")
}

func migrateDown(cmd *cobra.Command, args []string) {
	DB_URL := os.Getenv("DB_URL")
	DB_AUTH_TOKEN := os.Getenv("DB_AUTH_TOKEN")

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", DB_URL, DB_AUTH_TOKEN))
	if err != nil {
		log.Fatalf("Error opening a connection to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// create migration table new table
	var name string
	err = db.QueryRow(checkMigrationsTableSQL).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			// table does not exists, create table
			_, err = db.Exec(createMigrationsTableSQL)
			if err != nil {
				log.Fatalf("Failed to create migrations table: %v", err)
				os.Exit(1)
			}
		}
	}

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_down.sql") {
			var id int32
			err := db.QueryRow("SELECT id FROM migrations WHERE name = ?", file.Name()).Scan(&id)
			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("Failed to query migrations table: %v", err)
				os.Exit(1)
			}

			if err == sql.ErrNoRows {
				fmt.Printf("Reset migration: %s\n", file.Name())

                // apply migrations
                content, err := os.ReadFile(filepath.Join(migrationsPath, file.Name()))
                if err != nil {
                    log.Fatalf("Failed to apply migration from file %s: %v", file.Name(), err)
                    os.Exit(1)
                }

                _, err = db.Exec(string(content))
                if err != nil {
                    log.Fatalf("Failed to apply migration from file %s: %v", file.Name(), err)
                    os.Exit(1)
                }

				_, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", file.Name())
				if err != nil {
					log.Fatalf("Failed to insert into migrations table: %v", err)
					os.Exit(1)
				}
			}
		}
	}

    fmt.Println("Finished applying migrations!")
}

func generate(cmd *cobra.Command, args []string) {
	migrationName := args[0]
	timestamp := time.Now().UTC().Format("20060102150405")

	upMigration := fmt.Sprintf("%s/%s_%s_up.sql", migrationsPath, timestamp, migrationName)
	downMigration := fmt.Sprintf("%s/%s_%s_down.sql", migrationsPath, timestamp, migrationName)

	upFile, err := os.Create(upMigration)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer upFile.Close()

	_, err = upFile.WriteString("-- Write your UP migration SQL here.\n")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	downFile, err := os.Create(downMigration)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer downFile.Close()

	_, err = downFile.WriteString("-- Write your DOWN migration SQL here.\n")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Migration files '%s' and '%s' created successfully.\n", upMigration, downMigration)
}
