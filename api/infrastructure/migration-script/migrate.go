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
const migration_sql_folder = "./infrastructure/migration-script/sql"

func run_migration() error {
	// TODO: add "embed.FS" package
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
		"SELECT file_id FROM %s ORDER BY created_at DESC LIMIT 1",
		migration_table_name,
	)

	var recent_file_id string
	database.DbConn.Get(&recent_file_id, stmt)
	log.Printf("Starting migration with id:[%s]", recent_file_id)

	dir_entries, err := os.ReadDir(migration_sql_folder)
	if err != nil {
		return err
	}

	for _, entry := range dir_entries {
		file_id := strings.Split(entry.Name(), ".")[0]
		log.Printf("[LOG] file_id:%s", file_id)

		is_already_migrated := file_id <= recent_file_id
		if is_already_migrated {
			continue
		}

		file_path := fmt.Sprintf("%s/%s", migration_sql_folder, entry.Name())
		// buf, err := ioutil.ReadFile(file_path)
		buf, err := os.ReadFile(file_path)
		if err != nil {
			return err
		}

		sql_file_content := string(buf)
		_, err = database.DbConn.Query(sql_file_content)

		if err != nil {
			database.DbConn.Query("ROLLBACK")
			return err
		}

		log.Printf("Migrated successfully:[%s]", file_id)
	}

	database.CloseDb()
	return nil
}

func main() {
	err := run_migration()
	if err != nil {
		log.Printf("%v", err)
	}
}
