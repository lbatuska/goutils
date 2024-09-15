package typeutils

func Expect[T any](val Unwrappable[T], msg string) T {
	return val.Expect(msg)
}

func Unwrap[T any](val Unwrappable[T]) T {
	return val.Unwrap()
}

func Unwrap_or[T any](val Unwrappable[T]) (def T) {
	return val.Unwrap_or(def)
}

func Unwrap_or_default[T any](val Unwrappable[T]) T {
	return val.Unwrap_or_default()
}

func Unwrap_or_else[T any](val Unwrappable[T], f func() T) T {
	return val.Unwrap_or_else(f)
}

// Returns if the underlying data has a Value (false in case of None or Error)
func Has_value(val ValueContainer) bool {
	return val.Has_value()
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
