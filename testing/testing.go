package Testing

import "testing"

func AssertEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected == actual {
		t.Logf("✅ [%T](%+v) == [%T](%+v)", expected, expected, actual, actual)
		return
	}
	t.Errorf("❌ [%T](%+v) != [%T](%+v)", expected, expected, actual, actual)
}

func AssertNotEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected != actual {
		t.Logf("✅ [%T](%+v) != [%T](%+v)", expected, expected, actual, actual)
		return
	}
	t.Errorf("❌ [%T](%+v) == [%T](%+v)", expected, expected, actual, actual)
}

func AssertTrue(t *testing.T, expected bool) {
	t.Helper()
	if expected {
		t.Logf("✅ [%T](%+v) == [%T](%+v)", expected, expected, true, true)
		return
	}
	t.Errorf("❌ [%T](%+v) != [%T](%+v)", expected, expected, true, true)
}

func AssertFalse(t *testing.T, expected bool) {
	t.Helper()
	if !expected {
		t.Logf("✅ [%T](%+v) != [%T](%+v)", expected, expected, false, false)
		return
	}
	t.Errorf("❌ [%T](%+v) == [%T](%+v)", expected, expected, false, false)
}

func AssertError(t *testing.T, expected error) {
	t.Helper()
	if expected != nil {
		t.Logf("✅ [%T](%+v) != [%T](%+v)", expected, expected, nil, nil)
		return
	}
	t.Errorf("❌ [%T](%+v) == [%T](%+v)", expected, expected, nil, nil)
}

func AssertNotError(t *testing.T, expected error) {
	t.Helper()
	if expected == nil {
		t.Logf("✅ [%T](%+v) != [%T](%+v)", expected, expected, nil, nil)
		return
	}
	t.Errorf("❌ [%T](%+v) == [%T](%+v)", expected, expected, nil, nil)
}

func AssertNil[T any](t *testing.T, expected *T) {
	t.Helper()
	if expected == nil {
		t.Logf("✅ [%T](%+v) == [%T](%+v)", expected, expected, nil, nil)
		return
	}
	t.Errorf("❌ [%T](%+v) != [%T](%+v)", expected, expected, nil, nil)
}

func AssertNotNil[T any](t *testing.T, expected *T) {
	t.Helper()
	if expected != nil {
		t.Logf("✅ [%T](%+v) == [%T](%+v)", expected, expected, nil, nil)
		return
	}
	t.Errorf("❌ [%T](%+v) != [%T](%+v)", expected, expected, nil, nil)
}

func internalAssertPanicMessage(f func()) (msg string, b bool) {
	b = true
	defer func() {
		r := recover()
		msg = r.(string)
	}()
	f()

	return msg, false
}

func AssertPanicMessage(t *testing.T, f func(), msg string) {
	t.Helper()
	returnMsg, x := internalAssertPanicMessage(f)
	if x {
		if msg == returnMsg {
			t.Log("✅ Code correctly panicked with correct message!")
		} else {
			t.Error("❌ Code correctly panicked but the message was wrong!")
		}
	} else {
		t.Error("❌ Code was supposed to panic!")
	}
}

func internalAssertPanic(f func()) (b bool) {
	b = true
	defer func() {
		_ = recover()
	}()
	f()

	return false
}

func AssertPanic(t *testing.T, f func()) {
	t.Helper()
	x := internalAssertPanic(f)
	if x {
		t.Log("✅ Code correctly panicked!")
	} else {
		t.Error("❌ Code was supposed to panic!")
	}
}

func AssertNotPanic(t *testing.T, f func()) {
	t.Helper()
	x := internalAssertPanic(f)
	if !x {
		t.Log("✅ Code correctly did not panic!")
	} else {
		t.Error("❌ Code panicked!")
	}
}
