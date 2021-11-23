package ncore

import (
	"reflect"
)

type InitializeChecker interface {
	HasInitialized() bool
}

/// InitStruct reflect fields in struct and run initializer function
func InitStruct(s interface{}, initFn func(name string, i interface{}) error) error {
	// Reflect on struct element
	rv := reflect.ValueOf(s).Elem()

	// Iterate fields
	for i := 0; i < rv.NumField(); i++ {
		// Get field interface
		fieldValue := rv.Field(i).Interface()
		fieldName := rv.Type().Field(i).Name

		err := initFn(fieldName, fieldValue)
		if err != nil {
			return err
		}
	}

	return nil
}
