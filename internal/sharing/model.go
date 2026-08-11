package sharing

import "time"

type ResourceType string

const (
	ResourceFile   ResourceType = "file"
	ResourceFolder ResourceType = "folder"
)

type Permission string

const (
	PermissionViewer Permission = "viewer"
	PermissionEditor Permission = "editor"
)

type Share struct {
	ID             string       `json:"id"`
	OwnerID        string       `json:"owner_id"`
	RecipientID    string       `json:"recipient_id"`
	RecipientEmail string       `json:"recipient_email,omitempty"`
	ResourceType   ResourceType `json:"resource_type"`
	ResourceID     string       `json:"resource_id"`
	ResourceName   string       `json:"resource_name"`
	Permission     Permission   `json:"permission"`
	AllowDownload  bool         `json:"allow_download"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type PublicLink struct {
	ID            string       `json:"id"`
	OwnerID       string       `json:"owner_id,omitempty"`
	ResourceType  ResourceType `json:"resource_type"`
	ResourceID    string       `json:"resource_id"`
	ResourceName  string       `json:"resource_name"`
	Permission    Permission   `json:"permission"`
	AllowDownload bool         `json:"allow_download"`
	PasswordSet   bool         `json:"password_set"`
	ExpiresAt     *time.Time   `json:"expires_at,omitempty"`
	RevokedAt     *time.Time   `json:"revoked_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type CreateShareRequest struct {
	RecipientEmail string       `json:"recipient_email"`
	ResourceType   ResourceType `json:"resource_type"`
	ResourceID     string       `json:"resource_id"`
	Permission     Permission   `json:"permission"`
	AllowDownload  bool         `json:"allow_download"`
}

type UpdateShareRequest struct {
	Permission    Permission `json:"permission"`
	AllowDownload bool       `json:"allow_download"`
}

type CreatePublicLinkRequest struct {
	ResourceType  ResourceType `json:"resource_type"`
	ResourceID    string       `json:"resource_id"`
	Permission    Permission   `json:"permission"`
	AllowDownload bool         `json:"allow_download"`
	Password      string       `json:"password,omitempty"`
	ExpiresAt     *time.Time   `json:"expires_at,omitempty"`
}

type PublicLinkCreated struct {
	PublicLink
	Token string `json:"token"`
}

type PublicAccessRequest struct {
	Password string `json:"password"`
}

type PublicFile struct {
	ID           string
	ObjectKey    string
	OriginalName string
	ContentType  string
	SizeBytes    int64
}
