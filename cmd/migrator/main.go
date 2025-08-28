package main

import (
    "errors"
    "flag"
    "fmt"
    "log"

    "github.com/golang-migrate/migrate/v4"
    _"github.com/golang-migrate/migrate/v4/database/sqlite3"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    var storagePath, migrationsPath, migrationsTable string

    flag.StringVar(&storagePath, "storage-path", "", "Path to the storage file")
    flag.StringVar(&migrationsPath, "migrations-path", "", "Path to the migrations folder")
    flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Name of the migrations table")
    flag.Parse()

    if storagePath == "" {
        panic("storage-path is required")
    }

    if migrationsPath == "" {
        panic("migrations-path is required")
    }

    // Убедитесь, что не добавляете "file://" к storagePath
    // SQLite ожидает просто путь к файлу
    databaseURL := fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable)
    
    m, err := migrate.New(
        "file://"+migrationsPath,
        databaseURL,
    )
    if err != nil {
        log.Fatalf("Failed to create migrate instance: %v", err)
    }

    if err := m.Up(); err != nil {
        if errors.Is(err, migrate.ErrNoChange) {
            fmt.Println("No new migrations to apply")
            return
        }
        log.Fatalf("Failed to apply migrations: %v", err)
    }

    fmt.Println("Migrations applied successfully")
}