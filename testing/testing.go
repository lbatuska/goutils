package testing

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
		t.Logf("✅ [%T](%+v) != [%T](%+v)", expected, expected, true, true)
		return
	}
	t.Errorf("❌ [%T](%+v) == [%T](%+v)", expected, expected, true, true)
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
