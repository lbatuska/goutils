package typeutils

import (
	"errors"

	"github.com/lbatuska/goutils/assert"
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
	assert.NotNil(res)
	if res.err == nil {
		return res.value
	}
	panic(msg)
}

func (res *Result[T]) Unwrap() T {
	assert.NotNil(res)
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
