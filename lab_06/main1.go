package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Server struct {
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	Debug      bool     `json:"debug"`
	AllowedIPs []string `json:"allowed_ips"`
}

func ToJSON(v any) (string, error) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return serializeValue(val, "")
}

func serializeValue(val reflect.Value, indent string) (string, error) {
	switch val.Kind() {
	case reflect.String:
		return strconv.Quote(val.String()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil

	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil

	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return "[]", nil
		}
		var items []string
		nextIndent := indent + "\t"
		for i := 0; i < val.Len(); i++ {
			str, err := serializeValue(val.Index(i), nextIndent)
			if err != nil {
				return "", err
			}
			items = append(items, nextIndent+str)
		}
		return "[\n" + strings.Join(items, ",\n") + "\n" + indent + "]", nil

	case reflect.Struct:
		t := val.Type()
		var fields []string
		nextIndent := indent + "\t"

		for i := 0; i < val.NumField(); i++ {
			fieldType := t.Field(i)
			fieldVal := val.Field(i)

			if !fieldType.IsExported() {
				continue
			}

			jsonKey := fieldType.Tag.Get("json")
			if jsonKey == "" {
				jsonKey = fieldType.Name
			}

			valStr, err := serializeValue(fieldVal, nextIndent)
			if err != nil {
				return "", err
			}

			fields = append(fields, fmt.Sprintf("%s%q: %s", nextIndent, jsonKey, valStr))
		}

		return "{\n" + strings.Join(fields, ",\n") + "\n" + indent + "}", nil

	default:
		return "", fmt.Errorf("unsupported type: %s", val.Kind())
	}
}

func main() {
	srv := Server{
		Host:       "localhost",
		Port:       8080,
		Debug:      true,
		AllowedIPs: []string{"192.168.1.1", "10.0.0.1"},
	}

	jsonResult, err := ToJSON(srv)
	if err != nil {
		fmt.Printf("Помилка серіалізації: %v\n", err)
		return
	}
	
	fmt.Println(jsonResult)
}
