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

// ==========================================
// 1. РЕАЛІЗАЦІЯ TO_YAML (Практична №9)
// ==========================================

func ToYAML(v any) (string, error) {
	var sb strings.Builder
	err := serializeYAMLValue(reflect.ValueOf(v), 0, &sb, false)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func serializeYAMLValue(v reflect.Value, indent int, sb *strings.Builder, isSliceElement bool) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	indentStr := strings.Repeat("  ", indent)

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			tag := field.Tag.Get("json")
			if tag == "" || tag == "-" {
				tag = strings.ToLower(field.Name)
			} else {
				tag = strings.Split(tag, ",")[0]
			}

			if i == 0 && isSliceElement {
				sb.WriteString(fmt.Sprintf("- %s: ", tag))
			} else {
				sb.WriteString(fmt.Sprintf("%s%s: ", indentStr, tag))
			}

			if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Struct {
				sb.WriteString("\n")
				if err := serializeYAMLValue(fieldValue, indent+1, sb, false); err != nil {
					return err
				}
			} else {
				if err := serializeYAMLValue(fieldValue, 0, sb, false); err != nil {
					return err
				}
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if elem.Kind() == reflect.Struct {
				if err := serializeYAMLValue(elem, indent, sb, true); err != nil {
					return err
				}
			} else {
				sb.WriteString(fmt.Sprintf("%s- ", indentStr))
				if err := serializeYAMLValue(elem, 0, sb, false); err != nil {
					return err
				}
			}
		}

	case reflect.String:
		sb.WriteString(fmt.Sprintf("\"%s\"\n", v.String()))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sb.WriteString(fmt.Sprintf("%d\n", v.Int()))

	case reflect.Bool:
		sb.WriteString(fmt.Sprintf("%t\n", v.Bool()))

	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}
	return nil
}

// ==========================================
// 2. РЕАЛІЗАЦІЯ TO_JSON (Практична №6)
// ==========================================

func ToJSON(v any) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return serializeJSONValue(val, "")
}

func serializeJSONValue(val reflect.Value, indent string) (string, error) {
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
			str, err := serializeJSONValue(val.Index(i), nextIndent)
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

			valStr, err := serializeJSONValue(fieldVal, nextIndent)
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

	fmt.Println("--- ТЕСТ TO_YAML ---")
	yamlResult, _ := ToYAML(srv)
	fmt.Print(yamlResult)

	fmt.Println("\n--- ТЕСТ TO_JSON ---")
	jsonResult, _ := ToJSON(srv)
	fmt.Println(jsonResult)
}
