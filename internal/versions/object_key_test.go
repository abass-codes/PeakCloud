package versions

import (
	"strings"
	"testing"
)

func TestNewObjectKey(t *testing.T) {
	key := NewObjectKey("owner-123", "file-456", 2)

	if !strings.HasPrefix(
		key,
		"owner-123/versions/file-456/v2-",
	) {
		t.Fatalf("unexpected object key: %s", key)
	}
}

func TestNewObjectKeyUnique(t *testing.T) {
	first := NewObjectKey("owner", "file", 2)
	second := NewObjectKey("owner", "file", 2)

	if first == second {
		t.Fatal("expected unique object keys")
	}
}
