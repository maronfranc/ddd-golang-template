package database

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/maronfranc/poc-golang-ddd/util/development"
)

func InsertReturningId(table string, entity any) (string, error) {
	stmt, values := PrepareInsertReturningId(table, entity)

	var id string
	err := DbConn.Get(&id, stmt, values...)
	return id, err
}

func Insert(table string, entity any) error {
	keyStr, valuePlaceholder, values := prepareInsertQuery(entity)
	stmt := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES(%s)",
		table,
		keyStr,
		valuePlaceholder,
	)

	_, err := DbConn.Exec(stmt, values...)
	return err
}

// UpdateById
//   - table: database table name.
//   - id: entity id.
//   - entity: object with key and field to be inserted in the database.
//   - dbNullKeys: list of key names with column name to be set to NULL.
func UpdateById(table string, id string, entity any, dbNullKeys []string) error {
	// lastIndex is the placeholder `$N` number.
	setPlaceholder, values, lastIndex := prepareUpdateQuery(entity, dbNullKeys)
	stmt := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id=$%d",
		table, setPlaceholder, lastIndex+1)
	values = append(values, id)

	_, err := DbConn.Query(stmt, values...)
	return err
}

func UpdateWhere(table string, entity any, wvs []WhereValue, dbNullKeys []string) error {
	// lastIndex is the placeholder `$N` number.
	setPlaceholder, values, lastIndex := prepareUpdateQuery(entity, dbNullKeys)
	whereStrFields, whereValues := prepareWhereQuery(wvs, lastIndex)
	stmt := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table, setPlaceholder, whereStrFields)
	values = append(values, whereValues...)

	_, err := DbConn.Query(stmt, values...)
	return err
}

func PrepareInsertReturningId(table string, entity any) (string, []any) {
	keyStr, valuePlaceholder, values := prepareInsertQuery(entity)
	stmt := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES(%s) RETURNING \"id\"",
		table,
		keyStr,
		valuePlaceholder,
	)
	return stmt, values
}

func prepareInsertQuery(entity any) (string, string, []any) {
	// dbKeyStr `id,key_1,key_2`
	var dbKeyStr string
	// placeholder `$1,$2,$3`
	var placeholder string
	// values `[]interface{1, true, "AnyValue"}`
	var values []any

	// Handle pointer input
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Pointer {
		entityValue = entityValue.Elem()
		if !entityValue.IsValid() {
			return "", "", nil
		}
	}

	entityType := entityValue.Type()

	// Iterate over struct fields.
	for i := 0; i < entityValue.NumField(); i++ {
		field := entityType.Field(i)
		dbTag := field.Tag.Get("db")

		// Skip if no or empty `db` tag.
		if dbTag == "" {
			continue
		}

		// Process the `db` tag to get the actual field name.
		dbFieldName := strings.Split(dbTag, ",")[0]
		if dbFieldName == "" {
			continue
		}

		// Get the actual field value.
		fieldValue := entityValue.Field(i)

		// Skip fields that are nil pointers.
		if fieldValue.Kind() == reflect.Pointer && fieldValue.IsNil() {
			continue
		}

		if dbKeyStr != "" {
			dbKeyStr += ","
		}
		dbKeyStr += dbFieldName

		values = append(values, fieldValue.Interface())
	}

	// Build placeholder string
	placeholderCount := len(values)
	if placeholderCount > 0 {
		placeholders := make([]string, placeholderCount)
		for i := range placeholderCount {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
		placeholder = strings.Join(placeholders, ",")
	}

	// Remove leading comma.
	dbKeyStr = strings.TrimPrefix(dbKeyStr, ",")

	return dbKeyStr, placeholder, values
}

func prepareUpdateQuery(entity any, dbNullKeys []string) (string, []any, int) {
	// keyValueStr `key_1=$1,key_2=$2`
	var keyValueStr string
	// values `[]interface{1, true, "AnyValue"}`
	var values []any

	// Handle both struct and pointer to struct
	elem := reflect.ValueOf(entity)
	if elem.Kind() == reflect.Pointer {
		if elem.IsNil() {
			panic("entity cannot be nil pointer")
		}
		elem = elem.Elem()
	}

	// Ensure we're working with a struct
	if elem.Kind() != reflect.Struct {
		panic("entity must be a struct or pointer to struct")
	}

	lastIndex := 0
	for i := 0; i < elem.NumField(); i++ {
		keyField := elem.Type().Field(i)
		valueField := elem.Field(i)

		value := valueField.Interface()
		// dbKey is in the dto.StructName{Field `db:"key_name"`}.
		dbKey := keyField.Tag.Get("db")
		// jsonKey is in the dto.StructName{Field `json:"key_name,omitempty"`}.
		jsonKey := keyField.Tag.Get("json")

		development.Assert(dbKey != "", "Dto with key field error. Add `db:\"column_name\"`")

		isPrtNil := valueField.Kind() == reflect.Pointer && valueField.IsNil()
		// isEmpty is what define omit the values with nil.
		isEmpty := isPrtNil && strings.Contains(jsonKey, "omitempty")
		if isEmpty {
			continue
		}

		if lastIndex == 0 {
			keyValueStr = fmt.Sprintf("%s=$%d", dbKey, lastIndex+1)
		} else {
			keyValueStr += fmt.Sprintf(",%s=$%d", dbKey, lastIndex+1)
		}

		values = append(values, value)
		lastIndex = lastIndex + 1
	}

	for _, dbKey := range dbNullKeys {
		if lastIndex == 0 {
			keyValueStr = fmt.Sprintf("%s=$%d", dbKey, lastIndex+1)
		} else {
			keyValueStr += fmt.Sprintf(",%s=$%d", dbKey, lastIndex+1)
		}
		values = append(values, nil)
		lastIndex = lastIndex + 1
	}

	return keyValueStr, values, lastIndex
}

type WhereValue struct {
	Field string
	Value any
}

// prepareWhereQuery creates a WHERE clause query string and corresponding values.
//
// Parameters:
//   - wvs: slice of WhereValue containing field names and their corresponding values
//   - lastIndex: last query placeholder($2) index to start next at ($lastIndex+1).
// Example:
//	- Input: []WhereValue{{Field: "name", Value: "John"}, {Field: "age", Value: 30}}`
//	- Output: ("name=$1,age=$2", []any{"John", 30})
func prepareWhereQuery(wvs []WhereValue, lastIndex int) (string, []any) {
	var values []any
	var strs []string
	paramIndex := lastIndex + 1
	for _, where := range wvs {
		development.Assert(where.Field != "", "Please provide the repository table column")
		q := fmt.Sprintf("%s=$%d", where.Field, paramIndex)
		strs = append(strs, q)
		values = append(values, where.Value)
		paramIndex++
	}

	query := strings.Join(strs, ",")
	return query, values
}

func Count(table string) (int, error) {
	stmt := fmt.Sprintf("SELECT count(*) AS total FROM %s", table)
	var count int
	err := DbConn.Get(&count, stmt)
	return count, err
}

func SelectManyAndCount[T any](table string, page, perPage int) ([]T, *int, error) {
	fields := getStructTags[T]("db")
	fieldStr := strings.Join(fields, ",")

	offset := (page - 1) * perPage
	stmt := fmt.Sprintf(
		"SELECT %s FROM %s LIMIT %d OFFSET %d",
		fieldStr, table, perPage, offset)
	pgStmt := fmt.Sprintf("SELECT count(*) AS total FROM %s", table)

	var vs []T
	var total int
	var err1, err2 error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		err1 = DbConn.Select(&vs, stmt)
	}()

	go func() {
		defer wg.Done()
		err2 = DbConn.Get(&total, pgStmt)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, nil, err1
	}
	if err2 != nil {
		return nil, nil, err2
	}

	return vs, &total, nil
}

func prepareSelectWhere(table string, fields []string, wvs []WhereValue) (string, []any) {
	fieldStr := strings.Join(fields, ",")
	stmt := fmt.Sprintf(
		"SELECT %s FROM %s ",
		fieldStr, table)

	var values []any
	if len(wvs) > 0 {
		whereQ, vs := prepareWhereQuery(wvs, 0)
		stmt += fmt.Sprintf(" WHERE %s", whereQ)
		values = vs
	}

	return stmt, values
}

// SelectById
// Return: `nil, nil` when item is not found.
func SelectById[T any](table string, id string) (*T, error) {
	fields := getStructTags[T]("db")
	fieldStr := strings.Join(fields, ",")
	stmt := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1 LIMIT 1",
		fieldStr, table)

	var v T
	err := DbConn.Get(&v, stmt, id)
	if err != nil {
		isNotFound := err.Error() == "sql: no rows in result set"
		if isNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &v, nil
}

func SelectOne[T any](table string, wvs []WhereValue) (*T, error) {
	fields := getStructTags[T]("db")
	stmt, values := prepareSelectWhere(table, fields, wvs)
	stmt += " LIMIT 1"

	var v T
	err := DbConn.Get(&v, stmt, values...)
	if err != nil {
		isNotFound := err.Error() == "sql: no rows in result set"
		if isNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &v, nil
}

func DeleteById(table string, id string) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id=$1", table)
	_, err := DbConn.Query(stmt, id)
	return err
}

// Function to get tag values from a struct
//   - For example a struct with `db:"id"` `db:"email"` return `[id email]`
func getStructTags[T any](tagName string) []string {
	var tags []string

	// Get the reflection value of the struct type (not an instance)
	var t = reflect.TypeFor[T]()
	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		development.Assert(tag != "", "Dto tag is empty. Add a `db` tag in the struct")
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}
