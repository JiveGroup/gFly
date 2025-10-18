package security

import (
	"github.com/gflydev/core/log"
	"github.com/gflydev/utils/str"
	"html"
	"reflect"
	"regexp"
	"strings"
)

var scriptTagPattern = regexp.MustCompile(`(?is)<script.*?>.*?</script>`)

// SanitizeStruct recursively sanitizes string fields to mitigate XSS payloads.
func SanitizeStruct(target any) {
	if target == nil {
		return
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer {
		return
	}

	sanitizeValue(val.Elem())
}

func SanitizeString(input string) string {
	if input == "" {
		return ""
	}
	clean := scriptTagPattern.ReplaceAllString(input, "")
	clean = str.Trim(clean)
	clean = html.UnescapeString(clean)
	clean = strings.ReplaceAll(clean, "\x00", "")
	return clean
}

func sanitizeValue(val reflect.Value) {
	if !val.IsValid() {
		return
	}

	switch val.Kind() {
	case reflect.Pointer:
		if !val.IsNil() {
			sanitizeValue(val.Elem())
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.CanSet() {
				sanitizeValue(field)
			} else if field.CanAddr() {
				// Non-settable fields may still be pointers/structs we can sanitize through address.
				sanitizeValue(field.Addr())
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.CanAddr() {
				sanitizeValue(elem.Addr())
			} else {
				sanitizeValue(elem)
			}
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			elem := val.MapIndex(key)
			if elem.CanInterface() {
				// Sanitize strings in map values.
				if elem.Kind() == reflect.String {
					clean := SanitizeString(elem.String())
					val.SetMapIndex(key, reflect.ValueOf(clean))
				}
			}
		}
	case reflect.String:
		clean := SanitizeString(val.String())
		val.SetString(clean)
	default:
		log.Tracef("unhandled default case for value type %v", val.Kind())
	}
}
