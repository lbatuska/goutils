package Assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

const (
	colorRed   = "\033[0;31m"
	colorNone  = "\033[0m"
	colorGreen = "\033[32m"
)

func file_func_line() (string, string, int) {
	pc, f, l, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc).Name()
	f = filepath.Base(f)
	fn = fn[strings.LastIndex(fn, "/")+1:]
	return f, fn, l
}

// Returns false for reflect.Invalid, you should handle that before calling this function
func IsNillable(kind reflect.Kind) bool {
	switch kind {
	// based on reflect/type.go -> Kind
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

func NotNil[T any](v T) {
	val := reflect.ValueOf(v)
	kind := val.Kind()
	if reflect.Invalid == kind {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sValue should not be nil!(It is reflect.Invalid)%s %s/%s():%d => [%T](%+v)%s",
			colorRed, colorGreen, f, fn, l, v, v, colorNone))
	}
	if IsNillable(kind) {
		if val.IsNil() {
			f, fn, l := file_func_line()
			panic(fmt.Sprintf("%sValue should not be nil!%s %s/%s():%d => [%T](%+v)%s",
				colorRed, colorGreen, f, fn, l, v, v, colorNone))
		}
	}
	// If it's neither Invalid nor Nil-able it cannot be nil
}

func Nil[T any](v T) {
	val := reflect.ValueOf(v)
	kind := val.Kind()
	if reflect.Invalid == kind {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sValue should be nil!(It is reflect.Invalid)%s %s/%s():%d => [%T](%+v)%s",
			colorRed, colorGreen, f, fn, l, v, v, colorNone))
	}

	if IsNillable(kind) {
		if !val.IsNil() {
			f, fn, l := file_func_line()
			panic(fmt.Sprintf("%sValue should be nil!%s %s/%s():%d => [%T](%+v)%s",
				colorRed, colorGreen, f, fn, l, v, v, colorNone))
		}
	}
	// If it's neither Invalid nor Nil-able it cannot be nil
	f, fn, l := file_func_line()
	panic(fmt.Sprintf("%sValue should be nil!%s %s/%s():%d => [%T](%+v)%s",
		colorRed, colorGreen, f, fn, l, v, v, colorNone))
}

func NotNilPtr[T any](v *T) {
	if v == nil {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sPointer value should not be nil!%s %s/%s():%d => [%T](%+v)%s",
			colorRed, colorGreen, f, fn, l, v, v, colorNone))
	}
}

func NilPtr[T any](v *T) {
	if v != nil {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sPointer value should be nil!%s %s/%s():%d => [%T](%+v)%s",
			colorRed, colorGreen, f, fn, l, v, v, colorNone))
	}
}

func Assert(v bool) {
	if !v {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sAssert (true) failed%s %s/%s():%d%s",
			colorRed, colorGreen, f, fn, l, colorNone))
	}
}

func True(v bool) {
	if !v {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sAssert (true) failed%s %s/%s():%d%s",
			colorRed, colorGreen, f, fn, l, colorNone))
	}
}

func AssertNot(v bool) {
	if v {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sAssert (false) failed%s %s/%s():%d%s",
			colorRed, colorGreen, f, fn, l, colorNone))
	}
}

func False(v bool) {
	if v {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sAssert (false) failed%s %s/%s():%d%s",
			colorRed, colorGreen, f, fn, l, colorNone))
	}
}

func Equal[T comparable](v, v2 T) {
	if v != v2 {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sAssert equals failed%s %s/%s():%d [%T](%+v) != [%T](%+v)%s",
			colorRed, colorGreen, f, fn, l, v, v, v2, v2, colorNone))
	}
}

func NotEqual[T comparable](v, v2 T) {
	if v == v2 {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sAssert not equals failed%s %s/%s():%d [%T](%+v) == [%T](%+v)%s",
			colorRed, colorGreen, f, fn, l, v, v, v2, v2, colorNone))
	}
}
