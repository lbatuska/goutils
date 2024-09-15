package typeutils

// Created to abstract over Is_some and Is_ok
type ValueContainer interface {
	Has_value() bool
}

type Unwrappable[T any] interface {
	Expect(string) T           // panics with a provided custom message
	Unwrap() T                 // panics with a generic message
	Unwrap_or(T) T             // returns the provided default value
	Unwrap_or_default() T      // returns the default value of the type T
	Unwrap_or_else(func() T) T // returns the result of evaluating the provided function
}

// Both an Optional and Result is an Option
type Option[T any] interface {
	ValueContainer
	Unwrappable[T]
}

type OptionalI[T any] interface {
	Is_some() bool
	Is_none() bool
	Ok_or(error) Result[T]
	Ok_or_else(func() error) Result[T]
	Unwrappable[T]
}

type ResultI[T any] interface {
	Is_ok() bool
	Is_err() bool
	Ok() Optional[T]
	Err() Optional[error]
	Unwrappable[T]
}

// Ensure compile time the interfaces are implemented
var (
	_ Option[any]    = (*Optional[any])(nil)
	_ Option[any]    = (*Result[any])(nil)
	_ OptionalI[any] = (*Optional[any])(nil)
	_ ResultI[any]   = (*Result[any])(nil)
	_ ValueContainer = (*Optional[any])(nil)
	_ ValueContainer = (*Result[any])(nil)
)
