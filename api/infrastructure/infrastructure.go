package infrastructure

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DbConn *sqlx.DB

// func Connect(user, password, dbhost, dbport, dbname string) error {
func Connect(connectionString string) error {
	conn, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	DbConn = conn
	return nil
}

func InsertReturningId(
	table string, s interface{},
) (string, error) {
	keyStr, valuePlaceholder, values := getStructPlaceholder(s)
	iStmt := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES(%s) RETURNING id",
		table,
		keyStr,
		valuePlaceholder,
	)
	var id string
	err := DbConn.Get(&id, iStmt, values...)
	return id, err
}

// UpdateById UPDATE table SET key1 = ? WHERE id = ?;
func UpdateById(table, id string, s any) error {
	setPlaceholder, values, lastIndex := getUpdatePlaceholder(s)
	iStmt := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id=$%d",
		table, setPlaceholder, lastIndex+1)
	values = append(values, id)
	_, err := DbConn.Query(iStmt, values...)
	return err
}

func DeleteById(table, id string) error {
	iStmt := fmt.Sprintf("DELETE FROM %s WHERE id=$1", table)
	_, err := DbConn.Query(iStmt, id)
	return err
}

func SelectNow() string {
	var timePing string
	DbConn.Get(&timePing, "SELECT now()::varchar")
	return timePing
}

func SelectOffset(table string, fields []string, page, perPage int) []interface{} {
	var o []interface{}
	fieldStr := strings.Join(fields, ",")
	offset := (page - 1) * perPage
	iStmt := fmt.Sprintf(
		"SELECT %s FROM %s LIMIT %d OFFSET %d",
		fieldStr, table, perPage, offset)
	DbConn.Get(&o, iStmt)
	return o
}

func SelectById[T any](table, id string, fields []string) *T {
	var o T
	selStr := strings.Join(fields, ",")
	iStmt := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id=$1 LIMIT 1",
		selStr, table)
	DbConn.Get(&o, iStmt, id)
	return &o
}

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
