package versions

import "testing"

func TestNullString(t *testing.T) {
	if got := nullString(""); got != nil {
		t.Fatalf(
			"expected nil for empty string, got %#v",
			got,
		)
	}

	const value = "etag-123"

	if got := nullString(value); got != value {
		t.Fatalf(
			"expected %q, got %#v",
			value,
			got,
		)
	}
}
