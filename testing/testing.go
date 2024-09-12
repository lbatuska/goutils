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

func AssertTrue(t *testing.T, expected bool) { AssertEqual(t, true, expected) }

func AssertError(t *testing.T, expected error) { AssertTrue(t, expected != nil) }
