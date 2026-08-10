package files

import (
	"strings"
	"testing"
)

func TestNewObjectKeyIncludesOwnerNamespace(t *testing.T) {
	key := NewObjectKey("user-123")

	if !strings.HasPrefix(key, "users/user-123/") {
		t.Fatalf("unexpected object key: %s", key)
	}
}

func TestNewObjectKeyIsUnique(t *testing.T) {
	first := NewObjectKey("user-123")
	second := NewObjectKey("user-123")

	if first == second {
		t.Fatal("expected unique object keys")
	}
}
