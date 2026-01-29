package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

const (
	EnvVarPattern = `\$\{([A-Za-z_][A-Za-z0-9_]*)\}`
)

// readConfigFile reads configuration file content
func readConfigFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return data, nil
}

// getFileExtension returns the lowercase file extension
func getFileExtension(path string) string {
	return strings.ToLower(filepath.Ext(path))
}

// unmarshalJSON parses JSON data into config
func unmarshalJSON(data []byte, config *Config) error {
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}
	return nil
}

// errUnsupportedFormat returns error for unsupported config format
func errUnsupportedFormat(ext string) error {
	return fmt.Errorf("unsupported config file format: %s", ext)
}

// expandEnvVars expands environment variable references in the config
// This is a universal function that processes all string fields recursively
func expandEnvVars(config *Config) error {
	return expandValue(reflect.ValueOf(config).Elem())
}

// expandValue recursively processes all string fields in a struct or map
func expandValue(v reflect.Value) error {
	// Skip invalid or unexported values
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		if v.CanSet() {
			v.SetString(expandAllEnvVars(v.String()))
		}
	case reflect.Struct:
		// Skip unexported fields
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// Check if field is exported and can be set
			if field.CanInterface() && field.CanAddr() {
				if err := expandValue(field); err != nil {
					return err
				}
			}
		}
	case reflect.Map:
		// Skip nil maps
		if v.IsNil() {
			return nil
		}
		// Create a new map to hold modified values
		mapType := v.Type()
		elemType := mapType.Elem()
		newMap := reflect.MakeMapWithSize(mapType, v.Len())

		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			if !value.IsValid() {
				continue
			}

			// Create a copy of the value to expand
			newValue := reflect.New(elemType).Elem()
			newValue.Set(value)

			// Expand the copy
			if err := expandValue(newValue); err != nil {
				return err
			}

			newMap.SetMapIndex(key, newValue)
		}

		if v.CanSet() {
			v.Set(newMap)
		}
	case reflect.Slice:
		// Skip nil slices
		if v.IsNil() {
			return nil
		}
		// Process each element
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			if element.CanAddr() && element.CanSet() {
				if err := expandValue(element); err != nil {
					return err
				}
			}
		}
	case reflect.Pointer:
		// Skip nil pointers
		if v.IsNil() {
			return nil
		}
		// Dereference and process
		if v.Elem().CanAddr() {
			return expandValue(v.Elem())
		}
	}

	return nil
}

// expandAllEnvVars expands all environment variable references in a string
// Supports multiple ${VAR} patterns in a single string
func expandAllEnvVars(s string) string {
	// Use regex to find all ${VAR_NAME} patterns
	re := regexp.MustCompile(EnvVarPattern)

	// Replace function that looks up each variable in the environment
	result := re.ReplaceAllStringFunc(s, func(match string) string {
		// Extract the variable name (remove ${ and })
		varName := match[2 : len(match)-1]
		varValue := os.Getenv(varName)

		if varValue == "" {
			// If environment variable is not set, keep the template
			return match
		}

		return varValue
	})

	return result
}
