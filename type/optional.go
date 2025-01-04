package Type

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

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
	Assert.NotNil(opt)
	opt.present = false
	// DB had a null value
	if src == nil {
		return nil
	}

	// If T is a scanner
	if scanner, ok := any(opt.value).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}
		opt.present = true
		return nil
	}

	if scanres := opt.scanBuiltin(src); scanres.IsSome() {
		return scanres.Unwrap()
	}

	return fmt.Errorf("Unsupported type %T or differs from Optional[%T], and the type doesn't implement sql.Scanner!", src, opt.value)
}

func (opt *Optional[T]) scanBuiltin(src interface{}) Optional[error] {
	opt.present = false
	// First handle the special cases where we allow conversion between types
	// This is usually just parsing []byte into type
	if scanres := opt.scanTimeSpecial(src); scanres.IsSome() {
		return scanres
	}
	if scanres := opt.scanStringSpecial(src); scanres.IsSome() {
		return scanres
	}

	srcVal := reflect.ValueOf(src)
	optType := reflect.TypeOf(opt.value)

	optElemType := optType
	if optElemType.Kind() == reflect.Pointer {
		optElemType = optElemType.Elem()
	}

	srcElemType := srcVal.Type()
	if srcElemType.Kind() == reflect.Pointer {
		srcElemType = srcElemType.Elem()
	}

	if srcElemType != optElemType {
		return Some(fmt.Errorf("Optional[%T] (aka %T) differs from %T!", opt.value, opt.value, src))
	}

	if optType.Kind() == reflect.Pointer {
		if srcVal.Kind() == reflect.Pointer {
			opt.value = srcVal.Interface().(T)
		} else {
			newPtr := reflect.New(optElemType)
			newPtr.Elem().Set(srcVal)
			opt.value = newPtr.Interface().(T)
		}
	} else {
		if srcVal.Kind() == reflect.Pointer {
			opt.value = srcVal.Elem().Interface().(T)
		} else {
			opt.value = srcVal.Interface().(T)
		}
	}
	opt.present = true
	return Some[error](nil)
}

func (opt *Optional[T]) scanStringSpecial(src interface{}) Optional[error] {
	opt.present = false
	switch v := any(opt.value).(type) {
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
			reflect.ValueOf(&opt.value).Elem().Set(reflect.ValueOf(string(s)))
			goto ok
		case *[]byte:
			reflect.ValueOf(&opt.value).Elem().Set(reflect.ValueOf(string(*s)))
			goto ok
		}
	}
	return None[error]()
ok:
	opt.present = true
	return Some[error](nil)
}

func (opt *Optional[T]) scanTimeSpecial(src interface{}) Optional[error] {
	opt.present = false
	switch v := any(opt.value).(type) {
	case *time.Time:
		switch t := src.(type) {
		case []byte:
			parsedTime, err := time.Parse(time.RFC3339, string(t))
			if err == nil {
				*v = parsedTime
				goto ok
			} else {
				return Some(err)
			}
		case *[]byte:
			parsedTime, err := time.Parse(time.RFC3339, string(*t))
			if err == nil {
				*v = parsedTime
				goto ok
			} else {
				return Some(err)
			}
		}
	case time.Time:
		switch t := src.(type) {
		case []byte:
			parsedTime, err := time.Parse(time.RFC3339, string(t))
			if err == nil {
				reflect.ValueOf(&opt.value).Elem().Set(reflect.ValueOf(parsedTime))
				goto ok
			} else {
				return Some(err)
			}
		case *[]byte:
			parsedTime, err := time.Parse(time.RFC3339, string(*t))
			if err == nil {
				reflect.ValueOf(&opt.value).Elem().Set(reflect.ValueOf(parsedTime))
				goto ok
			} else {
				return Some(err)
			}
		}

	}
	return None[error]()
ok:
	opt.present = true
	return Some[error](nil)
}

func (opt Optional[T]) MarshalJSON() ([]byte, error) {
	if !opt.present {
		// Return null for `omitempty` compatibility
		return []byte("null"), nil
		// panic("Tried to marshal an Optional that did not have a value!")
	}
	return json.Marshal(opt.value)
}

func (opt *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		opt.present = false
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	opt.value = value
	opt.present = true
	return nil
}

func (o Optional[T]) Value() (driver.Value, error) {
	if !o.present {
		return nil, nil
	}

	switch v := any(o.Value).(type) {

	case driver.Valuer:
		return v.Value()

	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64,
		complex64, complex128:
		return v, nil

	case *string, *bool,
		*int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64, *uintptr,
		*float32, *float64,
		*complex64, *complex128:
		return v, nil

	case fmt.Stringer:
		return v.String(), nil

	default:
		return nil, errors.New("unsupported type for Optional")
	}
}
