package infrastructure

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DbConn *sqlx.DB

func Start(envfile string) error {
	EnvLoad(envfile)
	connStr, err := EnvGet("PG_CONNECTION_STR")
	if err != nil {
		return err
	}
	connectDb(connStr)
	return nil
}
func connectDb(connectionString string) error {
	conn, err := sqlx.Open("postgres", connectionString)
	DbConn = conn
	return err
}
func CloseDb() {
	log.Print("Closing infrastructure connection")
	DbConn.Close()
}

func InsertReturningId(table string, s interface{}) (string, error) {
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
func UpdateById(table, id string, s any) error {
	setPlaceholder, values, lastIndex := getUpdatePlaceholder(s)
	// UpdateById UPDATE table SET key1 = ? WHERE id = ?;
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
func SelectById[T any](table, id string, fields []string) (*T, error) {
	var o T
	selStr := strings.Join(fields, ",")
	stmt := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id=$1 LIMIT 1",
		selStr, table)
	err := DbConn.Get(&o, stmt, id)
	return &o, err
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
func getUpdatePlaceholder(s any) (string, []any, int) {
	var keyValueStr string
	var values []any
	var lastIndex int
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
	// getUpdatePlaceholder (key1=?,key2=?,key3=?) []any
	return keyValueStr, values, lastIndex
}

func EnvGetFile() string {
	n := flag.String("env", "test", "-env flags: ['test','dev','prod']")
	flag.Parse()
	fileName := fmt.Sprintf(".%s.env", *n)
	return fileName
}
func EnvLoad(name string) error {
	// TODO: https://github.com/joho/godotenv/issues/126#issuecomment-1474645022
	return godotenv.Load(name)
}
func EnvGet(k string) (string, error) {
	value, exists := os.LookupEnv(k)
	if !exists {
		msg := fmt.Sprintf("Key(%s) not in environment", k)
		return "", errors.New(msg)
	}
	return value, nil
}

func EnvGetAsBool(k string) (bool, error) {
	vStr, err := EnvGet(k)
	if err != nil {
		return false, err
	}
	vBool, err := strconv.ParseBool(vStr)
	return vBool, err
}
