package typeutils

func Expect[T any](val Unwrappable[T], msg string) T {
	return val.Expect(msg)
}

func Unwrap[T any](val Unwrappable[T]) T {
	return val.Unwrap()
}

func UnwrapOr[T any](val Unwrappable[T]) (def T) {
	return val.UnwrapOr(def)
}

func UnwrapOrDefault[T any](val Unwrappable[T]) T {
	return val.UnwrapOrDefault()
}

func UnwrapOrElse[T any](val Unwrappable[T], f func() T) T {
	return val.UnwrapOrElse(f)
}

// Returns if the underlying data has a Value (false in case of None or Error)
func HasValue(val ValueContainer) bool {
	return val.HasValue()
}

func ResultWrap[T any](val T, err error) Result[T] {
	if err == nil {
		return Ok(val)
	}
	return Err[T](err)
}

func ResultWrapb[T any](err error, val T) Result[T] {
	if err == nil {
		return Ok(val)
	}
	return Err[T](err)
}

func Ptr[T any](v T) *T {
	return &v
}
