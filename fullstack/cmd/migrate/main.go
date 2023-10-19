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
    "sort"

	"github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
	"github.com/spf13/cobra"
)

//go:embed sql/create_migrations_table.sql
var createMigrationsTableSQL string

//go:embed sql/check_migrations_table.sql
var checkMigrationsTableSQL string

//go:embed sql/check_migration.sql
var checkMigrationSQL string

//go:embed sql/record_migration.sql
var recordMigrationSQL string

//go:embed sql/delete_migration.sql
var deleteMigrationSQL string

var migrationsPath string
var dbURL string
var dbAuthToken string

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

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Performs migrations",
		Args:  cobra.NoArgs,
		Run:   migrate,
	}
	migrateCmd.Flags().StringVarP(&migrationsPath, "path", "p", "./", "Path where the migration files are located")
	migrateCmd.Flags().StringVarP(&dbURL, "url", "u", "", "The Database URL. If env DB_URL is defined then this is not needed")
	migrateCmd.Flags().StringVarP(&dbAuthToken, "token", "t", "", "The Database auth token. If env DB_AUTH_TOKEN is defined then this is not needed")

	var rollbackCmd = &cobra.Command{
		Use:   "rollback",
		Short: "Rollbacks migrations",
		Args:  cobra.NoArgs,
		Run:   rollback,
	}
	rollbackCmd.Flags().StringVarP(&migrationsPath, "path", "p", "./", "Path where the migration files are located")
	rollbackCmd.Flags().StringVarP(&dbURL, "url", "u", "", "The Database URL. If env DB_URL is defined then this is not needed")
	rollbackCmd.Flags().StringVarP(&dbAuthToken, "token", "t", "", "The Database auth token. If env DB_AUTH_TOKEN is defined then this is not needed")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(migrateCmd)
    rootCmd.AddCommand(rollbackCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

func migrate(cmd *cobra.Command, args []string) {
    checkEnv(&dbURL, &dbAuthToken)

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbURL, dbAuthToken))
	if err != nil {
		log.Fatalf("Error opening a connection to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

    err = checkMigrationTable(db)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

    sort.Slice(files, func(i, j int) bool {
        // ascending order
        return files[i].Name() < files[j].Name()
    })

	var suffix = "_up.sql"
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			baseName := strings.TrimSuffix(file.Name(), suffix)
			var id int32
			err := db.QueryRow(checkMigrationSQL, baseName).Scan(&id)
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

                _, err = db.Exec(recordMigrationSQL, baseName)
                if err != nil {
                    log.Fatalf("Failed to insert migration record into migrations table: %v", err)
                    os.Exit(1)
                }
			}
		}
	}

	fmt.Println("Finished applying migrations ✅")
}

func rollback(cmd *cobra.Command, args []string) {
    checkEnv(&dbURL, &dbAuthToken)

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbURL, dbAuthToken))
	if err != nil {
		log.Fatalf("Error opening a connection to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

    err = checkMigrationTable(db)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

    sort.Slice(files, func(i, j int) bool {
        // descending order
        return files[i].Name() > files[j].Name()
    })

    suffix := "_down.sql"
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			baseName := strings.TrimSuffix(file.Name(), suffix)
			var id int32
			err := db.QueryRow(checkMigrationSQL, baseName).Scan(&id)
			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("Failed to query migrations table: %v", err)
				os.Exit(1)
			}

            if err == nil && id > 0 {
				fmt.Printf("Resetting migration: %s\n", file.Name())

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
                _, err = db.Exec(deleteMigrationSQL, id)
                if err != nil {
                    log.Fatalf("Failed to remove migration record from migrations table: %v", err)
                    os.Exit(1)
                }
            }
		}
	}

	fmt.Println("Finished rolling back migrations ✅")
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

func checkMigrationTable(db *sql.DB) error {
	// check if migrations table exists or not
    _, err := db.Exec(checkMigrationsTableSQL)
	if err != nil {
		if err == sql.ErrNoRows {
			// table does not exists, create table
			_, err = db.Exec(createMigrationsTableSQL)
			if err != nil {
				return fmt.Errorf("Failed to create migrations table: %v", err)
			}
		}
	}

    return nil
}

func checkEnv(dbURL, dbAuthToken *string) {
	if *dbURL == "" {
		*dbURL = os.Getenv("DB_URL")
	}

	if *dbAuthToken == "" {
		*dbAuthToken = os.Getenv("DB_AUTH_TOKEN")
	}
}
