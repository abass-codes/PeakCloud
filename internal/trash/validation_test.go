package trash

import (
	"errors"
	"testing"
)

func TestValidateResourceType(t *testing.T) {
	tests := []struct {
		name         string
		resourceType ResourceType
		wantErr      error
	}{
		{
			name:         "file",
			resourceType: ResourceFile,
			wantErr:      nil,
		},
		{
			name:         "folder",
			resourceType: ResourceFolder,
			wantErr:      nil,
		},
		{
			name:         "empty",
			resourceType: "",
			wantErr:      ErrInvalidType,
		},
		{
			name:         "unknown",
			resourceType: "unknown",
			wantErr:      ErrInvalidType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateResourceType(test.resourceType)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf(
					"expected %v, got %v",
					test.wantErr,
					err,
				)
			}
		})
	}
}
