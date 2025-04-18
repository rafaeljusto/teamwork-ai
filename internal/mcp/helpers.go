package mcp

import (
	"fmt"
	"slices"
)

// ParamGroup applies a series of functions to a map of parameters.
func ParamGroup(params map[string]any, funcs ...ParamFunc) error {
	for _, fn := range funcs {
		if err := fn(params); err != nil {
			return fmt.Errorf("error binding parameter: %w", err)
		}
	}
	return nil
}

// ParamFunc defines a function type that takes a map of parameters and
// returns an error. This is used to define functions that can retrieve
// parameters from a map, converting them to a specific type and applying
// middleware functions if necessary.
type ParamFunc func(map[string]any) error

// ParamMiddleware defines a function type that takes a pointer to a specific
// type and returns a boolean indicating whether to continue processing and an
// error if any issue occurs. This is used to apply middleware functions to
// parameters before they are set to the target.
type ParamMiddleware[T any] func(*T) (bool, error)

// RequiredParam retrieves a required parameter from a map, converting it to the
// specified type. It returns an error if the key is not found or if the type
// conversion fails. If the target is nil, it returns an error. It also allows
// for middleware functions to be applied to the value before setting it to the
// target. Each middleware function should return a boolean indicating whether
// to continue processing and an error if any issue occurs.
func RequiredParam[T any](target *T, key string, middlewares ...ParamMiddleware[T]) ParamFunc {
	return func(params map[string]any) error {
		return param(params, target, key, false, middlewares...)
	}
}

// OptionalParam retrieves an optional parameter from a map, converting it to
// the specified type. It returns an error if the type conversion fails. If the
// target is nil, it returns an error. It also allows for middleware functions
// to be applied to the value before setting it to the target. Each middleware
// function should return a boolean indicating whether to continue processing
// and an error if any issue occurs.
func OptionalParam[T any](target *T, key string, middlewares ...ParamMiddleware[T]) ParamFunc {
	return func(params map[string]any) error {
		return param(params, target, key, true, middlewares...)
	}
}

// OptionalPointerParam retrieves an optional parameter from a map and sets
// it to a pointer target. It converts the value to the specified type and
// applies middleware functions to the value before setting it. If the target
// is nil, it returns an error. The middleware functions should return a
// boolean indicating whether to continue processing and an error if any issue
// occurs. If the parameter is not found, it does not set the target pointer.
func OptionalPointerParam[T any](target **T, key string, middlewares ...ParamMiddleware[T]) ParamFunc {
	return func(params map[string]any) error {
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
}

func param[T any](
	params map[string]any,
	target *T,
	key string,
	optional bool,
	middlewares ...ParamMiddleware[T],
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

// RequiredNumericParam retrieves a required numeric parameter from a map,
// converting it to the target numeric type. It returns an error if the key is
// not found or if the type conversion fails. If the target is nil, it returns
// an error.
func RequiredNumericParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	target *T,
	key string,
	middlewares ...ParamMiddleware[T],
) ParamFunc {
	return func(params map[string]any) error {
		return numericParam(params, target, key, false, middlewares...)
	}
}

// OptionalNumericParam retrieves an optional numeric parameter from a map,
// converting it to the target numeric type. It returns an error if the type
// conversion fails. If the target is nil, it returns an error.
func OptionalNumericParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	target *T,
	key string,
	middlewares ...ParamMiddleware[T],
) ParamFunc {
	return func(params map[string]any) error {
		return numericParam(params, target, key, true, middlewares...)
	}
}

// OptionalNumericPointerParam retrieves an optional numeric parameter from a
// map and sets it to a pointer target. It converts the value to the specified
// numeric type and applies middleware functions to the value before setting it.
// If the target is nil, it returns an error.
func OptionalNumericPointerParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	target **T,
	key string,
	middlewares ...ParamMiddleware[T],
) ParamFunc {
	return func(params map[string]any) error {
		if target == nil {
			return fmt.Errorf("target cannot be nil")
		}
		var temp T
		var set bool
		middlewares = append(middlewares, func(*T) (bool, error) { set = true; return true, nil })
		if err := numericParam(params, &temp, key, true, middlewares...); err != nil {
			return err
		}
		if set {
			*target = &temp
		}
		return nil
	}
}

func numericParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	params map[string]any,
	target *T,
	key string,
	optional bool,
	middlewares ...ParamMiddleware[T],
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
	vType := T(v)
	for _, middleware := range middlewares {
		var err error
		if ok, err = middleware(&vType); err != nil || !ok {
			return err
		}
	}
	*target = vType
	return nil
}

// OptionalListParam retrieves an optional list parameter from a map, converting
// each item to the specified type. It returns an error if the key is not found
// or if the type conversion fails. If the target is nil, it returns an error.
func OptionalListParam[T any](target *[]T, key string) ParamFunc {
	return func(params map[string]any) error {
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
}

// OptionalNumericListParam retrieves an optional list of numeric parameters
// from a map, converting each item to the specified numeric type. It returns
// an error if the key is not found or if the type conversion fails. If the
// target is nil, it returns an error.
func OptionalNumericListParam[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	target *[]T, key string,
) ParamFunc {
	return func(params map[string]any) error {
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
}

// RestrictValues restricts the values of a parameter to a predefined set of
// allowed values. It can be used as a middleware function in the Param or
// OptionalParam functions.
func RestrictValues[T comparable](allowedValues ...T) ParamMiddleware[T] {
	return func(value *T) (bool, error) {
		if value == nil {
			return true, nil
		}
		if slices.Contains(allowedValues, *value) {
			return true, nil
		}
		return false, fmt.Errorf("value %v is not allowed, must be one of %v", *value, allowedValues)
	}
}
