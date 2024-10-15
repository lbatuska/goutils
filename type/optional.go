package Type

import "github.com/lbatuska/goutils/assert"

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
	assert.NotNil(opt)
	if opt.present {
		return opt.value
	}
	panic(msg)
}

func (opt *Optional[T]) Unwrap() T {
	assert.NotNil(opt)
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
