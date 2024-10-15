package Type

type Optional[T any] struct {
	value   T
	present bool
}

type Result[T any] struct {
	value T
	err   error
}
