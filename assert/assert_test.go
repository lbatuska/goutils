package Assert

import (
	"testing"

	Testing "github.com/lbatuska/goutils/testing"
)

func Test_notNil(t *testing.T) {
	var something *interface{}
	Testing.AssertPanic(t, func() {
		NotNil(something)
	})

	Testing.AssertPanic(t, func() {
		NotNilPtr(something)
	})

	somethingElse := make([]int, 5)
	Testing.AssertNotPanic(t, func() {
		NotNil(&somethingElse)
	})

	Testing.AssertNotPanic(t, func() {
		NotNilPtr(&somethingElse)
	})
}

func Test_nil(t *testing.T) {
	var something *interface{}
	Testing.AssertPanic(t, func() {
		Nil(something)
	})

	Testing.AssertNotPanic(t, func() {
		NilPtr(something)
	})

	somethingElse := make([]int, 5)
	Testing.AssertPanic(t, func() {
		Nil(&somethingElse)
	})

	Testing.AssertPanic(t, func() {
		NilPtr(&somethingElse)
	})
}

func Test_assert(t *testing.T) {
	something := false
	Testing.AssertPanic(t, func() {
		Assert(something)
	})

	somethingElse := true
	Testing.AssertNotPanic(t, func() {
		Assert(somethingElse)
	})
}

func Test_true(t *testing.T) {
	something := false
	Testing.AssertPanic(t, func() {
		True(something)
	})

	somethingElse := true
	Testing.AssertNotPanic(t, func() {
		True(somethingElse)
	})
}

func Test_assertNot(t *testing.T) {
	something := false
	Testing.AssertNotPanic(t, func() {
		AssertNot(something)
	})

	somethingElse := true
	Testing.AssertPanic(t, func() {
		AssertNot(somethingElse)
	})
}

func Test_false(t *testing.T) {
	something := false
	Testing.AssertNotPanic(t, func() {
		False(something)
	})

	somethingElse := true
	Testing.AssertPanic(t, func() {
		False(somethingElse)
	})
}

func Test_equal(t *testing.T) {
	something := "this is something"
	something2 := "this is something"
	Testing.AssertNotPanic(t, func() {
		Equal(something, something2)
	})

	somethingElse := "this is something else"
	somethingElse2 := "this is something completely else"
	Testing.AssertPanic(t, func() {
		Equal(somethingElse, somethingElse2)
	})
}

func Test_notEqual(t *testing.T) {
	something := "this is something"
	something2 := "this is something"
	Testing.AssertPanic(t, func() {
		NotEqual(something, something2)
	})

	somethingElse := "this is something else"
	somethingElse2 := "this is something completely else"
	Testing.AssertNotPanic(t, func() {
		NotEqual(somethingElse, somethingElse2)
	})
}
