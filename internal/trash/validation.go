package trash

func ValidateResourceType(resourceType ResourceType) error {
	switch resourceType {
	case ResourceFile, ResourceFolder:
		return nil
	default:
		return ErrInvalidType
	}
}
