package Type

import (
	"fmt"
	"reflect"
	"unsafe"
)

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

// meant to be used as defer Type.CatchUnwrap(Type.Ptr(&res)) or Type.CatchUnwrap(&res) if res is already a pointer
// where res is a pointer to an Option or Result returned by a function (initialized to not be nil)
// func X() (res *Optional[int]) {
// 	res = None[int]()
// 	defer CatchUnwrap(&res)
//
//  // Some possibly unsafe unwrapping of values
//  return res
// }
// === OR ===
// func X() (res Optional[int]) {
// 	res = None[int]()
//  defer Type.CatchUnwrap(Type.Ptr(&res))
//
//  // Some possibly unsafe unwrapping of values
//  return res
// }

func CatchUnwrap(ret interface{}) {
	r := recover()
	if r == nil {
		return
	}
	vp := reflect.ValueOf(r)
	if vp.Kind() != reflect.Pointer || vp.IsNil() {
		panic(r)
	}

	if _, ok := r.(OptionalerMarker); ok {
		if setOptionalNone(ret) {
			return
		}
	}

	if _, ok := r.(ResulterMarker); ok {
		if setResultError(ret, fmt.Errorf("Tried to unwrap a failed result!")) {
			return
		}
	}

	panic(r)
}

func setOptionalNone(ret interface{}) bool {
	v := reflect.ValueOf(ret)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return false
	}
	elem := v.Elem().Elem()
	if elem.Kind() != reflect.Struct {
		return false
	}
	presentField := elem.FieldByName("present")
	if !presentField.IsValid() {
		return false
	}
	if presentField.Kind() == reflect.Bool {
		presentField = reflect.NewAt(presentField.Type(),
			unsafe.Pointer(presentField.UnsafeAddr())).Elem()
		presentField.SetBool(false) // Setting it to false (None)
	} else {
		return false
	}
	return true
}

func setResultError(ret interface{}, err error) bool {
	v := reflect.ValueOf(ret)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return false
	}
	elem := v.Elem().Elem()
	if elem.Kind() != reflect.Struct {
		return false
	}
	errField := elem.FieldByName("err")
	if !errField.IsValid() {
		return false
	}
	if errField.Kind() == reflect.Interface {
		errField = reflect.NewAt(errField.Type(),
			unsafe.Pointer(errField.UnsafeAddr())).Elem()
		errField.Set(reflect.ValueOf(err)) // Setting the error
	} else {
		return false
	}
	return true
}
