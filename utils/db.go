package utils

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
)

// Shamelessly copied from chatgpt
func ScanRowToStruct(row pgx.Rows, structPtr interface{}) error {
	// Ensure the structPtr is a pointer
	ptrVal := reflect.ValueOf(structPtr)
	if ptrVal.Kind() != reflect.Ptr || ptrVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, got %T", structPtr)
	}

	// Get the struct fields (a slice of reflect.StructField)
	structType := ptrVal.Elem().Type()
	numFields := structType.NumField()

	// Create a slice to hold pointers to the struct fields where the values will be stored
	destValues := make([]interface{}, numFields)

	// Get the column names
	columns := row.FieldDescriptions()
	columnNames := make([]string, len(columns))
	for i, col := range columns {
		columnNames[i] = string(col.Name)
	}

	// Map column names to struct fields (assuming the struct fields have the same name as the columns)
	for i := 0; i < numFields; i++ {
		structField := structType.Field(i)
		// You could handle snake_case to camelCase conversion here if needed
		fieldName := structField.Name

		// Find the index of the column in the result set that matches the struct field
		var columnIndex int = -1
		for j, colName := range columnNames {
			if colName == fieldName {
				columnIndex = j
				break
			}
		}

		// If we find a match, prepare the struct field for scanning
		if columnIndex >= 0 {
			destValues[i] = reflect.ValueOf(structPtr).Elem().Field(i).Addr().Interface()
		} else {
			// If no matching column is found, set the field to zero value (optional)
			destValues[i] = reflect.New(structField.Type).Interface()
		}
	}

	// Scan the row into the struct
	if err := row.Scan(destValues...); err != nil {
		return fmt.Errorf("failed to scan row into struct: %v", err)
	}

	return nil
}
