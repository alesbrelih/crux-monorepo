package pkg

import (
	"fmt"
	"reflect"
)

// traverses only first depth and only on structs for now
func GetStructStringValues(i interface{}) []string {
	values := []string{}

	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Struct {
		return values
	}

	numOfFields := value.NumField()
	for i := 0; i < numOfFields; i++ {
		values = append(values, fmt.Sprintf("%s", value.Field(i)))
	}

	return values
}
