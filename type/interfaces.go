package Type

import "database/sql"

// Created to abstract over Is_some and Is_ok
type ValueContainer interface {
	HasValue() bool
}

type Unwrappable[T any] interface {
	Expect(string) T         // panics with a provided custom message
	Unwrap() T               // panics with a generic message
	UnwrapOr(T) T            // returns the provided default value
	UnwrapOrDefault() T      // returns the default value of the type T
	UnwrapOrElse(func() T) T // returns the result of evaluating the provided function
}

// Both an Optional and Result is an Optioner
type Optioner[T any] interface {
	ValueContainer
	Unwrappable[T]
}

type Optionaler[T any] interface {
	IsSome() bool
	IsNone() bool
	OkOr(error) Result[T]
	OkOrElse(func() error) Result[T]
	Unwrappable[T]
}

type Resulter[T any] interface {
	IsOk() bool
	IsErr() bool
	Ok() Optional[T]
	Err() Optional[error]
	Unwrappable[T]
}

// Ensure compile time the interfaces are implemented
var (
	_ Optioner[any]   = (*Optional[any])(nil)
	_ Optioner[any]   = (*Result[any])(nil)
	_ Optionaler[any] = (*Optional[any])(nil)
	_ Resulter[any]   = (*Result[any])(nil)
	_ ValueContainer  = (*Optional[any])(nil)
	_ ValueContainer  = (*Result[any])(nil)
	_ sql.Scanner     = (*Optional[any])(nil)
	// _ sql.Scanner     = (*Result[any])(nil)
)
