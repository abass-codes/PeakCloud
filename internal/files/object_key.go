package files

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func NewObjectKey(ownerID string) string {
	now := time.Now().UTC()

	return fmt.Sprintf(
		"users/%s/%04d/%02d/%s",
		ownerID,
		now.Year(),
		now.Month(),
		uuid.NewString(),
	)
}
