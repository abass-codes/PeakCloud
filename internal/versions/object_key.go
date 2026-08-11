package versions

import (
	"fmt"

	"github.com/google/uuid"
)

func NewObjectKey(
	ownerID string,
	fileID string,
	versionNumber int,
) string {
	return fmt.Sprintf(
		"%s/versions/%s/v%d-%s",
		ownerID,
		fileID,
		versionNumber,
		uuid.NewString(),
	)
}
