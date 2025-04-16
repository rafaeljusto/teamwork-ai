package mcp

import "fmt"

// OptionalParam retrieves an optional parameter from a map, returning the value
// and a boolean indicating if it was found.
func OptionalParam[T any](params map[string]any, key string) (T, bool, error) {
	value, ok := params[key]
	if !ok {
		var zero T
		return zero, false, nil
	}
	v, ok := value.(T)
	if !ok {
		var zero T
		return zero, false, fmt.Errorf("invalid type for %s: expected %T, got %T", key, zero, value)
	}
	return v, true, nil
}

// OptionalNumericParam retrieves an optional numeric parameter from a map,
// converting it to the specified numeric type. It returns the value, a boolean
// indicating if it was found, and an error if the conversion fails.
func OptionalNumericParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	params map[string]any, key string,
) (T, bool, error) {
	value, ok := params[key]
	if !ok {
		var zero T
		return zero, false, nil
	}
	v, ok := value.(float64)
	if !ok {
		var zero T
		return zero, false, fmt.Errorf("invalid type for %s: expected %T, got %T", key, zero, value)
	}
	return T(v), true, nil
}

// OptionalListParam retrieves an optional list parameter from a map, converting
// each item to the specified type. It returns the list, a boolean indicating if
// it was found, and an error if the conversion fails.
func OptionalListParam[T any](params map[string]any, key string) ([]T, bool, error) {
	value, ok := params[key]
	if !ok {
		return nil, false, nil
	}
	array, ok := value.([]any)
	if !ok {
		return nil, false, fmt.Errorf("invalid type for %s: expected []any, got %T", key, value)
	}
	var result []T
	for _, item := range array {
		v, ok := item.(T)
		if !ok {
			var zero T
			return nil, false, fmt.Errorf("invalid type in %s: expected %T, got %T", key, zero, item)
		}
		result = append(result, v)
	}
	return result, true, nil
}

// OptionalNumericListParam retrieves an optional list of numeric parameters
// from a map, converting each item to the specified numeric type. It returns
// the list, a boolean indicating if it was found, and an error if the
// conversion fails.
func OptionalNumericListParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	params map[string]any, key string,
) ([]T, bool, error) {
	value, ok := params[key]
	if !ok {
		return nil, false, nil
	}
	array, ok := value.([]any)
	if !ok {
		return nil, false, fmt.Errorf("invalid type for %s: expected []any, got %T", key, value)
	}
	var result []T
	for _, item := range array {
		v, ok := item.(float64)
		if !ok {
			return nil, false, fmt.Errorf("invalid type in %s: expected float64, got %T", key, item)
		}
		result = append(result, T(v))
	}
	return result, true, nil
}
