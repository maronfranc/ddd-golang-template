package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/maronfranc/poc-golang-ddd/infrastructure"
	"github.com/maronfranc/poc-golang-ddd/infrastructure/database"
)

//go:embed sql/*
var migrationFiles embed.FS

var (
	migration_table_name = "pg_migrations"
	migration_sql_dir    = "sql"
	file_up_pattern      = ".up."
	file_down_pattern    = ".down."
)

func setupDatabase() error {
	envfile, err := infrastructure.EnvGetFileName()
	if err != nil {
		return err
	}

	err = infrastructure.EnvLoad(envfile)
	if err != nil {
		return err
	}

	// Start database connection.
	err = database.Start(envfile)
	if err != nil {
		return err
	}

	return nil
}

func createMigrationTableIfNotExists() error {
	// Check if migration table exists, create if not.
	tableCheckStmt := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = '%s'
		)`, migration_table_name)
	var tableExists bool
	err := database.DbConn.Get(&tableExists, tableCheckStmt)
	if err != nil {
		return err
	}

	if !tableExists {
		createTableStmt := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id SERIAL PRIMARY KEY,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				file_id VARCHAR UNIQUE NOT NULL
			)`, migration_table_name)
		_, err := database.DbConn.Exec(createTableStmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func runMigration() error {
	err := setupDatabase()
	if err != nil {
		return err
	}
	defer database.CloseDb()

	stmt := fmt.Sprintf(
		"SELECT file_id FROM %s ORDER BY file_id DESC LIMIT 1",
		migration_table_name,
	)

	var recent_file_id string
	database.DbConn.Get(&recent_file_id, stmt)
	if recent_file_id == "" {
		log.Printf("No migrations to run.")
	} else {
		log.Printf("Starting migration with id: `%s`", recent_file_id)
	}

	// Get all migration files.
	entries, err := migrationFiles.ReadDir(migration_sql_dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("No migration files found in %s.", migration_sql_dir)
	}

	err = createMigrationTableIfNotExists()
	if err != nil {
		return err
	}

	// Process each migration file.
	for _, entry := range entries {
		file_id := strings.Split(entry.Name(), ".")[0]

		is_already_migrated := file_id <= recent_file_id
		if is_already_migrated {
			continue
		}
		not_valid_up_pattern := !strings.Contains(entry.Name(), file_up_pattern)
		if not_valid_up_pattern {
			continue
		}

		log.Printf("• file_id: %s", file_id)

		// Read from embedded filesystem.
		file_path := fmt.Sprintf("%s/%s", migration_sql_dir, entry.Name())
		buf, err := migrationFiles.ReadFile(file_path)
		if err != nil {
			return err
		}

		tx, err := database.DbConn.Begin()
		if err != nil {
			return err
		}

		sql_file_content := string(buf)
		_, err = tx.Exec(sql_file_content)
		if err != nil {
			tx.Rollback()
			return err
		}

		insertMigrationStmt := fmt.Sprintf("INSERT INTO %s (file_id) VALUES ($1)", migration_table_name)
		_, err = tx.Exec(insertMigrationStmt, file_id)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

// runDownMigration undo the most recent migration.
// TODO: accept `file_id` params to undo many migrations at once.
func runDownMigration() error {
	err := setupDatabase()
	if err != nil {
		return err
	}
	defer database.CloseDb()

	// Get the most recent migration.
	stmt := fmt.Sprintf(
		"SELECT file_id FROM %s ORDER BY file_id DESC LIMIT 1",
		migration_table_name,
	)

	var recent_file_id string
	database.DbConn.Get(&recent_file_id, stmt)
	if recent_file_id == "" {
		log.Printf("No migrations to roll back.")
		return nil
	}

	log.Printf("Rolling back migration with id: `%s`", recent_file_id)

	// Look for files that match the pattern with the recent_file_id.
	entries, err := migrationFiles.ReadDir(migration_sql_dir)
	if err != nil {
		return err
	}

	var downFilePath string
	for _, entry := range entries {
		if strings.Contains(entry.Name(), recent_file_id) &&
			strings.Contains(entry.Name(), file_down_pattern) &&
			strings.HasSuffix(entry.Name(), ".sql") {
			downFilePath = fmt.Sprintf("%s/%s", migration_sql_dir, entry.Name())
			break
		}
	}

	if downFilePath == "" {
		return fmt.Errorf("no down migration file found for id: %s", recent_file_id)
	}

	// Read from embedded filesystem.
	buf, err := migrationFiles.ReadFile(downFilePath)
	if err != nil {
		return err
	}

	tx, err := database.DbConn.Begin()
	if err != nil {
		return err
	}

	// Execute the down SQL.
	sql_file_content := string(buf)
	_, err = tx.Exec(sql_file_content)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the migration record.
	deleteStmt := fmt.Sprintf("DELETE FROM %s WHERE file_id = $1", migration_table_name)
	_, err = tx.Exec(deleteStmt, recent_file_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Printf("Usage: %s [up|down]", os.Args[0])
		os.Exit(1)
	}

	switch args[0] {
	case "up":
		err := runMigration()
		if err != nil {
			panic(err)
		}
		log.Print("All migrations up applied successfully.")
	case "down":
		err := runDownMigration()
		if err != nil {
			panic(err)
		}
		log.Print("All migrations down applied successfully.")
	default:
		log.Printf("Unknown command: %s", args[0])
		log.Printf("Valid args: %s [up|down]", os.Args[0])
		os.Exit(1)
	}
}
