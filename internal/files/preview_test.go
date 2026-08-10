package files

import "testing"

func TestClassifyPreview(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		filename    string
		kind        PreviewKind
		previewable bool
	}{
		{
			name:        "image",
			contentType: "image/png",
			filename:    "photo.png",
			kind:        PreviewImage,
			previewable: true,
		},
		{
			name:        "pdf",
			contentType: "application/pdf",
			filename:    "report.pdf",
			kind:        PreviewPDF,
			previewable: true,
		},
		{
			name:        "plain text",
			contentType: "text/plain",
			filename:    "notes.txt",
			kind:        PreviewText,
			previewable: true,
		},
		{
			name:        "go source",
			contentType: "text/plain",
			filename:    "main.go",
			kind:        PreviewCode,
			previewable: true,
		},
		{
			name:        "typescript source",
			contentType: "application/octet-stream",
			filename:    "page.tsx",
			kind:        PreviewCode,
			previewable: true,
		},
		{
			name:        "video",
			contentType: "video/mp4",
			filename:    "demo.mp4",
			kind:        PreviewVideo,
			previewable: true,
		},
		{
			name:        "audio",
			contentType: "audio/mpeg",
			filename:    "song.mp3",
			kind:        PreviewAudio,
			previewable: true,
		},
		{
			name:        "unsupported archive",
			contentType: "application/zip",
			filename:    "archive.zip",
			kind:        PreviewUnsupported,
			previewable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyPreview(tt.contentType, tt.filename)

			if result.Kind != tt.kind {
				t.Fatalf(
					"expected kind %q, got %q",
					tt.kind,
					result.Kind,
				)
			}

			if result.Previewable != tt.previewable {
				t.Fatalf(
					"expected previewable %v, got %v",
					tt.previewable,
					result.Previewable,
				)
			}
		})
	}
}
