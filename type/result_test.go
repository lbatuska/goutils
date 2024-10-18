package Type

import (
	"errors"
	"reflect"
	"testing"

	Testing "github.com/lbatuska/goutils/testing"
)

var (
	r         = Ok("something")
	s         = Ok(5)
	errResult = Err_t(errors.New("some error"), 5)
	nilResult = (*Result[int])(nil)
)

func Test_resultConstructors(t *testing.T) {
	err := errors.New("some error")
	w := Err[int](err)
	x := Err_t(err, 5)

	Testing.AssertEqual(t, "something", r.value)
	Testing.AssertEqual(t, 5, s.value)
	Testing.AssertEqual(t, err, w.err)
	Testing.AssertEqual(t, err, x.err)
}

func Test_resultFunctions(t *testing.T) {
	Testing.AssertTrue(t, r.IsOk())
	Testing.AssertTrue(t, s.IsOk())
	Testing.AssertFalse(t, r.IsErr())
	Testing.AssertFalse(t, s.IsErr())
	Testing.AssertTrue(t, errResult.IsErr())
	Testing.AssertTrue(t, nilResult.IsErr())
	Testing.AssertFalse(t, errResult.IsOk())
	Testing.AssertFalse(t, nilResult.IsOk())
	Testing.AssertTrue(t, r.HasValue())
	Testing.AssertTrue(t, s.HasValue())
	Testing.AssertFalse(t, errResult.HasValue())
	Testing.AssertFalse(t, nilResult.HasValue())
}

func Test_resultExpect(t *testing.T) {
	noPanic := func() {
		_ = r.Expect("test")
	}
	panic := func() {
		_ = errResult.Expect("test")
	}
	panic2 := func() {
		_ = nilResult.Expect("test")
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertPanicMessage(t, panic, "test")
	Testing.AssertPanic(t, panic2)
	Testing.AssertEqual(t, r.Unwrap(), "something")
}

func Test_resultUnwrap(t *testing.T) {
	noPanic := func() {
		_ = r.Unwrap()
	}
	panic := func() {
		_ = errResult.Unwrap()
	}
	panic2 := func() {
		_ = nilResult.Unwrap()
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertPanic(t, panic)
	Testing.AssertPanic(t, panic2)
	Testing.AssertEqual(t, r.Unwrap(), "something")
}

func Test_resultUnwrapOr(t *testing.T) {
	noPanic := func() {
		_ = r.UnwrapOr("test")
	}
	noPanic2 := func() {
		_ = errResult.UnwrapOr(1)
	}
	noPanic3 := func() {
		_ = nilResult.UnwrapOr(1)
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertNotPanic(t, noPanic2)
	Testing.AssertNotPanic(t, noPanic3)
	Testing.AssertEqual(t, r.UnwrapOr("test"), "something")
	Testing.AssertEqual(t, errResult.UnwrapOr(1), 1)
	Testing.AssertEqual(t, nilResult.UnwrapOr(1), 1)
}

func Test_resultUnwrapOrDefault(t *testing.T) {
	noPanic := func() {
		_ = r.UnwrapOrDefault()
	}
	noPanic2 := func() {
		_ = errResult.UnwrapOrDefault()
	}
	noPanic3 := func() {
		_ = nilResult.UnwrapOrDefault()
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertNotPanic(t, noPanic2)
	Testing.AssertNotPanic(t, noPanic3)
	Testing.AssertEqual(t, r.UnwrapOrDefault(), "something")
	Testing.AssertEqual(t, errResult.UnwrapOrDefault(), 0)
	Testing.AssertEqual(t, nilResult.UnwrapOrDefault(), 0)
}

func Test_resultUnwrapOrElse(t *testing.T) {
	returnString := func() string {
		return "this is a string"
	}
	returnInt := func() int {
		return 10
	}

	noPanic := func() {
		_ = r.UnwrapOrElse(returnString)
	}
	noPanic2 := func() {
		_ = errResult.UnwrapOrElse(returnInt)
	}
	noPanic3 := func() {
		_ = nilResult.UnwrapOrElse(returnInt)
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertNotPanic(t, noPanic2)
	Testing.AssertNotPanic(t, noPanic3)
	Testing.AssertEqual(t, r.UnwrapOrElse(returnString), "something")
	Testing.AssertEqual(t, errResult.UnwrapOrElse(returnInt), 10)
	Testing.AssertEqual(t, nilResult.UnwrapOrElse(returnInt), 10)
}

func Test_resultExpectErr(t *testing.T) {
	panic := func() {
		_ = r.ExpectErr("test")
	}
	noPanic := func() {
		_ = errResult.ExpectErr("test")
	}
	noPanic2 := func() {
		_ = nilResult.ExpectErr("test")
	}

	Testing.AssertPanic(t, panic)
	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertNotPanic(t, noPanic2)
	Testing.AssertEqual(t, errResult.err, errResult.ExpectErr("test"))
}

func Test_resultUnwrapErr(t *testing.T) {
	panic := func() {
		_ = r.UnwrapErr()
	}
	noPanic := func() {
		_ = errResult.UnwrapErr()
	}
	noPanic2 := func() {
		_ = nilResult.UnwrapErr()
	}

	Testing.AssertPanic(t, panic)
	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertNotPanic(t, noPanic2)
	Testing.AssertEqual(t, errResult.err, errResult.UnwrapErr())
}

func Test_resultOk(t *testing.T) {
	Testing.AssertEqual(t, u, r.Ok())
	Testing.AssertEqual(t, none, errResult.Ok())
	Testing.AssertEqual(t, none, nilResult.Ok())
}

func Test_resultErr(t *testing.T) {
	rOptional := r.Err()
	errOptional := errResult.Err()
	errOptional2 := nilResult.Err()
	err := errors.New("some error")

	Testing.AssertTrue(t, rOptional.IsNone())
	Testing.AssertTrue(t, errOptional.IsSome())
	Testing.AssertTrue(t, errOptional2.IsSome())
	Testing.AssertEqual(t, reflect.TypeOf(err), reflect.TypeOf(errOptional.value))
	Testing.AssertEqual(t, reflect.TypeOf(err), reflect.TypeOf(errOptional2.value))
}
