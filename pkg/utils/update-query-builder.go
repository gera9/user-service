package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func UpdateQueryBuilder(table string, o any) (string, int, []any) {
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("UPDATE %s SET", table))

	args := []any{}
	pos := 1

	reflectValue := reflect.ValueOf(o)
	reflectType := reflect.TypeOf(o)

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectType.Field(i)
		value := reflectValue.Field(i)

		if value.IsZero() {
			continue
		}

		if pos > 1 {
			sb.WriteString(",")
		}

		sb.WriteString(fmt.Sprintf(" %s = $%d", field.Name, pos))
		args = append(args, value.Interface())
		pos++
	}

	return sb.String(), pos, args
}
