package mcp

import "fmt"

// Param retrieves a required parameter from a map, converting it to the
// specified type. It returns an error if the key is not found or if the type
// conversion fails. If the target is nil, it returns an error. It also allows
// for middleware functions to be applied to the value before setting it to the
// target. Each middleware function should return a boolean indicating whether
// to continue processing and an error if any issue occurs.
func Param[T any](params map[string]any, target *T, key string, middlewares ...func(*T) (bool, error)) error {
	return param(params, target, key, false, middlewares...)
}

// OptionalParam retrieves an optional parameter from a map, converting it to
// the specified type. It returns an error if the type conversion fails. If the
// target is nil, it returns an error. It also allows for middleware functions
// to be applied to the value before setting it to the target. Each middleware
// function should return a boolean indicating whether to continue processing
// and an error if any issue occurs.
func OptionalParam[T any](params map[string]any, target *T, key string, middlewares ...func(*T) (bool, error)) error {
	return param(params, target, key, true, middlewares...)
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
	var temp T
	var set bool
	middlewares = append(middlewares, func(*T) (bool, error) { set = true; return true, nil })
	if err := param(params, &temp, key, true, middlewares...); err != nil {
		return err
	}
	if set {
		*target = &temp
	}
	return nil
}

func param[T any](
	params map[string]any,
	target *T,
	key string,
	optional bool,
	middlewares ...func(*T) (bool, error),
) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	value, ok := params[key]
	if !ok {
		if optional {
			return nil
		}
		return fmt.Errorf("parameter %s is required", key)
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

// OptionalNumericParam retrieves a required numeric parameter from a map,
// converting it to the target numeric type. It returns an error if the key is
// not found or if the type conversion fails. If the target is nil, it returns
// an error.
func NumericParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	params map[string]any, target *T, key string,
) error {
	return numericParam(params, target, key, false)
}

// OptionalNumericParam retrieves an optional numeric parameter from a map,
// converting it to the target numeric type. It returns an error if the type
// conversion fails. If the target is nil, it returns an error.
func OptionalNumericParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	params map[string]any, target *T, key string,
) error {
	return numericParam(params, target, key, true)
}

func numericParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	params map[string]any, target *T, key string, optional bool,
) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	value, ok := params[key]
	if !ok {
		if optional {
			return nil
		}
		return fmt.Errorf("parameter %s is required", key)
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
