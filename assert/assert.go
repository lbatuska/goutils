package assert

import (
	"fmt"
	"path/filepath"
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

func NotNil[T any](v *T) {
	if v == nil {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sValue should not be nil!%s %s/%s():%d => [%T](%+v)%s",
			colorRed, colorGreen, f, fn, l, v, v, colorNone))
	}
}

func Nil[T any](v *T) {
	if v != nil {
		f, fn, l := file_func_line()
		panic(fmt.Sprintf("%sValue should be nil!%s %s/%s():%d => [%T](%+v)%s",
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
