// EXPEREMENTAL
package btrsql

import (
	"database/sql"
	"fmt"
	"reflect"
)

func ScanToStruct(rows *sql.Rows, structPtr interface{}) error {
	structType := reflect.TypeOf(structPtr)

	mappings, err := mapColumns(rows, structType)
	if err != nil {
		return err
	}

	err = rows.Scan(mappings.mappedFieldPtrs...)
	if err != nil {
		return err
	}

	mappings.setFields(structPtr)

	return nil
}

type structMappings struct {
	mappedFields    []*field
	mappedFieldPtrs []interface{}
}

func (s *structMappings) setFields(destPtr interface{}) {
	destValue := reflect.ValueOf(destPtr).Elem()

	for i := range s.mappedFields {
		instanceValue := reflect.ValueOf(s.mappedFieldPtrs[i]).Elem().Elem()
		s.setNestedField(destValue, s.mappedFields[i].Indices, instanceValue)
	}
}
func (s *structMappings) setNestedField(root reflect.Value, pathIndices []int, value reflect.Value) {
	destField := root
	for i := range pathIndices {
		if destField.Kind() == reflect.Ptr {
			if destField.IsNil() {

				if !value.IsValid() {
					return
				}

				newValue := reflect.New(destField.Type().Elem())
				destField.Set(newValue)
			}

			destField = destField.Elem()
		}

		destField = destField.Field(pathIndices[i])
	}

	if !value.IsValid() {
		destField.Set(reflect.Zero(destField.Type()))

	} else if destField.Kind() == reflect.Ptr {
		newValue := reflect.New(destField.Type().Elem())
		newValue.Elem().Set(value)
		destField.Set(newValue)

	} else {
		destField.Set(value)
	}
}

func mapColumns(rows *sql.Rows, sType reflect.Type) (*structMappings, error) {
	columns, err := rows.Columns()
	if err != nil {
		return &structMappings{}, err
	}

	layout := getLayout(sType)

	s := &structMappings{
		mappedFieldPtrs: make([]interface{}, len(columns)),
		mappedFields:    make([]*field, len(columns)),
	}

	for i := range columns {
		field := layout.fieldsByName[columns[i]]
		if field == nil {
			return s, fmt.Errorf("no destination field for '%s'", columns[i])
		}

		s.mappedFieldPtrs[i] = reflect.New(reflect.PtrTo(field.Type)).Interface()
		s.mappedFields[i] = field
	}

	return s, nil
}
