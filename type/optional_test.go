package Type

import (
	"errors"
	"testing"

	Testing "github.com/lbatuska/goutils/testing"
)

var (
	u           = Some("something")
	v           = Some("something")
	w           = Some("something else")
	x           = Some(5)
	y           = Some(5)
	z           = Some(6)
	none        = None[int]()
	nilOptional = (*Optional[int])(nil)
)

func Test_optionalScan(t *testing.T) {
	a := None[string]()
	b := None[*string]()
	c := None[string]()
	d := None[int]()
	a.Scan("a")
	b.Scan("b")
	c.Scan(Ptr("c"))
	Testing.AssertEqual(t, "a", a.Unwrap())
	Testing.AssertEqual(t, "b", *b.Unwrap())
	Testing.AssertEqual(t, "c", c.Unwrap())
	Testing.AssertPanic(t, func() { d.Unwrap() })
}

func Test_optionalSome(t *testing.T) {
	Testing.AssertEqual(t, x, y)
	Testing.AssertEqual(t, v, u)
	Testing.AssertNotEqual(t, y, z)
	Testing.AssertNotEqual(t, u, w)
	Testing.AssertNotEqual(t, none, y)
}

func Test_optionalFunctions(t *testing.T) {
	Testing.AssertTrue(t, u.IsSome())
	Testing.AssertFalse(t, u.IsNone())
	Testing.AssertTrue(t, u.HasValue())
	Testing.AssertTrue(t, none.IsNone())
	Testing.AssertFalse(t, none.IsSome())
	Testing.AssertFalse(t, none.HasValue())
	Testing.AssertTrue(t, nilOptional.IsNone())
	Testing.AssertFalse(t, nilOptional.IsSome())
	Testing.AssertFalse(t, nilOptional.HasValue())
}

func Test_optionalExpect(t *testing.T) {
	noPanic := func() {
		_ = u.Expect("test")
	}
	panic := func() {
		_ = none.Expect("test")
	}
	panic2 := func() {
		_ = nilOptional.Expect("test")
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertPanicMessage(t, panic, "test")
	Testing.AssertPanic(t, panic2)
	Testing.AssertEqual(t, u.Unwrap(), "something")
}

func Test_optionalUnwrap(t *testing.T) {
	noPanic := func() {
		_ = u.Unwrap()
	}
	panic := func() {
		_ = none.Unwrap()
	}
	panic2 := func() {
		_ = nilOptional.Unwrap()
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertPanic(t, panic)
	Testing.AssertPanic(t, panic2)
	Testing.AssertEqual(t, u.Unwrap(), "something")
}

func Test_optionalUnwrapOr(t *testing.T) {
	noPanic := func() {
		_ = u.UnwrapOr("test")
	}
	noPanic2 := func() {
		_ = none.UnwrapOr(1)
	}
	noPanic3 := func() {
		_ = nilOptional.UnwrapOr(1)
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertNotPanic(t, noPanic2)
	Testing.AssertNotPanic(t, noPanic3)
	Testing.AssertEqual(t, u.UnwrapOr("test"), "something")
	Testing.AssertEqual(t, none.UnwrapOr(1), 1)
	Testing.AssertEqual(t, nilOptional.UnwrapOr(1), 1)
}

func Test_optionalUnwrapOrDefault(t *testing.T) {
	noPanic := func() {
		_ = u.UnwrapOrDefault()
	}
	noPanic2 := func() {
		_ = none.UnwrapOrDefault()
	}
	noPanic3 := func() {
		_ = nilOptional.UnwrapOrDefault()
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertNotPanic(t, noPanic2)
	Testing.AssertNotPanic(t, noPanic3)
	Testing.AssertEqual(t, u.UnwrapOrDefault(), "something")
	Testing.AssertEqual(t, none.UnwrapOrDefault(), 0)
	Testing.AssertEqual(t, nilOptional.UnwrapOrDefault(), 0)
}

func Test_optionalUnwrapOrElse(t *testing.T) {
	returnString := func() string {
		return "this is a string"
	}
	returnInt := func() int {
		return 10
	}

	noPanic := func() {
		_ = u.UnwrapOrElse(returnString)
	}
	noPanic2 := func() {
		_ = none.UnwrapOrElse(returnInt)
	}
	noPanic3 := func() {
		_ = nilOptional.UnwrapOrElse(returnInt)
	}

	Testing.AssertNotPanic(t, noPanic)
	Testing.AssertNotPanic(t, noPanic2)
	Testing.AssertNotPanic(t, noPanic3)
	Testing.AssertEqual(t, u.UnwrapOrElse(returnString), "something")
	Testing.AssertEqual(t, none.UnwrapOrElse(returnInt), 10)
	Testing.AssertEqual(t, nilOptional.UnwrapOrElse(returnInt), 10)
}

func Test_optionalOkOr(t *testing.T) {
	uOk := u.OkOr(errors.New("some error"))
	noneErr := none.OkOr(errors.New("some error"))
	nilOptionalErr := nilOptional.OkOr(errors.New("some error"))

	Testing.AssertTrue(t, uOk.IsOk())
	Testing.AssertTrue(t, noneErr.IsErr())
	Testing.AssertTrue(t, nilOptionalErr.IsErr())
}

func Test_optionalOkOrElse(t *testing.T) {
	returnErr := func() error {
		return errors.New("some error")
	}

	uOk := u.OkOrElse(returnErr)
	noneErr := none.OkOrElse(returnErr)
	nilOptionalErr := nilOptional.OkOrElse(returnErr)

	Testing.AssertTrue(t, uOk.IsOk())
	Testing.AssertTrue(t, noneErr.IsErr())
	Testing.AssertTrue(t, nilOptionalErr.IsErr())
}
