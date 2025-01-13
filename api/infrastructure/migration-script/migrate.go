package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/maronfranc/poc-golang-ddd/infrastructure"
)

// TODO: add "embed.FS" package
const migration_table_name = "pg_migrations"
const migration_sql_folder = "./migration-script/sql"

func run_migration() error {
	envfile := infrastructure.EnvGetFile()
	err := infrastructure.Start(envfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Infrastructure start failed: %v\n", err)
		panic(err)
	}

	stmt := fmt.Sprintf(
		"SELECT file_id FROM %s ORDER BY created_at DESC LIMIT 1",
		migration_table_name,
	)

	var recent_file_id string
	infrastructure.DbConn.Get(&recent_file_id, stmt)
	log.Printf("Starting migration with id:[%s]", recent_file_id)

	dir_entries, err := os.ReadDir(migration_sql_folder)
	if err != nil {
		return err
	}

	for _, entry := range dir_entries {
		file_id := strings.Split(entry.Name(), ".")[0]

		log.Printf("[LOG] file_id:%s", file_id)
		log.Printf("[LOG] recent_file_id:%s", recent_file_id)

		is_already_migrated := file_id <= recent_file_id
		if is_already_migrated {
			continue
		}

		file_path := fmt.Sprintf("%s/%s", migration_sql_folder, entry.Name())
		buf, err := ioutil.ReadFile(file_path)
		if err != nil {
			return err
		}

		sql_file_content := string(buf)
		_, err = infrastructure.DbConn.Query(sql_file_content)

		if err != nil {
			infrastructure.DbConn.Query("ROLLBACK")
			return err
		}

		log.Printf("Migrated successfully:[%s]", file_id)
	}

	infrastructure.CloseDb()
	return nil
}

func main() {
	err := run_migration()
	if err != nil {
		log.Printf("%v", err)
	}
}
