package main

import (
	"errors"
	"testing"
	"time"

	. "github.com/lbatuska/goutils/logger"
	. "github.com/lbatuska/goutils/testing"
	. "github.com/lbatuska/goutils/typeutils"
)

func TestTypes(t *testing.T) {
	Create(&FileLoggerImpl{})
	go LoggerInstance().StartLogger()
	LoggerInstance().Write("TEST MSG")
	time.Sleep(time.Second * 2)
	x := Some("Test")
	AssertTrue(t, x.IsSome())
	AssertTrue(t, !x.IsNone())
	AssertEqual(t, x.Unwrap(), "Test")
	AssertTrue(t, None_t(1).IsNone())
	err := errors.New("Test")
	y := Err_t(err, "t")
	AssertEqual(t, y.UnwrapErr(), err)
	AssertError(t, y.UnwrapErr())
	z := x.OkOr(err)
	AssertTrue(t, z.IsOk())
	AssertEqual(t, None[string]().OkOr(err).UnwrapErr(), err)
	AssertTrue(t, HasValue(x))
	AssertTrue(t, HasValue(y.Err()))
	AssertTrue(t, HasValue(y.Err().OkOr(err)))
	AssertEqual(t, err, y.Err().OkOr(err).Unwrap())
	AssertEqual(t, err, y.Err().OkOr(errors.New("")).Unwrap())
	AssertNotEqual(t, err, Err[string](errors.New("asd")).UnwrapErr())
}
