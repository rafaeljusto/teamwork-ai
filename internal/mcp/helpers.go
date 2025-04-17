package mcp

import "fmt"

// OptionalParam retrieves an optional parameter from a map, converting it to
// the specified type. It returns an error if the key is not found or if the
// type conversion fails. If the target is nil, it returns an error. It also
// allows for middleware functions to be applied to the value before setting it
// to the target. Each middleware function should return a boolean indicating
// whether to continue processing and an error if any issue occurs.
func OptionalParam[T any](params map[string]any, target *T, key string, middlewares ...func(*T) (bool, error)) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	value, ok := params[key]
	if !ok {
		return nil
	}
	v, ok := value.(T)
	if !ok {
		return fmt.Errorf("invalid type for %s: expected %T, got %T", key, *target, value)
	}
	for _, middleware := range middlewares {
		var err error
		if ok, err = middleware(&v); err != nil || !ok {
			return err
		}
	}
	*target = v
	return nil
}

// OptionalPointerParam retrieves an optional parameter from a map and sets
// it to a pointer target. It converts the value to the specified type and
// applies middleware functions to the value before setting it. If the target
// is nil, it returns an error. The middleware functions should return a
// boolean indicating whether to continue processing and an error if any issue
// occurs. If the parameter is not found, it does not set the target pointer.
func OptionalPointerParam[T any](
	params map[string]any,
	target **T,
	key string,
	middlewares ...func(*T) (bool, error),
) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	var set bool
	middlewares = append(middlewares, func(t *T) (bool, error) {
		set = true
		return true, nil
	})
	var temp T
	if err := OptionalParam(params, &temp, key, middlewares...); err != nil {
		return err
	}
	if set {
		*target = &temp
	}
	return nil
}

// OptionalNumericParam retrieves an optional numeric parameter from a map,
// converting it to the target numeric type. It returns an error if the key is
// not found or if the type conversion fails. If the target is nil, it returns
// an error.
func OptionalNumericParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	params map[string]any, target *T, key string,
) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	value, ok := params[key]
	if !ok {
		return nil
	}
	v, ok := value.(float64)
	if !ok {
		return fmt.Errorf("invalid type for %s: expected %T, got %T", key, *target, value)
	}
	*target = T(v)
	return nil
}

// OptionalListParam retrieves an optional list parameter from a map, converting
// each item to the specified type. It returns an error if the key is not found
// or if the type conversion fails. If the target is nil, it returns an error.
func OptionalListParam[T any](params map[string]any, target *[]T, key string) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	value, ok := params[key]
	if !ok {
		return nil
	}
	array, ok := value.([]any)
	if !ok {
		return fmt.Errorf("invalid type for %s: expected []any, got %T", key, value)
	}
	*target = make([]T, 0, len(array))
	for _, item := range array {
		v, ok := item.(T)
		if !ok {
			var zero T
			return fmt.Errorf("invalid type in %s: expected %T, got %T", key, zero, item)
		}
		*target = append(*target, v)
	}
	return nil
}

// OptionalNumericListParam retrieves an optional list of numeric parameters
// from a map, converting each item to the specified numeric type. It returns
// an error if the key is not found or if the type conversion fails. If the
// target is nil, it returns an error.
func OptionalNumericListParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	params map[string]any, target *[]T, key string,
) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	value, ok := params[key]
	if !ok {
		return nil
	}
	array, ok := value.([]any)
	if !ok {
		return fmt.Errorf("invalid type for %s: expected []any, got %T", key, value)
	}
	*target = make([]T, 0, len(array))
	for _, item := range array {
		v, ok := item.(float64)
		if !ok {
			return fmt.Errorf("invalid type in %s: expected float64, got %T", key, item)
		}
		*target = append(*target, T(v))
	}
	return nil
}
