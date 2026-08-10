package files

import "time"

type File struct {
	ID           string    `json:"id"`
	OwnerID      string    `json:"-"`
	ObjectKey    string    `json:"-"`
	OriginalName string    `json:"name"`
	ContentType  string    `json:"content_type"`
	SizeBytes    int64     `json:"size_bytes"`
	ETag         string    `json:"etag,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
