package trash

import "testing"

func TestResourceTypesAreStable(t *testing.T) {
	tests := []struct {
		name  string
		value ResourceType
		valid bool
	}{
		{
			name:  "file",
			value: ResourceFile,
			valid: true,
		},
		{
			name:  "folder",
			value: ResourceFolder,
			valid: true,
		},
		{
			name:  "empty",
			value: ResourceType(""),
			valid: false,
		},
		{
			name:  "unknown",
			value: ResourceType("unknown"),
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateResourceType(test.value)

			if test.valid && err != nil {
				t.Fatalf(
					"expected %q to be valid: %v",
					test.value,
					err,
				)
			}

			if !test.valid && err == nil {
				t.Fatalf(
					"expected %q to be invalid",
					test.value,
				)
			}
		})
	}
}
