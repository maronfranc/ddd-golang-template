package infrastructure

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DbConn *sqlx.DB

func ConnectDb(connectionString string) error {
	conn, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	DbConn = conn
	return nil
}
func CloseDb() {
	log.Print("Closing infrastructure connection")
	DbConn.Close()
}

func InsertReturningId(
	table string, s interface{},
) (string, error) {
	keyStr, valuePlaceholder, values := getStructPlaceholder(s)
	stmt := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES(%s) RETURNING id",
		table,
		keyStr,
		valuePlaceholder,
	)
	var id string
	err := DbConn.Get(&id, stmt, values...)
	return id, err
}

// UpdateById UPDATE table SET key1 = ? WHERE id = ?;
func UpdateById(table, id string, s any) error {
	setPlaceholder, values, lastIndex := getUpdatePlaceholder(s)
	stmt := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id=$%d",
		table, setPlaceholder, lastIndex+1)
	values = append(values, id)
	_, err := DbConn.Query(stmt, values...)
	return err
}

func DeleteById(table, id string) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id=$1", table)
	_, err := DbConn.Query(stmt, id)
	return err
}

func SelectNow() string {
	var timePing string
	DbConn.Get(&timePing, "SELECT now()::varchar")
	return timePing
}

// func SelectPagination[T any](table string, s T, page, perPage int) ([]T, int) {
func SelectPagination[T any](table string, fields []string, page, perPage int) ([]T, int) {
	fieldStr := strings.Join(fields, ",")
	offset := (page - 1) * perPage
	stmt := fmt.Sprintf(
		"SELECT %s FROM %s LIMIT %d OFFSET %d",
		fieldStr, table, perPage, offset)
	pgStmt := fmt.Sprintf("SELECT count(*) AS total FROM %s", table)

	var vs []T
	var total int
	DbConn.Select(&vs, stmt)
	DbConn.Get(&total, pgStmt)
	return vs, total
}

func SelectById[T any](table, id string, fields []string) *T {
	var o T
	selStr := strings.Join(fields, ",")
	stmt := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id=$1 LIMIT 1",
		selStr, table)
	DbConn.Get(&o, stmt, id)
	return &o
}

// func SelectByIdJoin[T any](table, id string, fields []string, tableToJoin string) *T {
// 	var o T
// 	selStr := strings.Join(fields, ",")
// 	stmt := fmt.Sprintf(
// 		"SELECT %s FROM %s t1"+
// 			"INNER JOIN %s t2 ON t1.id = t2." +
// 			" WHERE id=$1 LIMIT 1",
// 		selStr, table, tableToJoin)
// 	// stmt += "WHERE id=$1 LIMIT 1"
// 	DbConn.Get(&o, stmt, id)
// 	return &o
// }

func connectionString(user, password, dbhost, dbport, dbname string) string {
	// "postgres://user:password@dbhost:dbport/my_db"
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, dbhost, dbport, dbname)
}

func GetConnValues() string {
	const dbhost = "localhost"
	const dbport = "5432"
	const user = "my_user"
	const password = "pass123"
	const dbname = "my_db"
	return connectionString(user, password, dbhost, dbport, dbname)
}

// getStructPlaceholder (key1,key2,key3) (?,?,?) and (value1,value2,value3)
func getStructPlaceholder(s any) (string, string, []any) {
	var keyValueStr string
	var placeholder string
	var values []any
	elem := reflect.ValueOf(s).Elem()
	for i := 0; i < elem.NumField(); i++ {
		key := strings.ToLower(elem.Type().Field(i).Name)
		value := elem.Field(i).Interface()
		if i == 0 {
			keyValueStr = key
			placeholder = fmt.Sprintf("$%d", i+1)
		} else {
			keyValueStr += fmt.Sprintf(",%s", key)
			placeholder += fmt.Sprintf(",$%d", i+1)
		}
		values = append(values, value)
	}
	return keyValueStr, placeholder, values
}

// getUpdatePlaceholder (key1=?,key2=?,key3=?) []any
func getUpdatePlaceholder(s any) (string, []any, int) {
	var keyValueStr string
	var values []any
	var lastIndex int
	// UpdateById UPDATE table SET key1=? WHERE id=?;
	elem := reflect.ValueOf(s).Elem()
	for i := 0; i < elem.NumField(); i++ {
		key := strings.ToLower(elem.Type().Field(i).Name)
		value := elem.Field(i).Interface()
		if i == 0 {
			keyValueStr = fmt.Sprintf("%s=$%d", key, i+1)
		} else {
			keyValueStr += fmt.Sprintf(",%s=$%d", key, i+1)
		}
		values = append(values, value)
		lastIndex = i + 1
	}
	return keyValueStr, values, lastIndex
}

// getStructKeys (key1,key2,keyN)
// func getStructKeys(s any) string {
// 	var keyValueStr string
// 	elem := reflect.ValueOf(s).Elem()
// 	for i := 0; i < elem.NumField(); i++ {
// 		key := strings.ToLower(elem.Type().Field(i).Name)
// 		if i == 0 {
// 			keyValueStr = key
// 		} else {
// 			keyValueStr += fmt.Sprintf(",%s", key)
// 		}
// 	}
// 	return keyValueStr
// }

// // getStructFieldNames (key1,key2) and (:key1,:key2)
// func getStructFieldNames(s interface{}) (string, string) {
// 	var fStr, commaStr string
// 	elem := reflect.ValueOf(s).Elem()
// 	for i := 0; i < elem.NumField(); i++ {
// 		key := strings.ToLower(elem.Type().Field(i).Name)
// 		if i == 0 {
// 			fStr = key
// 			commaStr = fmt.Sprintf(":%s", key)
// 		} else {
// 			fStr += fmt.Sprintf(",%s", key)
// 			commaStr += fmt.Sprintf(",:%s", key)
// 		}
// 	}
// 	return fStr, commaStr
// }

// // getStructKeyValue key1=:key1,key2=:key2
// func getStructKeyValue(s any) (string, []any) {
// 	var keyValueStr string
// 	var values []any
// 	elem := reflect.ValueOf(s).Elem()
// 	for i := 0; i < elem.NumField(); i++ {
// 		key := strings.ToLower(elem.Type().Field(i).Name)
// 		value := elem.Field(i).Interface()
// 		if i == 0 {
// 			keyValueStr = fmt.Sprintf("%s=:%s", key, key)
// 		} else {
// 			keyValueStr += fmt.Sprintf(",%s=:%s", key, key)
// 		}
// 		values = append(values, value)
// 	}
// 	return keyValueStr, values
// }
