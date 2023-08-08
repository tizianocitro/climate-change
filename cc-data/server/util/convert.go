package util

import (
	"encoding/json"
	"reflect"
)

func Convert(a, b interface{}) error {
	js, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, b)
}

func ConvertStructToKeys(s interface{}) []string {
	var keys []string

	// Get the type of the struct
	structType := reflect.TypeOf(s)

	// Loop through the fields of the struct
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Append the field name to the keys slice
		keys = append(keys, field.Name)
	}

	return keys
}
