package versions

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abass-codes/peakcloud/internal/files"
	"github.com/gin-gonic/gin"
)

func TestWriteError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "not found",
			err:        ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "forbidden",
			err:        ErrForbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "invalid version",
			err:        ErrInvalidVersion,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "file too large",
			err:        files.ErrFileTooLarge,
			wantStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:       "internal error",
			err:        errors.New("boom"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			writeError(ctx, test.err)

			if recorder.Code != test.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					test.wantStatus,
					recorder.Code,
				)
			}
		})
	}
}
