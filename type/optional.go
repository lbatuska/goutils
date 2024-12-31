package Type

import (
	"database/sql"
	"fmt"

	Assert "github.com/lbatuska/goutils/assert"
)

// CTORS BEGIN
func Some[T any](value T) Optional[T] {
	return Optional[T]{value, true}
}

func None[T any]() Optional[T] {
	return Optional[T]{present: false}
}

// Because None has no type passing a value of desired type might be preferred syntax over providing type on the function call
func None_t[T any](T) Optional[T] {
	return Optional[T]{present: false}
}

// CTORS END

func (opt *Optional[T]) IsSome() bool {
	if opt == nil {
		return false
	}
	return opt.present
}

func (opt *Optional[T]) IsNone() bool {
	if opt == nil {
		return true
	}
	return !opt.present
}

func (opt *Optional[T]) HasValue() bool {
	if opt == nil {
		return false
	}
	return opt.IsSome()
}

// UNWRAPPABLE INTERFACE BEGIN
func (opt *Optional[T]) Expect(msg string) T {
	Assert.NotNil(opt)
	if opt.present {
		return opt.value
	}
	panic(msg)
}

func (opt *Optional[T]) Unwrap() T {
	Assert.NotNil(opt)
	if opt.present {
		return opt.value
	}
	panic("Tried unwrapping an Optional that did not have a value!")
}

func (opt *Optional[T]) UnwrapOr(val T) T {
	if opt != nil {
		if opt.present {
			return opt.value
		}
	}
	return val
}

func (opt *Optional[T]) UnwrapOrDefault() T {
	if opt != nil {
		if opt.present {
			return opt.value
		}
	}
	var res T
	return res
}

func (opt *Optional[T]) UnwrapOrElse(f func() T) T {
	if opt != nil {
		if opt.present {
			return opt.value
		}
	}
	return f()
}

// UNWRAPPABLE INTERFACE END

// transforms Some(v) to Ok(v), and None to Err(err)
func (opt *Optional[T]) OkOr(err error) Result[T] {
	if opt != nil {
		if opt.present {
			return Ok(opt.value)
		}
	}

	return Err[T](err)
}

// transforms Some(v) to Ok(v), and None to a value of Err using the provided function
func (opt *Optional[T]) OkOrElse(f func() error) Result[T] {
	if opt != nil {
		if opt.present {
			return Ok(opt.value)
		}
	}
	return Err[T](f())
}

func (opt *Optional[T]) Scan(src interface{}) error {
	// DB had a null value
	if src == nil {
		opt.present = false
		return nil
	}
	// If T is a scanner
	if scanner, ok := any(&opt.value).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			opt.present = false
			return err
		}
		opt.present = true
		return nil
	}
	// We implement parsing for some builtin types
	opt.present = false
	switch v := any(&opt.value).(type) {

	case *string:
		if str, ok := src.(string); ok {
			*v = str
			opt.present = true
			return nil
		}
		if b, ok := src.([]byte); ok {
			*v = string(b)
			opt.present = true
			return nil
		}

	case *int:
		if s, ok := src.(int); ok {
			*v = s
			opt.present = true
			return nil
		}
	case *int32:
		if s, ok := src.(int32); ok {
			*v = s
			opt.present = true
			return nil
		}
	case *int64:
		if s, ok := src.(int64); ok {
			*v = s
			opt.present = true
			return nil
		}

	case *bool:
		switch s := src.(type) {
		case bool:
			*v = s
			opt.present = true
			return nil
			// We could technically allow this, however we try to avoid implicit conversions to ensure type safety.
			// case string:
			// 	if s == "1" || s == "true" || s == "t" {
			// 		*v = true
			// 		opt.present = true
			// 		return nil
			// 	}
			// 	if s == "0" || s == "false" || s == "f" {
			// 		*v = false
			// 		opt.present = true
			// 		return nil
			// 	}
		}

	case *float64:
		switch s := src.(type) {
		case float64:
			*v = s
			opt.present = true
			return nil
		case float32:
			*v = float64(s)
			opt.present = true
			return nil
		}

	case *float32:
		switch s := src.(type) {
		case float32:
			*v = s
			opt.present = true
			return nil
		}
	}
	// We couldnt parse the value
	opt.present = false
	return fmt.Errorf("unsupported type %T or differs from Optional[%T], and the type doesn't implement sql.Scanner", src, opt.value)
}
