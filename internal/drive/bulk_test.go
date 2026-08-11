package drive

import (
	"errors"
	"testing"

	"github.com/abass-codes/peakcloud/internal/authorization"
)

func TestResourceTypeForItemFile(t *testing.T) {
	got, err := resourceTypeForItem("file")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != authorization.ResourceFile {
		t.Fatalf("expected file resource, got %q", got)
	}
}

func TestResourceTypeForItemFolder(t *testing.T) {
	got, err := resourceTypeForItem("folder")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != authorization.ResourceFolder {
		t.Fatalf("expected folder resource, got %q", got)
	}
}

func TestResourceTypeForItemRejectsUnknownType(t *testing.T) {
	_, err := resourceTypeForItem("unknown")
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrUnsupportedItemType) {
		t.Fatalf(
			"expected ErrUnsupportedItemType, got %v",
			err,
		)
	}
}
