package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/maronfranc/poc-golang-ddd/infrastructure"
	"github.com/maronfranc/poc-golang-ddd/infrastructure/database"
)

const migration_table_name = "pg_migrations"
const migration_sql_dir = "./infrastructure/migration-script/sql"
const command_pattern = ".up."

func create_migration_table_if_not_exists() error {
	// Check if migration table exists, create if not.
	tableCheckStmt := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = '%s'
		)`, migration_table_name)

	var tableExists bool
	database.DbConn.Get(&tableExists, tableCheckStmt)
	if !tableExists {
		createTableStmt := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id SERIAL PRIMARY KEY,
				created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				file_id VARCHAR UNIQUE NOT NULL
			)`, migration_table_name)
		_, err := database.DbConn.Exec(createTableStmt)
		if err != nil {
			return fmt.Errorf("failed to create migration table: %w", err)
		}
	}

	return nil
}

func run_migration() error {
	// TODO: add "embed.FS" package.
	envfile, err := infrastructure.EnvGetFileName()
	err = infrastructure.EnvLoad(envfile)
	if err != nil {
		panic(err)
	}

	err = database.Start(envfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Infrastructure start failed: %v\n", err)
		panic(err)
	}

	stmt := fmt.Sprintf(
		"SELECT file_id FROM %s ORDER BY file_id DESC LIMIT 1",
		migration_table_name,
	)

	var recent_file_id string
	database.DbConn.Get(&recent_file_id, stmt)
	if recent_file_id == "" {
		log.Printf("Starting migration from the first file.")
	} else {
		log.Printf("Starting migration with id: `%s`", recent_file_id)
	}

	dir_entries, err := os.ReadDir(migration_sql_dir)
	if err != nil {
		return err
	}
	if len(dir_entries) == 0 {
		return fmt.Errorf("No migration files found in %s.", migration_sql_dir)
	}

	create_migration_table_if_not_exists()

	for _, entry := range dir_entries {
		file_id := strings.Split(entry.Name(), ".")[0]

		is_already_migrated := file_id <= recent_file_id
		if is_already_migrated {
			continue
		}
		not_valid_up_pattern := !strings.Contains(entry.Name(), command_pattern)
		if not_valid_up_pattern {
			continue
		}

		log.Printf("[LOG] file_id: %s", file_id)

		file_path := fmt.Sprintf("%s/%s", migration_sql_dir, entry.Name())
		buf, err := os.ReadFile(file_path)
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

		// Insert migration record.
		insertStmt := fmt.Sprintf("INSERT INTO %s (file_id) VALUES ($1)", migration_table_name)
		_, err = tx.Exec(insertStmt, file_id)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	database.CloseDb()
	return nil
}

func main() {
	err := run_migration()
	if err != nil {
		panic(err)
	}

	log.Print("All migrations applied successfully!")
}
