package files

import (
	"path/filepath"
	"strings"
)

const MaxTextPreviewBytes int64 = 1024 * 1024

type PreviewKind string

const (
	PreviewImage       PreviewKind = "image"
	PreviewPDF         PreviewKind = "pdf"
	PreviewText        PreviewKind = "text"
	PreviewCode        PreviewKind = "code"
	PreviewVideo       PreviewKind = "video"
	PreviewAudio       PreviewKind = "audio"
	PreviewUnsupported PreviewKind = "unsupported"
)

type PreviewInfo struct {
	Kind        PreviewKind `json:"kind"`
	Previewable bool        `json:"previewable"`
	Inline      bool        `json:"inline"`
}

func ClassifyPreview(contentType string, filename string) PreviewInfo {
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	extension := strings.ToLower(filepath.Ext(filename))

	switch {
	case strings.HasPrefix(contentType, "image/"):
		return PreviewInfo{
			Kind:        PreviewImage,
			Previewable: true,
			Inline:      true,
		}

	case contentType == "application/pdf" || extension == ".pdf":
		return PreviewInfo{
			Kind:        PreviewPDF,
			Previewable: true,
			Inline:      true,
		}

	case strings.HasPrefix(contentType, "video/"):
		return PreviewInfo{
			Kind:        PreviewVideo,
			Previewable: true,
			Inline:      true,
		}

	case strings.HasPrefix(contentType, "audio/"):
		return PreviewInfo{
			Kind:        PreviewAudio,
			Previewable: true,
			Inline:      true,
		}

	case isCodeExtension(extension):
		return PreviewInfo{
			Kind:        PreviewCode,
			Previewable: true,
			Inline:      false,
		}

	case strings.HasPrefix(contentType, "text/") ||
		contentType == "application/json" ||
		contentType == "application/xml" ||
		contentType == "application/javascript":
		return PreviewInfo{
			Kind:        PreviewText,
			Previewable: true,
			Inline:      false,
		}

	default:
		return PreviewInfo{
			Kind:        PreviewUnsupported,
			Previewable: false,
			Inline:      false,
		}
	}
}

func isCodeExtension(extension string) bool {
	switch extension {
	case ".go",
		".py",
		".js",
		".jsx",
		".ts",
		".tsx",
		".java",
		".c",
		".h",
		".cpp",
		".cc",
		".cxx",
		".hpp",
		".cs",
		".rs",
		".rb",
		".php",
		".swift",
		".kt",
		".kts",
		".scala",
		".sh",
		".bash",
		".zsh",
		".fish",
		".sql",
		".html",
		".htm",
		".css",
		".scss",
		".sass",
		".less",
		".vue",
		".svelte",
		".yaml",
		".yml",
		".toml",
		".ini",
		".conf",
		".env",
		".md":
		return true
	default:
		return false
	}
}
