package Type

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

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
	Assert.NotNil(res)
	e := fmt.Errorf("Unsupported type %T or differs from Result[%T], and the type doesn't implement sql.Scanner!",
		src, res.value)
	res.err = e
	// DB had a null value
	if src == nil {
		return nil
	}

	// If T is a scanner
	if scanner, ok := any(res.value).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}
		res.err = nil
		return nil
	}

	if scanres := res.scanBuiltin(src); scanres.IsSome() {
		return scanres.Unwrap()
	}

	return e
}

func (res *Result[T]) scanBuiltin(src interface{}) Optional[error] {
	res.err = nil
	// First handle the special cases where we allow conversion between types
	// This is usually just parsing []byte into type
	if scanres := res.scanTimeSpecial(src); scanres.IsSome() {
		return scanres
	}
	if scanres := res.scanStringSpecial(src); scanres.IsSome() {
		return scanres
	}

	srcVal := reflect.ValueOf(src)
	optType := reflect.TypeOf(res.value)

	optElemType := optType
	if optElemType.Kind() == reflect.Pointer {
		optElemType = optElemType.Elem()
	}

	srcElemType := srcVal.Type()
	if srcElemType.Kind() == reflect.Pointer {
		srcElemType = srcElemType.Elem()
	}

	if srcElemType != optElemType {
		e := fmt.Errorf("Result[%T] (aka %T) differs from %T!", res.value, res.value, src)
		res.err = e
		return Some(e)
	}

	if optType.Kind() == reflect.Pointer {
		if srcVal.Kind() == reflect.Pointer {
			res.value = srcVal.Interface().(T)
		} else {
			newPtr := reflect.New(optElemType)
			newPtr.Elem().Set(srcVal)
			res.value = newPtr.Interface().(T)
		}
	} else {
		if srcVal.Kind() == reflect.Pointer {
			res.value = srcVal.Elem().Interface().(T)
		} else {
			res.value = srcVal.Interface().(T)
		}
	}
	res.err = nil
	return Some[error](nil)
}

func (res *Result[T]) scanStringSpecial(src interface{}) Optional[error] {
	switch v := any(res.value).(type) {
	case *string:
		switch s := src.(type) {
		case []byte:
			*v = string(s)
			goto ok
		case *[]byte:
			*v = string(*s)
			goto ok
		}
	case string:
		switch s := src.(type) {
		case []byte:
			reflect.ValueOf(&res.value).Elem().Set(reflect.ValueOf(string(s)))
			goto ok
		case *[]byte:
			reflect.ValueOf(&res.value).Elem().Set(reflect.ValueOf(string(*s)))
			goto ok
		}
	}
	return None[error]()
ok:
	res.err = nil
	return Some[error](nil)
}

func (res *Result[T]) scanTimeSpecial(src interface{}) Optional[error] {
	switch v := any(res.value).(type) {
	case *time.Time:
		switch t := src.(type) {
		case []byte:
			parsedTime, err := time.Parse(time.RFC3339, string(t))
			if err == nil {
				*v = parsedTime
				goto ok
			} else {
				res.err = err
				return Some(err)
			}
		case *[]byte:
			parsedTime, err := time.Parse(time.RFC3339, string(*t))
			if err == nil {
				*v = parsedTime
				goto ok
			} else {
				res.err = err
				return Some(err)
			}
		}
	case time.Time:
		switch t := src.(type) {
		case []byte:
			parsedTime, err := time.Parse(time.RFC3339, string(t))
			if err == nil {
				reflect.ValueOf(&res.value).Elem().Set(reflect.ValueOf(parsedTime))
				goto ok
			} else {
				res.err = err
				return Some(err)
			}
		case *[]byte:
			parsedTime, err := time.Parse(time.RFC3339, string(*t))
			if err == nil {
				reflect.ValueOf(&res.value).Elem().Set(reflect.ValueOf(parsedTime))
				goto ok
			} else {
				res.err = err
				return Some(err)
			}
		}

	}
	return None[error]()
ok:
	res.err = nil
	return Some[error](nil)
}
