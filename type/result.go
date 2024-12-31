package Type

import (
	"database/sql"
	"errors"
	"fmt"

	Assert "github.com/lbatuska/goutils/assert"
)

// CTORS BEGIN
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func Err_t[T any](err error, x T) Result[T] {
	return Result[T]{err: err}
}

// CTORS END

func (res *Result[T]) IsOk() bool {
	if res == nil {
		return false
	}
	return res.err == nil
}

func (res *Result[T]) IsErr() bool {
	if res == nil {
		return true
	}
	return res.err != nil
}

func (res *Result[T]) HasValue() bool {
	if res == nil {
		return false
	}
	return res.IsOk()
}

// UNWRAPPABLE INTERFACE
func (res *Result[T]) Expect(msg string) T {
	Assert.NotNil(res)
	if res.err == nil {
		return res.value
	}
	panic(msg)
}

func (res *Result[T]) Unwrap() T {
	Assert.NotNil(res)
	if res.err == nil {
		return res.value
	}
	panic("Tried unwrapping a Result that had an error value!")
}

func (res *Result[T]) UnwrapOr(val T) T {
	if res != nil {
		if res.err == nil {
			return res.value
		}
	}
	return val
}

func (res *Result[T]) UnwrapOrDefault() T {
	if res != nil {
		if res.err == nil {
			return res.value
		}
	}
	var ret T
	return ret
}

func (res *Result[T]) UnwrapOrElse(f func() T) T {
	if res != nil {
		if res.err == nil {
			return res.value
		}
	}
	return f()
}

// UNWRAPPABLE INTERFACE

// This function panic on Ok instead of Err
func (res *Result[T]) ExpectErr(msg string) error {
	if res == nil {
		return errors.New("ExpectErr was called on a nil Result.")
	}

	if res.err != nil {
		return res.err
	}
	panic(msg)
}

// This function panic on Ok instead of Err
func (res *Result[T]) UnwrapErr() error {
	if res == nil {
		return errors.New("UnwrapErr was called on a nil Result.")
	}
	if res.err != nil {
		return res.err
	}
	panic("Expect_err was called with an Ok value")
}

// transforms Result into Option, mapping Ok(v) to Some(v) and Err(e) to None
func (res *Result[T]) Ok() Optional[T] {
	if res != nil {
		if res.err == nil {
			return Optional[T]{value: res.value, present: true}
		}
	}
	return Optional[T]{present: false}
}

// transforms Result into Option, mapping Err(e) to Some(e) and Ok(v) to None
func (res *Result[T]) Err() Optional[error] {
	if res == nil {
		return Optional[error]{
			value: errors.New("Err was called on a nil Result."), present: true,
		}
	}
	if res.err != nil {
		return Optional[error]{value: res.err, present: true}
	}
	return Optional[error]{present: false}
}

func (res *Result[T]) Scan(src interface{}) error {
	// DB had a null value
	if src == nil {
		res.err = errors.New("src was nil!")
		return nil
	}
	// If T is a scanner
	if scanner, ok := any(&res.value).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			res.err = err
			return err
		}
		res.err = nil
		return nil
	}
	// We implement parsing for some builtin types
	mismatchErr := fmt.Errorf("Type of src (%T) doesn't match type of Result[%T]!", src, res.value)

	switch v := any(&res.value).(type) {

	case *string:
		if str, ok := src.(string); ok {
			*v = str
			res.err = nil
			return nil
		}
		if b, ok := src.([]byte); ok {
			*v = string(b)
			res.err = nil
			return nil
		}
		res.err = mismatchErr
		return res.err

	case *int:
		if s, ok := src.(int); ok {
			*v = s
			res.err = nil
			return nil
		}
		res.err = mismatchErr
		return res.err
	case *int32:
		if s, ok := src.(int32); ok {
			*v = s
			res.err = nil
			return nil
		}
		res.err = mismatchErr
		return res.err
	case *int64:
		if s, ok := src.(int64); ok {
			*v = s
			res.err = nil
			return nil
		}
		res.err = mismatchErr
		return res.err

	case *bool:
		switch s := src.(type) {
		case bool:
			*v = s
			res.err = nil
			return nil
		}
		res.err = mismatchErr
		return res.err

	case *float64:
		switch s := src.(type) {
		case float64:
			*v = s
			res.err = nil
			return nil
		case float32:
			*v = float64(s)
			res.err = nil
			return nil

		}
		res.err = mismatchErr
		return res.err

	case *float32:
		switch s := src.(type) {
		case float32:
			*v = s
			res.err = nil
			return nil
		}
		res.err = mismatchErr
		return res.err
	}
	// We couldnt parse the value
	err := fmt.Errorf("Unsupported type %T, and the type doesn't implement sql.Scanner!", src)
	res.err = err
	return err
}
