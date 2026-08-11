package versions

import "time"

type Version struct {
	ID            string    `json:"id"`
	FileID        string    `json:"file_id"`
	VersionNumber int       `json:"version_number"`
	ObjectKey     string    `json:"-"`
	SizeBytes     int64     `json:"size_bytes"`
	ContentType   string    `json:"content_type"`
	ETag          string    `json:"etag,omitempty"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}
