package versions

import (
	"errors"
	"fmt"
	"testing"

	"github.com/abass-codes/peakcloud/internal/authorization"
)

func TestMapAuthorizationError(t *testing.T) {
	otherErr := errors.New("database unavailable")

	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "nil",
			err:  nil,
			want: nil,
		},
		{
			name: "resource not found",
			err:  authorization.ErrResourceNotFound,
			want: ErrNotFound,
		},
		{
			name: "wrapped resource not found",
			err: fmt.Errorf(
				"authorize: %w",
				authorization.ErrResourceNotFound,
			),
			want: ErrNotFound,
		},
		{
			name: "forbidden",
			err:  authorization.ErrForbidden,
			want: ErrForbidden,
		},
		{
			name: "wrapped forbidden",
			err: fmt.Errorf(
				"authorize: %w",
				authorization.ErrForbidden,
			),
			want: ErrForbidden,
		},
		{
			name: "other error",
			err:  otherErr,
			want: otherErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := mapAuthorizationError(test.err)

			if test.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", got)
				}
				return
			}

			if !errors.Is(got, test.want) {
				t.Fatalf(
					"expected %v, got %v",
					test.want,
					got,
				)
			}
		})
	}
}
