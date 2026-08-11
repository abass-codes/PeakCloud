package versions

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParseVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{
			name:    "valid",
			value:   "7",
			want:    7,
			wantErr: false,
		},
		{
			name:    "zero",
			value:   "0",
			wantErr: true,
		},
		{
			name:    "negative",
			value:   "-1",
			wantErr: true,
		},
		{
			name:    "not integer",
			value:   "abc",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(nil)
			ctx.Params = gin.Params{
				{
					Key:   "version",
					Value: test.value,
				},
			}

			got, err := parseVersion(ctx)

			if test.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != test.want {
				t.Fatalf(
					"expected %d, got %d",
					test.want,
					got,
				)
			}
		})
	}
}
