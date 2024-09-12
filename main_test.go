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
	Create(&FileLoggerimpl{})
	go Logger().StartLogger()
	Logger().Write("TEST MSG")
	time.Sleep(time.Second * 2)
	x := Some("Test")
	AssertTrue(t, x.Is_some())
	AssertTrue(t, !x.Is_none())
	AssertEqual(t, x.Unwrap(), "Test")
	AssertTrue(t, None_t(1).Is_none())
	err := errors.New("Test")
	y := Err_t(err, "t")
	AssertEqual(t, y.Unwrap_err(), err)
	AssertError(t, y.Unwrap_err())
	z := x.Ok_or(err)
	AssertTrue(t, z.Is_ok())
	AssertEqual(t, None[string]().Ok_or(err).Unwrap_err(), err)
	AssertTrue(t, Has_value(x))
	AssertTrue(t, Has_value(y.Err()))
	AssertTrue(t, Has_value(y.Err().Ok_or(err)))
	AssertEqual(t, err, y.Err().Ok_or(err).Unwrap())
	AssertEqual(t, err, y.Err().Ok_or(errors.New("")).Unwrap())
	AssertNotEqual(t, err, Err[string](errors.New("asd")).Unwrap_err())
}
