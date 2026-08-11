package versions

import "testing"

func TestValidateVersionNumber(t *testing.T) {
	tests := []struct {
		name          string
		versionNumber int
		wantErr       bool
	}{
		{
			name:          "valid first version",
			versionNumber: 1,
			wantErr:       false,
		},
		{
			name:          "valid later version",
			versionNumber: 7,
			wantErr:       false,
		},
		{
			name:          "zero invalid",
			versionNumber: 0,
			wantErr:       true,
		},
		{
			name:          "negative invalid",
			versionNumber: -1,
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateVersionNumber(test.versionNumber)

			if test.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !test.wantErr && err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	}
}
