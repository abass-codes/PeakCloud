package trash

import "time"

type ResourceType string

const (
	ResourceFile   ResourceType = "file"
	ResourceFolder ResourceType = "folder"
)

type Item struct {
	ID           string       `json:"id"`
	ResourceType ResourceType `json:"resource_type"`
	Name         string       `json:"name"`
	ContentType  string       `json:"content_type,omitempty"`
	SizeBytes    int64        `json:"size_bytes,omitempty"`
	DeletedAt    time.Time    `json:"deleted_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
